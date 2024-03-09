# 获取本机ip地址
ipaddr=$(hostname -I | awk '{print $1}')

# 打印ip地址
echo "IP Address: $ipaddr"

# 修改配置文件中所有127.0.0.1为本机ip
sed -i "s|{ipaddr}|$ipaddr|" ./kq_consumer/validation_email/config.yaml