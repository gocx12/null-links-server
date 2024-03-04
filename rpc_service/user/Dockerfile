FROM golang:1.20 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build

COPY . .
RUN go mod download
COPY rpc_service/user/etc /app/etc
RUN go build -ldflags="-s -w" -o /app/rpc_service_user rpc_service/user/user.go


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/rpc_service_user /app/rpc_service_user
# COPY --from=builder /app/etc /app/etc

CMD ["./rpc_service_user", "-f", "etc/user.yaml"]
