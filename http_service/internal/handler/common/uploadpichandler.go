package common

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"sync"
	"time"

	"null-links/http_service/internal/svc"

	"github.com/demdxx/gocast"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	// Register image handling libraries by importing them.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func UploadPicHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the multipart form in the request
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			logx.Error("parse multi part form error: ", err)
			http.Error(w, "parse multi part form error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Retrieve the file from form data
		file, header, err := r.FormFile("image") // "file" is the field name in the form
		if err != nil {
			logx.Error("Retrieve file error: ", err)
			http.Error(w, "Retrieve error:"+err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// save to object storage
		// Minio
		endpoint := svcCtx.Config.MinIO.Endpoint
		accessKeyID := svcCtx.Config.MinIO.AccessKeyID
		secretAccessKey := svcCtx.Config.MinIO.SecretAccessKey
		useSSL := svcCtx.Config.MinIO.UseSSL

		// Initialize minio client object.
		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			logx.Error("init minio client error: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		userId := gocast.ToInt64(r.FormValue("user_id"))
		businessId := gocast.ToInt32(r.FormValue("business_id"))
		lastPicUrl := r.FormValue("last_pic_url")

		// get bucket name
		bucketName, err := UploadPicBusinessId(businessId).bucketName()
		if err != nil {
			logx.Error("UploadPicHandler|unknown business_id: ", businessId)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		// generate unique filename
		objectName, err := generateUniqueFilename(file, userId)
		if err != nil {
			logx.Error("generate unique filename error: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		var wg sync.WaitGroup
		wg.Add(2)
		// upload the pic
		go func() {
			defer wg.Done()
			objectName += filepath.Ext(header.Filename)

			contentType := "application/octet-stream"

			info, err := minioClient.PutObject(r.Context(), bucketName, objectName, file, -1, minio.PutObjectOptions{ContentType: contentType})
			if err != nil {
				logx.Error("minio upload error: ", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				logx.Debugf("Successfully uploaded %s of size %v\n", objectName, info)
			}
		}()

		// remove last pic
		go func() {
			defer wg.Done()
			logx.Debug("lastPicUrl: ", lastPicUrl)
			if lastPicUrl != "" {
				lastFilename := path.Base(lastPicUrl)
				err = minioClient.RemoveObject(context.Background(), bucketName, lastFilename, minio.RemoveObjectOptions{})
				if err != nil {
					logx.Error("minio remove object error: ", err, ", bucketName: ", bucketName, ", objectName: ", lastFilename)
				}
			}
		}()
		wg.Wait()

		// TODO(chancy): 用户量上来后，转移到云存储

		// Return success
		fileInfo := FileInfo{
			Success: true,
			// Key:          info.Key,
			Url: svcCtx.Config.MinIO.DownloadHost + "/" + bucketName + "/" + objectName,
			// ThumbnailUrl: svcCtx.Config.MinIO.DownloadHost + "/" + bucketName + "/" + objectName,
			// Name:         objectName,
			// Type: contentType,
			// Size: info.Size,
			// Error:        "",
			// DeleteUrl:    svcCtx.Config.MinIO.DownloadHost + "/" + bucketName + "/" + objectName,
			// DeleteType:   "DELETE",
		}

		httpx.OkJsonCtx(r.Context(), w, fileInfo)
	}
}

func generateUniqueFilename(file multipart.File, userId int64) (string, error) {
	// 同一纳秒上传完全相同的文件，则生成的文件名是相同的
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// reset file pointer
	file.Seek(0, 0)

	hashInBytes := hash.Sum(nil)[:16]
	md5Hash := hex.EncodeToString(hashInBytes)

	return fmt.Sprintf("%s-%d-%d", md5Hash, time.Now().UnixNano(), userId), nil
}

// FileInfo describes a file that has been uploaded.
type FileInfo struct {
	Success      bool   `json:"success"`
	Key          string `json:"-"`
	Url          string `json:"url,omitempty"`
	ThumbnailUrl string `json:"thumbnail_url,omitempty"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Size         int64  `json:"size"`
	Error        string `json:"error,omitempty"`
	DeleteUrl    string `json:"delete_url,omitempty"`
	DeleteType   string `json:"delete_type,omitempty"`
}

type UploadPicBusinessId int32

const (
	WebSetCover  UploadPicBusinessId = 1
	WebLinkCover UploadPicBusinessId = 2
	Avatar       UploadPicBusinessId = 3
)

/**
 * @return the bucket name for the business id.
 */
func (id UploadPicBusinessId) bucketName() (string, error) {
	switch id {
	case WebSetCover:
		return "websetcover", nil
	case WebLinkCover:
		return "weblinkcover", nil
	case Avatar:
		return "avatar", nil
	default:
		return "unknown", errors.New("unknown business_id")
	}
}
