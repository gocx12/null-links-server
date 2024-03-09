# rsync -r ./http_service root@43.139.0.182:/data/null_tv/http_service

# 获取本机ip地址
ifconfig | grep "inet " | grep -v

# 修改所有


go run rpc_service/webset/webset.go

go run rpc_service/user/user.go

go run http_service/service.go

go run chat_service/service.go

go run kq_consumer/validation_email/validation_email.go

go run kq_consumer/weblink_cover/weblink_cover.go


