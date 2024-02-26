package chat

import (
	"context"

	"null-links/chat_service/internal/svc"
	"null-links/chat_service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportLogic {
	return &ReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReportLogic) Report(req *types.ChatReportReq) (resp *types.ChatReportResp, err error) {
	// todo: add your logic here and delete this line
	
	switch req.type:
	case 1:
	// 举报用户
	case 2:
	// 举报聊天
	case 3:
	// 举报webset
	default:
		logx.Error("unknow report type:", req.type)


	return
}
