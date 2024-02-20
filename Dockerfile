FROM golang:latest as builder

WORKDIR /app
COPY . .

RUN make build

FROM alpine:latest

COPY --from=builder /app/copilot-metrics-exporter /copilot-metrics-exporter

CMD ["/copilot-metrics-exporter"]