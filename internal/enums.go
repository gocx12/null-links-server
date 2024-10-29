package internal

import "errors"

type ChatStatusEnum int64

const (
	ChatValid    ChatStatusEnum = 1
	ChatDeleted  ChatStatusEnum = 2
	ChatReported ChatStatusEnum = 3
	ChatBanned   ChatStatusEnum = 4
)

func (e ChatStatusEnum) Code() int64 {
	switch e {
	case ChatValid:
		return int64(ChatValid)
	case ChatDeleted:
		return int64(ChatDeleted)
	case ChatReported:
		return int64(ChatReported)
	case ChatBanned:
		return int64(ChatBanned)
	default:
		return -1
	}
}

type TopicStatusEnum int64

const (
	TopicValid   TopicStatusEnum = 1
	TopicDeleted TopicStatusEnum = 2
)

func (e TopicStatusEnum) Code() int64 {
	switch e {
	case TopicValid:
		return int64(TopicValid)
	case TopicDeleted:
		return int64(TopicDeleted)
	default:
		return -1
	}
}

type LikeActionTypeEnum int32

const (
	DoLike       LikeActionTypeEnum = 1
	DoCancelLike LikeActionTypeEnum = 2
)

type LikeStatusEnum int32

const (
	Like   LikeStatusEnum = 1
	UnLike LikeStatusEnum = 2
)

func (e LikeStatusEnum) Code() int32 {
	switch e {
	case UnLike:
		return int32(UnLike)
	case Like:
		return int32(Like)
	default:
		return -1
	}
}

type UploadPicBusinessIdEnum int32

const (
	WebSetCover  UploadPicBusinessIdEnum = 1
	WebLinkCover UploadPicBusinessIdEnum = 2
	Avatar       UploadPicBusinessIdEnum = 3
)

/**
 * @return the bucket name for the business id.
 */
func (id UploadPicBusinessIdEnum) BucketName() (string, error) {
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
	WebsetPendReview WebsetStatusEnum = 1 // 待审核
	WebsetPublished  WebsetStatusEnum = 2 // 已发布
	WebsetRejected   WebsetStatusEnum = 3 // 审核未通过
	WebsetDeleted    WebsetStatusEnum = 4 // 已删除
)

func (e WebsetStatusEnum) Code() int32 {
	switch e {
	case WebsetPendReview:
		return int32(WebsetPendReview)
	case WebsetPublished:
		return int32(WebsetPublished)
	case WebsetRejected:
		return int32(WebsetRejected)
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

func (e WeblinkStatusEnum) Code() int32 {
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
