// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2

package types

type AdviceReq struct {
	Reason      string `json:"reason"`
	PicUrl      string `json:"pic_url"`
	ContactInfo string `json:"contact_info"`
}

type AdviceResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
}

type Chat struct {
	ChatID     int64  `json:"chat_id"`
	WebsetID   int64  `json:"webset_id"`
	UserID     int64  `json:"user_id"`
	UserName   string `json:"username"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	TopicId    int64  `json:"topic_id"`
	TopicTitle string `json:"topic_title"`
}

type ChatGetAllTopicReq struct {
	WebsetId int64 `form:"webset_id"`
}

type ChatGetAllTopicResp struct {
	StatusCode int32   `json:"status_code"`
	StatusMsg  string  `json:"status_msg,optional"`
	TopicList  []Topic `json:"topic_list"`
}

type ChatGetTopicReq struct {
	TopicId int64 `form:"topic_id"`
}

type ChatGetTopicResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
	ChatList   []Chat `json:"chat_list"`
	TopicTitle string `json:"topic_title"`
}

type ChatHistoryReq struct {
	WebsetID   int64  `form:"webset_id"`
	Type       int32  `form:"type"`
	Page       int32  `form:"page,optional"`
	PageSize   int32  `form:"page_size,optional"`
	LastChatId int64  `form:"last_chat_id,optional"`
	StartTime  string `form:"start_time,optional"`
	EndTime    string `form:"end_time,optional"`
	Keyword    string `form:"keyword,optional"`
}

type ChatHistoryResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
	ChatList   []Chat `json:"chat_list"`
}

type ChatLikeReq struct {
	ChatId int64 `form:"chat_id"`
}

type ChatLikeResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
}

type ChatLinkReq struct {
	ChatId  int64 `form:"chat_id"`
	TopicId int64 `form:"topic_id"`
}

type ChatLinkResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
}

type ChatWsReq struct {
	Token    string `form:"token,optional"`
	WebsetID int64  `form:"webset_id"`
	UserID   int64  `form:"user_id"`
}

type ChatWsResp struct {
}

type CheckUsernameReq struct {
	Username string `form:"username,optional"`
}

type CheckUsernameResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
	Result     int32  `json:"result"`
}

type FavoriteActionReq struct {
	ActionType int32 `form:"action_type"`
	WebsetID   int64 `form:"webset_id"`
}

type FavoriteActionResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
}

type FavoriteListReq struct {
}

type FavoriteListResp struct {
	StatusCode int32    `json:"status_code"`
	StatusMsg  string   `json:"status_msg,optional"`
	WebsetList []Webset `json:"webset_list"`
}

type FeedReq struct {
	LatestTime string `form:"latest_time,optional"`
	Page       int32  `form:"page,optional"`
	PageSize   int32  `form:"page_size,optional"`
}

type FeedResp struct {
	StatusCode int32         `json:"status_code"`
	StatusMsg  string        `json:"status_msg,optional"`
	NextTime   int64         `json:"next_time"`
	WebsetList []WebsetShort `json:"webset_list"`
}

type FriendUser struct {
	User    User   `form:"friend_user_info"`
	Message string `form:"message,optional"`
	MsgType int    `form:"msg_type"`
}

type GetValidationCodeReq struct {
	Email string `form:"email"`
}

type GetValidationCodeResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
}

type LikeActionReq struct {
	ActionType int32 `form:"action_type"`
	WebsetId   int64 `form:"webset_id"`
}

type LikeActionResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
}

type LoginReq struct {
	Username  string `json:"username,optional"`
	UserEmail string `json:"user_email,optional"`
	Password  string `json:"password"`
}

type LoginResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
	Token      string `json:"token"`
	UserID     int64  `json:"user_id"`
	Username   string `json:"username"`
	AvatarUrl  string `json:"avatar_url"`
}

type MessageActionReq struct {
	ActionType string `form:"action_type"`
	Content    string `form:"content"`
	ToUserId   string `form:"to_user_id"`
}

type MessageActionResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
}

type ModifyReq struct {
	Username  string `json:"username,optional"`
	AvatarUrl string `json:"avatar_url,optional"`
}

type ModifyResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
}

type PayInfoReq struct {
}

type PayInfoResp struct {
	StatusCode          int32             `json:"status_code"`
	StatusMsg           string            `json:"status_msg,optional"`
	Balance             float64           `json:"balance"`
	WithdrawHistoryList []WithdrawHistory `json:"withdraw_history_list"`
}

type PingReq struct {
}

type PingResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
	Data       string `json:"data"`
}

type PublishActionReq struct {
	ActionType  int32            `json:"action_type"`
	Category    int32            `json:"category"`
	AuthorId    int64            `json:"author_id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	CoverUrl    string           `json:"cover_url"`
	WebLinkList []WebLinkPublish `json:"weblink_list"`
	WebsetId    int64            `json:"webest_id,optional"`
}

type PublishActionResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
}

type PublishListReq struct {
	Page     int32 `form:"page"`
	PageSize int32 `form:"page_size"`
	Tag      int32 `form:"tag"`
}

type PublishListResp struct {
	StatusCode int32           `json:"status_code"`
	StatusMsg  string          `json:"status_msg,optional"`
	WebsetList []PublishWebset `json:"webset_list"`
	Total      int32           `json:"total"`
}

type PublishWebset struct {
	Id            int64  `json:"id"`
	Title         string `json:"title"`
	Describe      string `json:"describe"`
	CoverUrl      string `json:"cover_url"`
	CreatedAt     string `json:"created_at"`
	ViewCount     int64  `json:"view_count"`
	LikeCount     int64  `json:"like_count"`
	FavoriteCount int64  `json:"favorite_count"`
	Status        int32  `json::"status"`
}

type RegisterReq struct {
	Username       string `json:"username"`
	UserEmail      string `json:"user_email"`
	ValidationCode string `json:"validation_code"`
	Password       string `json:"password"`
}

type RegisterResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
	Token      string `json:"token"`
	UserId     int64  `json:"user_id"`
}

type RelationActionReq struct {
	ActionType string `form:"action_type"`
	ToUserId   int64  `form:"to_user_id"`
}

type RelationActionResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
}

type RelationFollowListReq struct {
}

type RelationFollowListResp struct {
	StatusCode string       `json:"status_code"`
	StatusMsg  string       `json:"status_msg,optional"`
	UserList   []FriendUser `json:"user_list"`
}

type RelationFollowerListReq struct {
}

type RelationFollowerListResp struct {
	StatusCode string `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
	UserList   []User `json:"follower_list"`
}

type RelationFriendListReq struct {
}

type RelationFriendListResp struct {
	StatusCode string `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
	UserList   []User `json:"friend_list"`
}

type ReportReq struct {
	UserId     int64  `json:"user_id"`
	BusinessId int32  `json:"business_id"`
	Id         int64  `json:"id"`
	Reason     string `json:"reason"`
}

type ReportResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
}

type Topic struct {
	TopicId    int64  `json:"topic_id"`
	TopicTitle string `json:"topic_title"`
}

type UploadFileReq struct {
	UserId int64 `json:"user_id"`
}

type UploadFileResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
}

type UploadPicReq struct {
}

type UploadPicResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
	Success    bool   `json:"success"`
}

type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	AvatarUrl     string `json:"avatar_url,optional"`
	BackgroundUrl string `json:"background_url"`
	FollowCount   int64  `json:"follow_count,optional"`
	FollowerCount int64  `json:"follower_count,optional"`
	IsFollow      bool   `json:"is_follow"`
	Signature     string `json:"signature,optional"`
	WorkCount     int64  `json:"work_count,optional"`
}

type UserInfoReq struct {
	UserId int64 `form:"user_id"`
}

type UserInfoResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
	User       User   `json:"user_info"`
}

type UserShort struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url,optional"`
}

type WebLink struct {
	Id         int64  `json:"id"`
	Describe   string `json:"describe"`
	Url        string `json:"url"`
	AuthorInfo User   `json:"author_info"`
	CoverURL   string `json:"cover_url"`
}

type WebLinkPublish struct {
	Url         string `json:"url"`
	Description string `json:"description"`
	CoverUrl    string `json:"cover_url,optional"`
}

type Webset struct {
	Id            int64     `json:"id"`
	Title         string    `json:"title"`
	Describe      string    `json:"describe"`
	AuthorInfo    User      `json:"author_info"`
	CoverURL      string    `json:"cover_url"`
	ViewCount     int64     `json:"view_count"`
	LikeCount     int64     `json:"like_count"`
	IsLike        bool      `json:"is_like"`
	FavoriteCount int64     `json:"favorite_count"`
	IsFavorite    bool      `json:"is_favorite"`
	WebLinkList   []WebLink `json:"weblink_list"`
}

type WebsetInfoReq struct {
	UserId   int64 `form:"user_id"`
	WebsetId int64 `form:"webset_id"`
}

type WebsetInfoResp struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,optional"`
	WebsetInfo Webset `json:"webset_info"`
}

type WebsetShort struct {
	Id            int64     `json:"id"`
	Title         string    `json:"title"`
	Describe      string    `json:"describe"`
	CoverUrl      string    `json:"cover_url"`
	CreatedAt     string    `json:"created_at"`
	AuthorInfo    UserShort `json:"author_info"`
	ViewCount     int64     `json:"view_count"`
	LikeCount     int64     `json:"like_count"`
	FavoriteCount int64     `json:"favorite_count"`
	IsLike        bool      `json:"is_like"`
}

type WithdrawHistory struct {
	Id     string
	Amount float64
	Time   string
}

type WithdrawReq struct {
}

type WithdrawResp struct {
	StatusCode         int32  `json:"status_code"`
	StatusMsg          string `json:"status_msg,optional"`
	WXFriendCodePicUrl string `json:"wx_friend_code_pic_url"`
}
