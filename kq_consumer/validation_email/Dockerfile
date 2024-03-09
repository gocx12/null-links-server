FROM golang:1.20 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.io,direct

WORKDIR /build

COPY . .
# 执行前置脚本
RUN chmod +x ./kq_consumer/validation_email/pre.sh
RUN /bin/bash ./kq_consumer/validation_email/pre.sh
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/validation_email ./kq_consumer/validation_email/validation_email.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/validation_email /app/validation_email
COPY --from=builder /build/kq_consumer/validation_email/config.yaml /app/config.yaml


CMD ["./validation_email", "-f", "config.yaml"]