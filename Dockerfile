FROM golang:1.18 as builder
WORKDIR /go/src/gin-consul
COPY . .
ENV SERVICE_NAME = pci-gateway

RUN CGO_ENABLED=0 GOOS=linux  go build -o /go-consul-client main.go
FROM busybox
COPY --from=builder /go-consul-client /go-consul-client
