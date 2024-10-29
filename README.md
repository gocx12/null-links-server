# null-links-server
闹链 N站

以链接将我们连接

Link us with links

# API Document
https://app.apifox.com/project/3613606

# deploy

## environment
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


```bash
vim /usr/local/openresty/nginx/conf/nginx.conf
/usr/local/openresty/bin/openresty -s reload
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


5. kafka UI

1.configure new cluster

name: 随便填
host: kafka
port: 9092

create topic:

6. Mysql

```bash
# 修改密码。注意原来该user的host，需要保持一致才能修改成功
alter user "admin"@"localhost" identified with mysql_native_password by "新密码"

# 修改host
update user set host='localhost' where user='admin';
grant all privileges on *.* to 'admin'@'localhost' with grant option;

flush privilege
```




## start up

### 1. nohup
sudo nohup ./http_service/service > nohup_http_service.log 2>&1 &

### 2. supervisor
以下是在 Ubuntu 下使用 Supervisor 来部署 Go 程序的步骤：

1. 安装 Supervisor
   - 使用以下命令安装 Supervisor：
     ```bash
     sudo apt-get install supervisor
     ```

2. 配置 Supervisor
   - 创建一个 Supervisor 配置文件，例如 `/etc/supervisor/conf.d/null-links.conf`：
     ```ini
     [program:null-links-server]
     command=/data/null-links-server/build/http_service/service -f /data/null-links-server/build/http_service/etc/service.yaml
     directory=/data/null-links-server/build/http_service/
     autostart=true
     autorestart=true
     stderr_logfile=/data/logs/null-links-server.err.log
     stdout_logfile=/data/logs/null-links-server.out.log
     ```
   - 在这个配置中：
     - `command` 指定了要启动的程序及其参数。
     - `directory` 设置了工作目录。
     - `autostart` 和 `autorestart` 确保程序在系统启动时自动启动，并在失败时自动重新启动。
     - `stderr_logfile` 和 `stdout_logfile` 指定了程序的标准错误和标准输出的日志文件位置。

3. 重新加载 Supervisor 配置
   - 使用以下命令重新加载 Supervisor 配置，使新的配置生效：
     ```bash
     sudo supervisorctl reread
     sudo supervisorctl update
     ```

4. 启动程序
   - 使用以下命令启动程序：
     ```bash
     sudo supervisorctl start null-links-server
     ```

5. 检查程序状态
   - 使用以下命令检查程序的状态：
     ```bash
     sudo supervisorctl status null-links-server
     ```

现在，你的 Go 程序应该在 Supervisor 的管理下运行了。你可以通过查看日志文件 `/var/log/null-links-server.out.log` 和 `/var/log/null-links-server.err.log` 来了解程序的运行情况。



# modify
```bash
goctl api go --api http_service/api/main.api --dir http_service

goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_user --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_favorite --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_like --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_relation --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_webset --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_weblink --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_topic --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_chat --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_balance --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_pay_history --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_business --url "root:123456@tcp(127.0.0.1:3306)/db_null_links" &&
goctl model mysql datasource -d http_service/internal/infrastructure/model -t t_advice --url "root:123456@tcp(127.0.0.1:3306)/db_null_links"
```


