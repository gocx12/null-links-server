# rsync -r ./http_service root@43.139.0.182:/data/null_tv/http_service

# 获取本机ip地址
MACHINE="159.138.156.216"
ipaddr=$(hostname -I | awk '{print $1}')
MODE="online"

build(){
  echo '===== Mode:build ====='
  go mod tidy
  # 刪除原有的build目录
  rm -rf ./build

  # 创建目录
  mkdir -p ./build/http_service/etc
  mkdir -p ./build/chat_service/etc
  mkdir -p ./build/rpc_service/user/etc
  mkdir -p ./build/rpc_service/webset/etc
  mkdir -p ./build/kq_consumer/weblink_cover
  mkdir -p ./build/kq_consumer/validation_email

  # 如果MODE=dev，则拷贝_dev配置文件
  if [ $MODE = "dev" ]; then
    cp -r ./http_service/etc/*_dev.yaml ./build/http_service/etc/chat_service.yaml
    cp -r ./chat_service/etc/*_dev.yaml ./build/chat_service/etc/service.yaml
    cp -r ./rpc_service/user/etc/*_dev.yaml ./build/rpc_service/user/etc/user.yaml
    cp -r ./rpc_service/webset/etc/*_dev.yaml ./build/rpc_service/webset/etc/service.yaml
    cp -r ./kq_consumer/weblink_cover/*_dev.yaml ./build/kq_consumer/weblink_cover/config.yaml
    cp -r ./kq_consumer/validation_email/*_dev.yaml ./build/kq_consumer/validation_email/config.yaml
  else
    cp -r ./http_service/etc/*_online.yaml ./build/http_service/etc/service.yaml
    cp -r ./chat_service/etc/*_online.yaml ./build/chat_service/etc/service.yaml
    cp -r ./rpc_service/user/etc/*_online.yaml ./build/rpc_service/user/etc/user.yaml
    cp -r ./rpc_service/webset/etc/*_online.yaml ./build/rpc_service/webset/etc/webset.yaml
    cp -r ./kq_consumer/weblink_cover/*_online.yaml ./build/kq_consumer/weblink_cover/config.yaml
    cp -r ./kq_consumer/validation_email/*_online.yaml ./build/kq_consumer/validation_email/config.yaml
  fi
  cp -r ./kq_consumer/validation_email/validation_code_page.html ./build/kq_consumer/validation_email/validation_code_page.html

  go build -ldflags="-s -w" -o ./build/chat_service/chat_service ./chat_service/service.go
  go build -ldflags="-s -w" -o ./build/kq_consumer/weblink_cover/weblink_cover ./kq_consumer/weblink_cover/weblink_cover.go
  go build -ldflags="-s -w" -o ./build/kq_consumer/validation_email/validation_email ./kq_consumer/validation_email/validation_email.go
  go build -ldflags="-s -w" -o ./build/rpc_service/user/rpc_user ./rpc_service/user/user.go
  go build -ldflags="-s -w" -o ./build/rpc_service/webset/rpc_webset ./rpc_service/webset/webset.go
  go build -ldflags="-s -w" -o ./build/http_service/http_service ./http_service/service.go
}

run() {
  echo '===== Mode:run ====='
  # 启动服务
  nohup ./build/rpc_service/user/rpc_user -f ./build/rpc_service/user/etc/user.yaml > /dev/null 2>&1 &
  nohup ./build/rpc_service/webset/rpc_webset -f ./build/rpc_service/webset/etc/webset.yaml > /dev/null 2>&1 &

  nohup ./build/http_service/http_service -f ./build/http_service/etc/service.yaml > /dev/null 2>&1 &
  nohup ./build/chat_service/chat_service -f ./build/chat_service/etc/service.yaml  > /dev/null 2>&1 &

  nohup ./build/kq_consumer/weblink_cover/weblink_cover -f ./build/kq_consumer/weblink_cover/config.yaml > /dev/null 2>&1 &
  nohup ./build/kq_consumer/validation_email/validation_email -f ./build/kq_consumer/validation_email/config.yaml -t ./build/kq_consumer/validation_email/validation_code_page.html > /dev/null 2>&1 &
}

syncFile() {
  echo '===== Mode:syncFile ====='
  rsync -avzP --progress  ./build.sh root@$MACHINE:/data/null-links-server && echo "Sync successful"
  rsync -avzP --progress  ./build root@$MACHINE:/data/null-links-server && echo "Sync successful"
}

while [ -n "$1" ]
do
    case "$1" in
    "--run")
        run
        ;;
    "--build")
        build
        # syncFile
        ;;
    "--sync")
        syncFile
        ;;
    "--publish")
        build
        syncFile
        ;;
    # "--vendor")
    #     buildWithVendor
    #     ;;
    # "--build_cron")
    #     if [ ! -d $VendorDir ];then
    #         buildCronWithDownload
    #     else
    #         buildCronWithVendor
    #     fi
    #     syncFile
    #     ;;
    # "--build")
    #     if [ ! -d $VendorDir ];then
    #         buildWithDownload
    #     else
    #         buildWithVendor
    #     fi
    #     ;;
    # "--start")
    #     echo '=====Mode:start ====='
    #     COUNT=`ps -ef |grep -v "grep" |grep -v "defunct"|grep -c $BINARY_NAME`
    #     if [ $COUNT -lt 1 ]; then
    #         nohup $BINARY_FILE > /dev/null 2>&1 &
    #     fi
    #     ;;
    "--stop")
        echo '=====Mode:stop ====='
        ps -ef | grep chat_service | grep -v "grep" | awk '{print $2}' | xargs kill
        ps -ef | grep http_service | grep -v "grep" | awk '{print $2}' | xargs kill
        ps -ef | grep rpc_user | grep -v "grep" | awk '{print $2}' | xargs kill
        ps -ef | grep rpc_webset | grep -v "grep" | awk '{print $2}' | xargs kill
        ps -ef | grep weblink_cover | grep -v "grep" | awk '{print $2}' | xargs kill
        ps -ef | grep validation_email | grep -v "grep" | awk '{print $2}' | xargs kill
        ;;
    # "--rebuild")
    #     echo '=====Mode:rebuild with Shutdown====='
    #     if [ ! -d $VendorDir ];then
    #         buildWithDownload
    #     else
    #         buildWithVendor
    #     fi
    #     runProject
    #     ;;
    # "--force_rebuild")
    #     echo '=====Mode:force_rebuild====='
    #     if [ ! -d $VendorDir ];then
    #         buildWithDownload
    #     else
    #         buildWithVendor
    #     fi
    #     ps -ef | grep $BINARY_NAME | grep -v "grep" |  awk '{print $2}' | xargs kill -s 9
    #     nohup $BINARY_FILE > /dev/null 2>&1 &
    #     ;;
    # "--reload")
    #     PID_FILE=$(cat $CURDIR/conf/app.yaml | grep PidFile | egrep -o -o [a-Z_]+.txt)
    #     pid=$(cat $PROJECTDIR/etc/${PID_FILE})
    #     ps -p "$pid"
    #     if [ $? -eq 0 ]; then
    #         kill -HUP $pid # SIGHUP
    #     else
    #         nohup $BINARY_FILE > /dev/null 2>&1 &
    #     fi
    #     ;;
    *)
        echo '没有包含第一参数'
        ;;
    esac
    shift
done