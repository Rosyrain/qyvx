package controller

import (
	"encoding/xml"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"io"
	"net/http"
	"qyvx/dao/mysql"
	"qyvx/logic"
	"qyvx/models"
	"qyvx/pkg/utools"
	"qyvx/pkg/wxbizmsgcrypt"
)

var (
	token          = "dPv2RKxhOBdgC1YmxuUmvLVl"
	receiverId     = "wwf4f0871502d60e9e"
	encodingAeskey = "DcJwJc2nHfwxIlkfMWbbsdrDcOBRST6SRsTJu2hbCtN"
	wxcpt          = wxbizmsgcrypt.NewWXBizMsgCrypt(token, encodingAeskey, receiverId, wxbizmsgcrypt.XmlType)
)

func InviteHandler(c *gin.Context) {
	// 1.参数校验
	p := new(models.InviteParam)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 2.逻辑处理
	// 2.1 发送邀请
	if err := logic.Invite(p); err != nil {
		zap.L().Error("logic.Invite(p) failed", zap.Error(err))
		//	关于发送邀请的错误
		if errors.Is(err, utools.ErrorUrl) {
			ResponseError(c, CodeErrorGithubUrl)
			return
		}
		if errors.Is(err, utools.ErrorGithubUser) {
			ResponseError(c, CodeErrorGithubUrl)
			return
		}
		if errors.Is(err, utools.ErrorGetGithubID) {
			ResponseError(c, CodeErrorInvite)
			return
		}
		if errors.Is(err, logic.ErrorInvite) {
			ResponseError(c, CodeErrorInvite)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	zap.L().Info("logic.Invite(p) success")
	// 2.2 更新企业微信用户信息
	if err := logic.UpdateAlias(p); err != nil {
		zap.L().Error("logic.UpdateAlias(p) failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, codeUserExist)
			return
		}
		if errors.Is(err, logic.ErrorRefreshToken) {
			ResponseErrorWithMsg(c, CoderRefreshToken, "addressBook token")
		}
		if errors.Is(err, logic.ErrorAddressAccessToken) {
			ResponseErrorWithMsg(c, CodeErrorToken, "addressBook token")
			return
		}
		if errors.Is(err, logic.ErrorUpdateAlias) {
			ResponseError(c, CodeServerBusy)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3.返回响应
	ResponseSuccess(c, nil)
}

func HookHandler(c *gin.Context) {
	reqMsgSign := c.DefaultQuery("msg_signature", "")
	reqTimestamp := c.DefaultQuery("timestamp", "")
	reqNonce := c.DefaultQuery("nonce", "")

	body, _ := io.ReadAll(c.Request.Body)

	msg, cryptErr := wxcpt.DecryptMsg(reqMsgSign, reqTimestamp, reqNonce, body)
	//fmt.Println(msg)
	if cryptErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "DecryptMsg fail", "details": cryptErr.ErrMsg})
		return
	}

	var msgContent models.MsgContent
	err := xml.Unmarshal(msg, &msgContent)
	if nil != err {
		zap.L().Error("umarshal failed", zap.Error(err))
	}
	//	//fmt.Println("msgcontent:  ", msgContent, err)
	zap.L().Info("msgContent Info--", zap.Any("spNo: ", msgContent.ApprovalInfo.SpNo))
	zap.L().Info("msgContent Info--", zap.Any("spName: ", msgContent.ApprovalInfo.SpName))
	zap.L().Info("msgContent Info--", zap.Any("spStatus: ", msgContent.ApprovalInfo.SpStatus))
	zap.L().Info("msgContent Info--", zap.Any("qyvxID: ", msgContent.ApprovalInfo.Applyer.UserId))

	zap.L().Warn("Because qyvx-response is required in 10 seconds,so processing tasks in an asynchronous manner.")
	go func() {
		ResponseSuccess(c, nil)
	}()

	ok, inviteParam, err := logic.ParseApproveInfo(msgContent.ApprovalInfo.SpNo)
	if err != nil {
		zap.L().Error("logic.ParseApproveInfo failed", zap.Error(err))
		if errors.Is(err, logic.ErrorRefreshToken) {
			ResponseErrorWithMsg(c, CoderRefreshToken, "github token")
		}
		if errors.Is(err, logic.ErrorGithubAccessToken) {
			ResponseErrorWithMsg(c, CodeErrorToken, "github token")
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	if !ok {
		zap.L().Info("don`t need to exec invite task... ...")
		ResponseErrorWithMsg(c, CodeServerBusy, "The approval information does not require subsequent operations")
		return
	}

	zap.L().Info("exec invite task...")
	// 2.逻辑处理
	// 2.1 发送邀请
	if err := logic.Invite(inviteParam); err != nil {
		zap.L().Error("logic.Invite(p) failed", zap.Error(err))
		//	关于发送邀请的错误
		if errors.Is(err, utools.ErrorUrl) {
			ResponseError(c, CodeErrorGithubUrl)
			return
		}
		if errors.Is(err, utools.ErrorGithubUser) {
			ResponseError(c, CodeErrorGithubUrl)
			return
		}
		if errors.Is(err, utools.ErrorGetGithubID) {
			ResponseError(c, CodeErrorInvite)
			return
		}
		if errors.Is(err, logic.ErrorInvite) {
			ResponseError(c, CodeErrorInvite)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	zap.L().Info("logic.Invite(p) success")
	// 2.2 更新企业微信用户信息
	if err := logic.UpdateAlias(inviteParam); err != nil {
		zap.L().Error("logic.UpdateAlias(p) failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, codeUserExist)
			return
		}
		if errors.Is(err, logic.ErrorRefreshToken) {
			ResponseErrorWithMsg(c, CoderRefreshToken, "addressBook token")
		}
		if errors.Is(err, logic.ErrorAddressAccessToken) {
			ResponseErrorWithMsg(c, CodeErrorToken, "addressBook token")
			return
		}
		if errors.Is(err, logic.ErrorUpdateAlias) {
			ResponseError(c, CodeServerBusy)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 响应已经提前返回

}
