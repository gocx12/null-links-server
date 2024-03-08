FROM golang:1.20 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.io,direct

WORKDIR /build

COPY . .
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/webset ./rpc_service/webset/webset.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/webset /app/webset
COPY --from=builder /build/rpc_service/webset/etc /app/etc

CMD ["./webset", "-f", "etc/webset.yaml"]