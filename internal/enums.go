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
