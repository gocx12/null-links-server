package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/demdxx/gocast"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"null-links/kq_consumer/weblink_cover/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type Conf struct {
	DataSource         string
	WlCoverKqConsumser kq.KqConf
	MinIO              struct {
		Endpoint        string
		AccessKeyID     string
		SecretAccessKey string
		UseSSL          bool
		DownloadHost    string
	}
}

var (
	c            Conf
	configFile          = flag.String("f", "kq_consumer/weblink_cover/config.yaml", "the config file")
	format       string = ".png"
	weblinkModel model.TWeblinkModel
)

func main() {
	flag.Parse()
	conf.MustLoad(*configFile, &c)
	weblinkModel = model.NewTWeblinkModel(sqlx.NewMysql(c.DataSource))

	logx.Info("weblink screeshot kq consumer starting")
	q := kq.MustNewQueue(c.WlCoverKqConsumser, kq.WithHandle(sreenShot))
	defer q.Stop()
	q.Start()
}

func sreenShot(k, v string) error {
	logx.Debug("screeshot", "k: ", k, " v: ", v)

	// format: webset_id::link_id::url
	websetId := gocast.ToInt64(v)

	weblinkDb, err := weblinkModel.FindByWebsetId(context.Background(), websetId)
	if err != nil {
		logx.Error("get weblink from db error: ", err, " ,webset id:", weblinkDb)
		return err
	}
	var wg sync.WaitGroup
	for _, weblink := range weblinkDb {
		wg.Add(1)
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
			endpoint := c.MinIO.Endpoint
			accessKeyID := c.MinIO.AccessKeyID
			secretAccessKey := c.MinIO.SecretAccessKey
			useSSL := c.MinIO.UseSSL
			bucketName := "weblinkcover"

			// Initialize minio client object.
			minioClient, err := minio.New(endpoint, &minio.Options{
				Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
				Secure: useSSL,
			})
			if err != nil {
				logx.Error("init minio client error: ", err)
			}
			objectName, err := generateUniqueFilename(fileReader, websetId)
			if err != nil {
				logx.Error("generate unique filename error: ", err)
			}
			objectName += format
			contentType := "application/octet-stream"
			_, err = minioClient.PutObject(context.Background(), bucketName, objectName, fileReader, -1, minio.PutObjectOptions{ContentType: contentType})
			if err != nil {
				logx.Error("minio upload error: ", err)
			}
			coverUrl := c.MinIO.DownloadHost + "/" + bucketName + "/" + objectName
			// 封面地址落库
			weblinkModel.UpdateCoverUrl(context.Background(), websetId, weblink.LinkId, coverUrl)
			wg.Done()
		}(weblink)
	}
	wg.Wait()
	return nil
}

func generateUniqueFilename(file io.ReadSeeker, id int64) (string, error) {
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
