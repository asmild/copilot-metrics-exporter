FROM scratch
COPY ca-certificates.crt /etc/ssl/certs/
COPY copilot-metrics-exporter /copilot-metrics-exporter
ENTRYPOINT ["/copilot-metrics-exporter"]