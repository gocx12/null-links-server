package common

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"null-links/http_service/internal/svc"

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
		uploadHandler(w, r)
	}
}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form in the request
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		logx.Error("parse multi part form error: ", err)
		http.Error(w, "parse multi part form error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the file from form data
	file, _, err := r.FormFile("image") // "file" is the field name in the form
	if err != nil {
		logx.Error("Retrieve file error: ", err)
		http.Error(w, "Retrieve error:"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// save to object storage
	// Minio
	endpoint := "127.0.0.1:9000"
	accessKeyID := "swjTlz02UIYY7iXCyCQM"
	secretAccessKey := "ayNq4sY19U9PM9TTBFuGKKvP62h6kfNpYx4DQeb5"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logx.Error("init minio client error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// Make a new bucket called testbucket.
	bucketName := "avatar"

	// Upload the test file
	// Change the value of filePath if the file is in another location
	objectName := "testdata"
	contentType := "application/octet-stream"

	// Upload the test file with FPutObject
	info, err := minioClient.PutObject(r.Context(), bucketName, objectName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		logx.Error("minio upload error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	logx.Debugf("Successfully uploaded %s of size %v\n", objectName, info)

	// Generate download link
	URL, err := minioClient.PresignedGetObject(context.Background(), bucketName, objectName, time.Hour*1, nil)
	if err != nil {
		logx.Error("minio get error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// TODO(chancy): 用户上来后，转移到云存储
	// Return success
	fileInfo := FileInfo{
		Success:      true,
		Key:          info.Key,
		Url:          URL.String(),
		ThumbnailUrl: "",
		Name:         objectName,
		Type:         contentType,
		Size:         info.Size,
		Error:        "",
		DeleteUrl:    URL.String(),
		DeleteType:   "DELETE",
	}

	httpx.OkJsonCtx(r.Context(), w, fileInfo)
}

func post(w http.ResponseWriter, r *http.Request) {
	// 从请求中读取文件
	fileInfos := make([]*FileInfo, 0)
	mr, err := r.MultipartReader()
	if err != nil {
		logx.Error("read multipart err: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Form, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		logx.Error("parse query err: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for {
		// var part *multipart.Part
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			logx.Error("read part err: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if name := part.FormName(); name != "" {

			if part.FileName() != "" {
				fileInfos = append(fileInfos, uploadFile(w, part))
			} else {
				r.Form[name] = append(r.Form[name], getFormValue(part))
			}
		}
	}

	js, err := json.Marshal(fileInfos)
	if err != nil {
		logx.Error("marshal fileInfos err: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if redirect := r.FormValue("redirect"); redirect != "" {
		http.Redirect(w, r, fmt.Sprintf(
			redirect,
			escape(string(js)),
		), http.StatusFound)
		return
	}
	jsonType := "application/json"
	if strings.Contains(r.Header.Get("Accept"), jsonType) {
		w.Header().Set("Content-Type", jsonType)
	}

	httpx.OkJsonCtx(r.Context(), w, fileInfos)
}

// Config holds configuration settings for an UploadHandler.
type Config struct {
	MinFileSize        int // bytes
	MaxFileSize        int // bytes
	AcceptFileTypes    *regexp.Regexp
	ExpirationTime     int // seconds
	ThumbnailMaxWidth  int // pixels
	ThumbnailMaxHeight int // pixels
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

const (
	IMAGE_TYPES = "image/(gif|p?jpeg|(x-)?png)"
)

var (
	DefaultConfig = Config{
		MinFileSize:        1,               // bytes
		MaxFileSize:        1024 * 1024 * 3, // bytes 3MB
		AcceptFileTypes:    regexp.MustCompile(IMAGE_TYPES),
		ExpirationTime:     300,
		ThumbnailMaxWidth:  80,
		ThumbnailMaxHeight: 80,
	}
	imageRegex        = regexp.MustCompile(IMAGE_TYPES)
	FileNotFoundError = errors.New("File Not Found")
)

// uploadFile handles the upload of a single file from a multipart form.
func uploadFile(w http.ResponseWriter, p *multipart.Part) (fi *FileInfo) {
	ctx := context.Background()

	fi = &FileInfo{
		Name: p.FileName(),
		Type: p.Header.Get("Content-Type"),
	}

	conf := DefaultConfig

	// Validate file type
	// if !conf.AcceptFileTypes.MatchString(fi.Type) {
	// 	fi.Error = "acceptFileTypes"
	// 	return
	// }

	isImage := imageRegex.MatchString(fi.Type)

	// Copy into buffers for save and thumbnail generation
	//
	// Max + 1 for LimitedReader size, so we can detect below if file size is
	// greater than max.
	lr := &io.LimitedReader{R: p, N: int64(conf.MaxFileSize + 1)}
	var bSave bytes.Buffer  // Buffer to be saved
	var bThumb bytes.Buffer // Buffer to be thumbnailed
	var wr io.Writer
	if isImage {
		wr = io.MultiWriter(&bSave, &bThumb)
	} else {
		wr = &bSave
	}
	_, err := io.Copy(wr, lr)
	if err != nil {
		logx.Error("Error copying file:", err)
		fi.Error = "readError"
		return
	}

	// Validate file size
	size := bSave.Len()
	if size < conf.MinFileSize {
		logx.Error("File failed validation: too small.", size, conf.MinFileSize)
		fi.Error = "minFileSize"
		return
	} else if size > conf.MaxFileSize {
		logx.Error("File failed validation: too large.", size, conf.MaxFileSize)
		fi.Error = "maxFileSize"
		return
	}

	// save to object storage
	// Minio
	endpoint := "127.0.0.1:9000"
	accessKeyID := "swjTlz02UIYY7iXCyCQM"
	secretAccessKey := "ayNq4sY19U9PM9TTBFuGKKvP62h6kfNpYx4DQeb5"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logx.Error(err)
	}
	// Make a new bucket called testbucket.
	bucketName := "avatar"
	// location := "cn-east-8"

	// err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	// if err != nil {
	// 	// Check to see if we already own this bucket (which happens if you run this twice)
	// 	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
	// 	if errBucketExists == nil && exists {
	// 		logx.Debugf("We already own %s\n", bucketName)
	// 	} else {
	// 		logx.Error(err)
	// 	}
	// } else {
	// 	logx.Infof("Successfully created %s\n", bucketName)
	// }

	// Upload the test file
	// Change the value of filePath if the file is in another location
	objectName := "testdata"
	filePath := "/tmp/testdata"
	contentType := "application/octet-stream"

	// Upload the test file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		logx.Error(err)
	}

	logx.Debugf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	// TODO(chancy): 用户上来后，转移到云存储

	// Set URLs in FileInfo
	u := &url.URL{
		Path: "h.Prefix" + "/",
	}
	uString := u.String()
	fi.Url = uString + escape(string(fi.Key)) + "/" +
		escape(string(fi.Name))
	fi.DeleteUrl = fi.Url
	fi.DeleteType = "DELETE"

	fi.ThumbnailUrl = uString + "thumbnails/" + escape(string(fi.Key))
	// Create thumbnail 略缩图
	// if isImage && size > 0 {
	// 	_, err = h.createThumbnail(fi, &bThumb)
	// 	if err != nil {
	// 		logx.Error("Error creating thumbnail:", err)
	// 	}
	// 	// If we wanted to save thumbnails to peristent storage, this would be
	// 	// a good spot to do it.
	// }
	return
}

func getFormValue(p *multipart.Part) string {
	var b bytes.Buffer
	io.CopyN(&b, p, int64(1<<20)) // Copy max: 1 MiB
	return b.String()
}

func escape(s string) string {
	return strings.Replace(url.QueryEscape(s), "+", "%20", -1)
}

// // createThumbnail generates a thumbnail and adds it to the cache.
// func createThumbnail(fi *FileInfo, r io.Reader) (data []byte, err error) {
// 	defer func() {
// 		if rec := recover(); rec != nil {
// 			logx.Error(rec)
// 			// 1x1 pixel transparent GIf, bas64 encoded:
// 			s := "R0lGODlhAQABAIAAAP///////yH5BAEKAAEALAAAAAABAAEAAAICTAEAOw=="
// 			data, _ = base64.StdEncoding.DecodeString(s)
// 			fi.ThumbnailUrl = "data:image/gif;base64," + s
// 		}
// 		h.Cache.Set(fi.Key, string(data), 0, 0, h.Conf.ExpirationTime)
// 	}()
// 	img, _, err := image.Decode(r)
// 	if err != nil {
// 		panic(err)
// 	}
// 	if bounds := img.Bounds(); bounds.Dx() > h.Conf.ThumbnailMaxWidth ||
// 		bounds.Dy() > h.Conf.ThumbnailMaxHeight {
// 		w, h := h.Conf.ThumbnailMaxWidth, h.Conf.ThumbnailMaxHeight
// 		if bounds.Dx() > bounds.Dy() {
// 			h = bounds.Dy() * h / bounds.Dx()
// 		} else {
// 			w = bounds.Dx() * w / bounds.Dy()
// 		}
// 		img = resize.Resize(img, img.Bounds(), w, h)
// 	}
// 	var b bytes.Buffer
// 	err = png.Encode(&b, img)
// 	if err != nil {
// 		panic(err)
// 	}
// 	data = b.Bytes()
// 	fi.ThumbnailUrl = "data:image/png;base64," +
// 		base64.StdEncoding.EncodeToString(data)
// 	return
// }
