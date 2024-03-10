FROM golang:1.20 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.io,direct

WORKDIR /build

COPY . .
# 执行前置脚本
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/chat_service ./chat_service/service.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/chat_service /app/chat_service

# 只需复制etc目录到当前工作目录即可
COPY --from=builder /build/chat_service/etc /app/etc

CMD ["./chat_service", "-f", "etc/service.yaml"]