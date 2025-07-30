
FROM golang:1.23.4-alpine3.21 AS builder
ARG BIN=copilot-metrics-exporter
ARG SRC=./cmd/copilot-exporter

WORKDIR /app
COPY . .

RUN echo "Building ${BIN}..." && \
    go get -d -v ./... && \
    go install -v ./... && \
    go test -v ./... && \
    go build -o ${BIN} ${SRC}/main.go

FROM alpine:3.21.3 AS runtime
ENV TZ=UTC

RUN addgroup -S exporter && adduser -S exporter -G exporter
USER exporter
COPY --from=builder /app/copilot-metrics-exporter /copilot-metrics-exporter

CMD ["/copilot-metrics-exporter"]

