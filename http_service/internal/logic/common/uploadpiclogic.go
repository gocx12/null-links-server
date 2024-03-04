package common

import (
	"context"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"

	// Register image handling libraries by importing them.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type UploadPicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadPicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadPicLogic {
	return &UploadPicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadPicLogic) UploadPic(req *types.UploadPicReq) (resp *types.UploadPicResp, err error) {
	// todo: add your logic here and delete this line

	resp.Success = false

	resp.Success = true
	return
}
