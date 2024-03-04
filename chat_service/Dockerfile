FROM golang:1.20 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

# COPY go.mod .
# COPY go.sum .
COPY . .
RUN go mod download
# RUN go build -ldflags="-s -w" -o /app/chat_service chat_service/chat_service/service.go
RUN go build -ldflags="-s -w" -o /app/chat_service ./chat_service/service.go 

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/chat_service /app/chat_service
COPY --from=builder /app/etc /app/etc

EXPOSE 8081
CMD ["./chat_service", "-f", "etc/service.yaml"]
