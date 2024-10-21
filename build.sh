MACHINE="159.138.156.216"
PROJECTDIR=$(pwd)
PROJECTNAME="null-link-server"
IMAGE="null-links-server"

# 打包镜像
build() {
  echo '===== Mode:build ====='
  rm -rf ./build
  mkdir -p ./build
  mkdir -p ./build/http_service
  mkdir -p ./build/kq_consumer/validation_email
  mkdir -p ./build/kq_consumer/weblink_cover
  # mkdir -p ./build/cron

  cp -r ./http_service/etc ./build/http_service
  cp ./kq_consumer/validation_email/config.yaml ./build/kq_consumer/validation_email/config.yaml
  cp ./kq_consumer/validation_email/validation_code_page.html ./build/kq_consumer/validation_email/validation_code_page.html
  cp ./kq_consumer/weblink_cover/config.yaml ./build/kq_consumer/weblink_cover/config.yaml
  
  cd $PROJECTDIR/build/http_service
  go build $PROJECTDIR/http_service/service.go
  
  cd $PROJECTDIR/build/kq_consumer/validation_email
  go build $PROJECTDIR/kq_consumer/validation_email/validation_email.go

  cd $PROJECTDIR/build/kq_consumer/weblink_cover
  go build $PROJECTDIR/kq_consumer/weblink_cover/weblink_cover.go

  # cd $PROJECTDIR/build/cron
  # go build  $PROJECTDIR/cron/main.go
}

# 上传文件
syncFile() {
  cd $PROJECTDIR
  echo '===== Mode:syncFile ====='
  tar -czvf build.tar.gz build/
  
  rsync -avzP --progress ./build.tar.gz root@$MACHINE:/data/null-links-server/ && echo "Sync successful"

  ssh root@$MACHINE "cd /data/null-links-server && tar -xzvf build.tar.gz"

  rm build.tar.gz
}

restart() {
  echo '===== Mode:restart ====='
  ssh root@$MACHINE "pm2 restart 0"
}

while [ -n "$1" ]
do
  case "$1" in
  "--run")
    echo '=====Mode:run ====='
     ssh root@$MACHINE "cd /data/null-links-web/build && nohup npm run start &"
    ;;
  "--build")
    build
    ;;
  "--sync")
    syncFile
    ;;
  "--publish")
    build
    syncFile
    ;;
  "unzip")
    echo '=====Mode:unzip ====='
    ssh root@$MACHINE "cd /data/null-links-web && tar -xzvf $IMAGE.tar.gz"
    ;;
  "--stop")
    echo '=====Mode:stop ====='
    docker stop $IMAGE
    ;;
  *)
    echo '没有包含第一参数'
    ;;
  esac
  shift
done