# null-links-server
闹链 N站

以链接将我们连接

Link us with links

# API Document
https://app.apifox.com/project/3613606

# deploy
1. 安装golang
```bash

```

2. 安装chrome内核，用于截图功能
```bash

```

3. docker容器化部署依赖
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
```
docker compose -f compose.yaml up --remove-orphans -d
```
