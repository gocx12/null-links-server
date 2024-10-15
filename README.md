# null-links-server
闹链 N站

以链接将我们连接

Link us with links

# API Document
https://app.apifox.com/project/3613606

# deploy
1. golang 安装
```bash

```

2. 安装chrome内核(用于截图功能)

```bash

```

3. openresty 安装
https://openresty.org/cn/linux-packages.html#ubuntu
```bash
sudo apt-get -y install --no-install-recommends wget gnupg ca-certificates lsb-release

wget -O - https://openresty.org/package/pubkey.gpg | sudo apt-key add -

echo "deb http://openresty.org/package/ubuntu $(lsb_release -sc) main" \
| sudo tee /etc/apt/sources.list.d/openresty.list

sudo apt-get update

sudo apt-get -y install openresty
```

4. docker容器化部署依赖
腾讯云教程: https://cloud.tencent.com/document/product/1207/45596

ubuntu docker engine 安装:
```bash
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo apt-get remove $pkg; done

# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

sudo docker run hello-world
```

docker compose 安装：
```bash
sudo apt-get update
sudo apt-get install docker-compose-plugin

docker compose version
```


容器启动
```bash
docker compose -f compose.yaml up --remove-orphans -d
```


# modify
```
goctl api go --api http_service/api/main.api --dir http_service

goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_user --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_webset --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_weblink --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_like --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_favorite --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_topic --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_chat --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_relation --url "root:123456@tcp(127.0.0.1:3306)/db_null_links"
```


# kafka UI
1. configure new cluster

name: 随便填
host: kafka
port: 9092