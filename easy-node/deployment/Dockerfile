FROM alpine:3.19

ENV GOEXPERIMENT=arenas
ENV QUILIBRIUM_SIGNATURE_CHECK=false
ENV GOMAXPROCS=4

COPY node /usr/local/bin/node

RUN chmod +x /usr/local/bin/node

WORKDIR /root

ENTRYPOINT ["node", "--prometheus-server", "0.0.0.0:8080"]