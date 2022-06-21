FROM registry.access.redhat.com/ubi8/go-toolset:1.17.7 as builder

COPY . ./
RUN go mod download
RUN go build

# Use minimal to get ca-certificates, otherwise communication with AWS fails
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

EXPOSE 9091

WORKDIR /
COPY --from=builder /opt/app-root/src/privatelink-infrastructure-monitor /opt/privatelink-infrastructure-monitor

ENTRYPOINT ["/opt/privatelink-infrastructure-monitor"]
