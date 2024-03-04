FROM golang:1.20 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build

COPY . .
RUN go mod download
COPY rpc_service/webset/etc /app/etc
RUN go build -ldflags="-s -w" -o /app/rpc_service_webset rpc_service/webset/webset.go


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/rpc_service_webset /app/rpc_service_webset
# COPY --from=builder /app/etc /app/etc

CMD ["./rpc_service_webset", "-f", "etc/webset.yaml"]


# FROM docker.registry.kgidc.cn/base_mirror/kg_centos_7.5:go1.20_v1 as build
# RUN mkdir -p /build
# COPY ./ /build
# WORKDIR /build
# RUN /bin/sh -c make.sh ...(编译参数)
# FROM docker.registry.kgidc.cn/base_mirror/kg_centos_7.5:openresty1.13_py2.7_v3
# RUN mkdir -p /data1/go_service/asrservice
# COPY --from=build /build/outer /data1/go_service/asrservice
# WORKDIR /data1/go_service/asrservice