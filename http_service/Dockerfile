FROM golang:1.20 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build

# ADD go.mod .
# ADD go.sum .
COPY . .
RUN go mod download

COPY http_service/etc /app/etc
RUN go build -ldflags="-s -w" -o /app/http_service http_service/service.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/http_service /app/http_service
COPY --from=builder /app/etc /app/etc

CMD ["./http_service", "-f", "etc/service.yaml"]
