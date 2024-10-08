package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeServerBusy
	CodeErrorGithubUrl
	CodeErrorGithubUser
	CodeErrorGetID
	CodeErrorUpdateAlias
	CodeErrorQyvxName
	CodeErrorInvite
	codeUserExist
	CoderRefreshToken
	CodeErrorToken
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:          "success",
	CodeInvalidParam:     "请求参数错误",
	CodeServerBusy:       "服务繁忙",
	CodeErrorGithubUrl:   "github url 错误",
	CodeErrorGithubUser:  "github 用户不存在",
	CodeErrorGetID:       "获取github ID失败",
	CodeErrorUpdateAlias: "更新企业微信昵称失败",
	CodeErrorQyvxName:    "获取企业微信名失败",
	CodeErrorInvite:      "发送邀请失败",
	codeUserExist:        "企业微信用户不存在",
	CoderRefreshToken:    "更新token失败",
	CodeErrorToken:       "token 无效",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
