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

	AddCertFailed:    "增加证书失败",
	GetCertFailed:    "获取证书失败",
	PutCertFailed:    "修改证书失败",
	DeleteCertFailed: "删除证书失败",

	AddProxyFailed:    "添加代理失败",
	GetProxyFailed:    "获取代理失败",
	PutProxyFailed:    "修改代理失败",
	DeleteProxyFailed: "删除代理失败",

	GetGwEventsFailed: "获取访问日志失败",

	GetWsEventsFailed: "获取网络日志失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
