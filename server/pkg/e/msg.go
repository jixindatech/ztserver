package e

var MsgFlags = map[int]string{
	SUCCESS: "OK",

	ERROR:   "失败",
	ERROR_1: "非法Token",
	ERROR_2: "其他用户登录",
	ERROR_3: "Token过期",

	InvalidParams: "请求参数错误",

	AddUserFailed:          "增加用户失败",
	GetUserFailed:          "获取用户失败",
	PutUserFailed:          "修改用户失败",
	DeleteUserFailed:       "删除用户失败",
	SaveUserResourceFailed: "修改用户资源失败",
	GetUserResourceFailed:  "获取用户资源失败",
	SendUserMailFailed:     "发送邮件失败",

	AddResourceFailed:    "增加资源失败",
	GetResourceFailed:    "获取资源失败",
	PutResourceFailed:    "修改资源失败",
	DeleteResourceFailed: "删除资源失败",

	AddEmailFailed: "增加邮箱失败",
	GetEmailFailed: "获取邮箱失败",
	PutEmailFailed: "修改邮箱失败",

	AddSSLFailed:    "增加证书失败",
	GetSSLFailed:    "获取证书失败",
	PutSSLFailed:    "修改证书失败",
	DeleteSSLFailed: "删除证书失败",

	AddUpstreamFailed:    "添加转发失败",
	GetUpstreamFailed:    "获取转发失败",
	PutUpstreamFailed:    "修改转发失败",
	DeleteUpstreamFailed: "删除转发失败",

	GetGwEventsFailed: "获取访问日志失败",

	GetWsEventsFailed: "获取网络日志失败",

	AddRouterFailed:    "添加路由失败",
	GetRouterFailed:    "获取路由失败",
	PutRouterFailed:    "修改路由失败",
	DeleteRouterFailed: "删除路由失败",

}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
