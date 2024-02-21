FROM golang:1.22.0 AS builder

WORKDIR /app
COPY . .

RUN make build

FROM alpine:3.19.1

COPY --from=builder /app/copilot-metrics-exporter /copilot-metrics-exporter

CMD ["/copilot-metrics-exporter"]