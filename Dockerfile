FROM golang:1.25.4 AS builder
ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY internal/ ./internal/

RUN go build -o snmp-simulator

FROM scratch
COPY --from=builder /app/snmp-simulator /snmp-simulator
ENTRYPOINT ["/snmp-simulator"]