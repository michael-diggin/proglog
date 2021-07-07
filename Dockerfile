FROM alpine as build

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
    wget -qO./grpc_health_probe \
    https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x ./grpc_health_probe
COPY ./proglog ./proglog

FROM gcr.io/distroless/static
WORKDIR /root/
COPY --from=build ./proglog ./proglog
COPY --from=build ./grpc_health_probe ./grpc_health_probe
ENTRYPOINT [ "./proglog" ]
