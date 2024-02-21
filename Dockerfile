
FROM golang:1.22.0-alpine3.19 AS builder
ARG BIN=copilot-metrics-exporter
ARG SRC=./cmd/copilot-exporter

WORKDIR /app
COPY . .

RUN echo "Building ${BIN}..." && \
    go get -d -v ./... && \
    go install -v ./... && \
    go build -o ${BIN} ${SRC}/main.go

FROM alpine:3.19.1 AS final
ENV TZ=UTC

COPY --from=builder /app/copilot-metrics-exporter /copilot-metrics-exporter

CMD ["/copilot-metrics-exporter"]

