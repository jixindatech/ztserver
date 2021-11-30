package e

const (
	SUCCESS = 200

	ERROR   = 300
	ERROR_1 = 301
	ERROR_2 = 302
	ERROR_3 = 303

	InvalidParams = 400

	AddUserFailed          = 1001
	GetUserFailed          = 1002
	PutUserFailed          = 1003
	DeleteUserFailed       = 1004
	SaveUserResourceFailed = 1005
	GetUserResourceFailed  = 1006
	SendUserMailFailed     = 1007

	AddResourceFailed    = 2001
	GetResourceFailed    = 2002
	PutResourceFailed    = 2003
	DeleteResourceFailed = 2004

	AddEmailFailed = 3001
	GetEmailFailed = 3002
	PutEmailFailed = 3003

	AddSSLFailed    = 4001
	GetSSLFailed    = 4002
	PutSSLFailed    = 4003
	DeleteSSLFailed = 4004

	AddUpstreamFailed    = 5001
	GetUpstreamFailed    = 5002
	PutUpstreamFailed    = 5003
	DeleteUpstreamFailed = 5004

	GetGwEventsFailed = 6001

	GetWsEventsFailed = 7001

	AddRouterFailed    = 8001
	GetRouterFailed    = 8002
	PutRouterFailed    = 8003
	DeleteRouterFailed = 8004
)
