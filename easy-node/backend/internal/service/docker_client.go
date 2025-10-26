package service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

type DockerClientService struct {
	client *client.Client
}

func NewDockerClientService() (*DockerClientService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &DockerClientService{client: cli}, nil
}

func (s *DockerClientService) Close() error {
	return s.client.Close()
}

// 获取容器状态
func (s *DockerClientService) GetContainerStatus(containerName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	containers, err := s.client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return "", fmt.Errorf("failed to list containers: %w", err)
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if strings.TrimPrefix(name, "/") == containerName {
				return container.State, nil
			}
		}
	}

	return "not_found", nil
}

// 启动容器
func (s *DockerClientService) StartContainer(containerName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	containerID, err := s.getContainerID(containerName)
	if err != nil {
		return err
	}

	if err := s.client.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container %s: %w", containerName, err)
	}

	zap.L().Info("Container started", zap.String("container", containerName))
	return nil
}

// 停止容器
func (s *DockerClientService) StopContainer(containerName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	containerID, err := s.getContainerID(containerName)
	if err != nil {
		return err
	}

	timeout := 30
	if err := s.client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		return fmt.Errorf("failed to stop container %s: %w", containerName, err)
	}

	zap.L().Info("Container stopped", zap.String("container", containerName))
	return nil
}

// 重启容器
func (s *DockerClientService) RestartContainer(containerName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	containerID, err := s.getContainerID(containerName)
	if err != nil {
		return err
	}

	timeout := 30
	if err := s.client.ContainerRestart(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		return fmt.Errorf("failed to restart container %s: %w", containerName, err)
	}

	zap.L().Info("Container restarted", zap.String("container", containerName))
	return nil
}

// 获取容器日志
func (s *DockerClientService) GetContainerLogs(containerName string, lines int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	containerID, err := s.getContainerID(containerName)
	if err != nil {
		return "", err
	}

	tail := fmt.Sprintf("%d", lines)
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
		Timestamps: true,
	}

	reader, err := s.client.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read container logs: %w", err)
	}

	return string(logs), nil
}

// 启动 Docker Compose 项目
func (s *DockerClientService) StartComposeProject(projectName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	// 获取项目的所有容器
	filterArgs := filters.NewArgs()
	filterArgs.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))
	containers, err := s.client.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filterArgs,
	})
	if err != nil {
		return fmt.Errorf("failed to list project containers: %w", err)
	}

	// 启动所有容器
	for _, ctr := range containers {
		if ctr.State != "running" {
			if err := s.client.ContainerStart(ctx, ctr.ID, types.ContainerStartOptions{}); err != nil {
				zap.L().Error("Failed to start container", zap.String("container", ctr.Names[0]), zap.Error(err))
			}
		}
	}

	zap.L().Info("Compose project started", zap.String("project", projectName))
	return nil
}

// 停止 Docker Compose 项目
func (s *DockerClientService) StopComposeProject(projectName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	// 获取项目的所有容器
	filterArgs := filters.NewArgs()
	filterArgs.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))
	containers, err := s.client.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filterArgs,
	})
	if err != nil {
		return fmt.Errorf("failed to list project containers: %w", err)
	}

	// 停止所有容器
	timeout := 30
	for _, ctr := range containers {
		if ctr.State == "running" {
			if err := s.client.ContainerStop(ctx, ctr.ID, container.StopOptions{Timeout: &timeout}); err != nil {
				zap.L().Error("Failed to stop container", zap.String("container", ctr.Names[0]), zap.Error(err))
			}
		}
	}

	zap.L().Info("Compose project stopped", zap.String("project", projectName))
	return nil
}

// 获取容器ID
func (s *DockerClientService) getContainerID(containerName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	containers, err := s.client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return "", fmt.Errorf("failed to list containers: %w", err)
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if strings.TrimPrefix(name, "/") == containerName {
				return container.ID, nil
			}
		}
	}

	return "", fmt.Errorf("container %s not found", containerName)
}

// 获取所有节点容器的状态
func (s *DockerClientService) GetAllNodeStatus(projectName string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filterArgs := filters.NewArgs()
	filterArgs.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))
	containers, err := s.client.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filterArgs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	statusMap := make(map[string]string)
	for _, ctr := range containers {
		if len(ctr.Names) > 0 {
			name := strings.TrimPrefix(ctr.Names[0], "/")
			statusMap[name] = ctr.State
		}
	}

	return statusMap, nil
}