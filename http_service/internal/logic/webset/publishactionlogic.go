package webset

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"null-links/http_service/internal/infrastructure/model"
	"null-links/http_service/internal/svc"
	"null-links/http_service/internal/types"

	"null-links/internal"

	"github.com/demdxx/gocast"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PublishActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishActionLogic {
	return &PublishActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 发布类型
type PublishActionTypeEnum int32

const (
	Publish PublishActionTypeEnum = 1
	Update  PublishActionTypeEnum = 2
	Delete  PublishActionTypeEnum = 3
)

// webset状态
type WebsetStatusEnum int32

const (
	WebsetPendReview   WebsetStatusEnum = 1 // 待审核
	WebsetPublished    WebsetStatusEnum = 2 // 已发布
	WebsetReviewUnpass WebsetStatusEnum = 3 // 审核未通过
	WebsetDeleted      WebsetStatusEnum = 4 // 已删除
)

func (e WebsetStatusEnum) code() int32 {
	switch e {
	case WebsetPendReview:
		return int32(WebsetPendReview)
	case WebsetPublished:
		return int32(WebsetPublished)
	case WebsetReviewUnpass:
		return int32(WebsetReviewUnpass)
	case WebsetDeleted:
		return int32(WebsetDeleted)
	default:
		return -1
	}
}

// weblink状态
type WeblinkStatusEnum int32

const (
	WeblinkPendReview   WeblinkStatusEnum = 1 // 待审核
	WeblinkPublished    WeblinkStatusEnum = 2 // 已发布
	WeblinkReviewUnpass WeblinkStatusEnum = 3 // 审核未通过
	WeblinkDeleted      WeblinkStatusEnum = 4 // 已删除
)

func (e WeblinkStatusEnum) code() int32 {
	switch e {
	case WeblinkPendReview:
		return int32(WeblinkPendReview)
	case WeblinkPublished:
		return int32(WeblinkPublished)
	case WeblinkReviewUnpass:
		return int32(WeblinkReviewUnpass)
	case WeblinkDeleted:
		return int32(WeblinkDeleted)
	default:
		return -1
	}
}

func (l *PublishActionLogic) PublishAction(req *types.PublishActionReq) (resp *types.PublishActionResp, err error) {
	resp = &types.PublishActionResp{}

	switch PublishActionTypeEnum(req.ActionType) {
	case Publish:
		err = l.doPublish(req)
	case Update:
		err = l.doUpdate(req)
	case Delete:
		err = l.doDelete(req)
	default:
		resp.StatusCode = internal.StatusParamErr
		resp.StatusMsg = "unkown action type"
	}

	if err != nil {
		logx.Error("PublishAction|err=", err)
		resp.StatusCode = internal.StatusErr
		resp.StatusMsg = "fail"
		return
	}

	resp.StatusCode = internal.StatusSuccess
	resp.StatusMsg = "success"
	return
}

func (l *PublishActionLogic) doPublish(req *types.PublishActionReq) (err error) {
	// 构造webset
	websetDb := model.TWebset{
		Title:       req.Title,
		Description: req.Description,
		AuthorId:    req.AuthorId,
		CoverUrl:    req.CoverUrl,
		Category:    0, // 分区，暂时不用
		ViewCnt:     0,
		LikeCnt:     0,
		FavoriteCnt: 0,
		Status:      gocast.ToInt64(WebsetPendReview.code()),
	}

	// 插入weblinks
	webLinkListDb := make([]model.TWeblink, 0, len(req.WebLinkList))
	for i, webLink := range req.WebLinkList {
		// url 格式化
		if !strings.HasPrefix(webLink.Url, "http") || !strings.HasPrefix(webLink.Url, "https") {
			webLink.Url = "https://" + webLink.Url
		}
		webLinkListDb = append(webLinkListDb, model.TWeblink{
			LinkId:      int64(i),
			WebsetId:    req.WebsetId,
			AuthorId:    req.AuthorId,
			Description: webLink.Description,
			Url:         webLink.Url,
			CoverUrl:    webLink.CoverUrl,
			Status:      gocast.ToInt64(WeblinkPendReview), // status==2 待审核
		})
	}

	// weblink 与 webset 在同一个数据库中，因此可以使用本地事务
	err = l.svcCtx.WebsetModel.GetConn().TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 1.插入webset
		r, err := l.svcCtx.WebsetModel.InsertTrans(l.ctx, &websetDb, session)
		if err != nil {
			return err
		}
		if rowsAffected, err := r.RowsAffected(); err != nil {
			return err
		} else if rowsAffected == 0 {
			return fmt.Errorf("insert webset failed, rows affected: %d", rowsAffected)
		}

		lastInsertId, err := r.LastInsertId()
		if err != nil {
			return err
		}
		for i := range webLinkListDb {
			webLinkListDb[i].WebsetId = lastInsertId
		}

		// 2.截图作为weblink封面
		// // kafka pusher
		// data := gocast.ToString(lastInsertId)
		// if err := l.svcCtx.WlCoverKqPusher.Push(data); err != nil {
		// 	logx.Error("wlCoverKqPusher push error:", err)
		// 	return err
		// }
		go func() {
			l.sreenShot(lastInsertId)
		}()

		// 3.批量插入weblink
		r, err = l.svcCtx.WeblinkModel.BulkInsertTrans(l.ctx, webLinkListDb, session)
		if err != nil {
			return err
		}
		if rowsAffected, err := r.RowsAffected(); err != nil {
			return err
		} else if rowsAffected == 0 {
			return fmt.Errorf("insert weblinks failed, rows affected: %d", rowsAffected)
		}

		return nil
	})

	if err != nil {
		logx.Error("insert weblinks failed, err: ", err)
		return
	}

	return nil
}

func (l *PublishActionLogic) doUpdate(req *types.PublishActionReq) (err error) {
	// 更新
	r, err := l.svcCtx.WebsetModel.UpdateWebsetInfo(l.ctx, &model.TWebset{
		Id:          req.WebsetId,
		Title:       req.Title,
		Description: req.Description,
		AuthorId:    req.AuthorId,
		CoverUrl:    req.CoverUrl,
		UpdatedAt:   time.Now(),
		Category:    0,
		Status:      gocast.ToInt64(WebsetPendReview.code()), // status==2 待审核
	})
	if rowsAffected, err := r.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return fmt.Errorf("insert weblinks failed, rows affected: %d", rowsAffected)
	}
	return
}

func (l *PublishActionLogic) doDelete(req *types.PublishActionReq) (err error) {
	// 软删除
	err = l.svcCtx.WebsetModel.UpdateStatus(l.ctx,
		gocast.ToInt64(WeblinkDeleted.code()),
		req.WebsetId,
	)
	return
}

func (l *PublishActionLogic) sreenShot(websetId int64) error {

	weblinkDb, err := l.svcCtx.WeblinkModel.FindByWebsetId(context.Background(), websetId)
	if err != nil {
		logx.Error("get weblink from db error: ", err, " ,webset id:", weblinkDb)
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(weblinkDb))
	for _, weblink := range weblinkDb {
		go func(weblink *model.TWeblink) {
			// 截图
			page := rod.New().MustConnect().MustPage(weblink.Url).MustWaitLoad()
			file, err := page.Screenshot(true, &proto.PageCaptureScreenshot{
				Format: "png",
			})
			fileReader := bytes.NewReader(file)
			if err != nil {
				logx.Error("screenshot error: ", err, " ,webset id: ", websetId, " ,link id:", weblink.LinkId)
			}

			// 上传图片
			endpoint := l.svcCtx.Config.MinIO.Endpoint
			accessKeyID := l.svcCtx.Config.MinIO.AccessKeyID
			secretAccessKey := l.svcCtx.Config.MinIO.SecretAccessKey
			useSSL := l.svcCtx.Config.MinIO.UseSSL
			bucketName := "weblinkcover"

			// Initialize minio client object.
			minioClient, err := minio.New(endpoint, &minio.Options{
				Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
				Secure: useSSL,
			})
			if err != nil {
				logx.Error("init minio client error: ", err)
			}
			objectName, err := l.generateUniqueFilename(fileReader, websetId)
			if err != nil {
				logx.Error("generate unique filename error: ", err)
			}
			objectName += ".png"
			contentType := "application/octet-stream"
			_, err = minioClient.PutObject(context.Background(), bucketName, objectName, fileReader, -1, minio.PutObjectOptions{ContentType: contentType})
			if err != nil {
				logx.Error("minio upload error: ", err)
			}
			coverUrl := l.svcCtx.Config.MinIO.DownloadHost + "/" + bucketName + "/" + objectName
			// 封面地址落库
			l.svcCtx.WeblinkModel.UpdateCoverUrl(context.Background(), websetId, weblink.LinkId, coverUrl)
			wg.Done()
		}(weblink)
	}
	wg.Wait()
	return nil
}

func (l *PublishActionLogic) generateUniqueFilename(file io.ReadSeeker, id int64) (string, error) {
	// 同一纳秒上传完全相同的文件，则生成的文件名是相同的
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// reset file pointer
	file.Seek(0, 0)

	hashInBytes := hash.Sum(nil)[:16]
	md5Hash := hex.EncodeToString(hashInBytes)

	return fmt.Sprintf("%s-%d-%d", md5Hash, time.Now().UnixNano(), id), nil
}
