FROM golang:1.20 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.io,direct

WORKDIR /build

COPY . .
RUN go mod download
RUN go build -ldflags="-s -w" -o /app/weblink_cover ./kq_consumer/weblink_cover/weblink_cover.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/weblink_cover /app/weblink_cover
COPY --from=builder /build/kq_consumer/weblink_cover/config.yaml /app/config.yaml

CMD ["./weblink_cover", "-f", "config.yaml"]