package consts

const (
	UserId       = "user_id" // 用户id
	Admin        = "$Admin^"
	Expire       = "exp"          // 过期值
	Root         = "root"         // 根节点
	Top          = "top"          // 置顶
	Activation   = "Activation"   // 启用
	Verification = "Verification" // 审核
	Delete       = "Delete"       // 删除，回收站
	Ban          = "Ban"          // 封禁
	Draft        = "Draft"        // 草稿
	Public       = "public"       // 公共
	Private      = "private"      // 私有
	Unread       = "unread"       // 未读
	Approved     = "approved"     // 批准
	Rejected     = "rejected"     // 驳回
	Layout       = "2006-01-02 15:04:05"
	// 通知类型
	NotifySystem  = "notify_system"
	NotifyComment = "notify_comment"
	NotifyArticle = "notify_article"
	NotifyCollect = "notify_collect"
	NotifyLike    = "notify_like"
	NotifySocial  = "notify_social"
	NotifyUsers   = "notify_users"

	ActionLike     = "like"
	ActionView     = "view"
	ActionShare    = "share"
	ActionFavorite = "favorite"
	ActionComment  = "comment"

	Online = "online"
)

// 社区redis缓存键
const (
	UserInfo    = "user_info"
	UserSetting = "user_setting"
	UserKey     = "user_key"
	UserPassKey = "user_pass_key"

	ArticleInfo          = "article_info"
	ArticleUserList      = "article_user_list"
	ArticleDraftList     = "article_draft_list"
	ArticlePopularity    = "article_popularity"
	ArticlePartition     = "article_partition"
	ArticlePartitionList = "article_partition_list"

	ArticleRecommendHistory = "article_recommend_history"
	ArticleHotTopics        = "article_hot_topics"

	CommentFloor     = "comment_floor"
	CommentLeaf      = "comment_leaf"
	CommentLeafCount = "comment_leaf_count"
	CommentLike      = "comment_like"

	Social            = "social"
	FollowerList      = "follower_list"
	FollowList        = "follow_list"
	MutualFriendsList = "mutual_friends_list"
	ArticleActionUser = "article_action_user"
	ArticleActionInfo = "article_action_info"
	GlobalSearchIndex = "global_search_index"
)

// queue name
const (
	QueueAction       = "action_queue"
	QueueActionFailed = "failed_action_queue"

	QueueVerifyArticle       = "verify_article_queue"
	QueueVerifyArticleFailed = "failed_verify_article_queue"
)

const (
	FriendStateless      int32 = iota // 无状态
	FriendStatusAccepted              // 已关注
	FriendStatusBlocked               // 已拉黑
	FriendStatusRemoved               // 已移除
)

func GetFriendStatusDescription(status int32) string {
	switch status {
	case FriendStateless:
		return "无状态"
	case FriendStatusAccepted:
		return "已关注"
	case FriendStatusBlocked:
		return "已拉黑"
	case FriendStatusRemoved:
		return "已移除"
	default:
		return "未知状态"
	}
}

const (
	Authorize = "Authorization"
	Secret    = "wxzx_zjyt_IM_secret_awhs71324gl3ksd879d"
	PassKey   = "pass_key"
	LocalHost = "127.0.0.1;localhost;0.0.0.0"
)
