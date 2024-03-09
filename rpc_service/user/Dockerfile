FROM golang:1.20 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.io,direct

WORKDIR /build

COPY . .
# 执行前置脚本
RUN chmod +x ./rpc_service/user/pre.sh
RUN /bin/bash ./rpc_service/user/pre.sh
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/user ./rpc_service/user/user.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/user /app/user
COPY --from=builder /build/rpc_service/user/etc /app/etc

CMD ["./user", "-f", "etc/user.yaml"]