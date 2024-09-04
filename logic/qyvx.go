package logic

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"qyvx/dao/mysql"
	"qyvx/models"
	snowflake "qyvx/pkg/snowflask"
	"qyvx/pkg/utools"
)

func Invite(param *models.InviteParam) error {
	qyvxID := param.QyvxID
	githubUrl := param.GithubUrl
	// 1. 判断数据库中是否有该用户
	Eerr := mysql.CheckUserExistByQyvxID(qyvxID)
	if !(errors.Is(Eerr, mysql.ErrorUserNotExist) || errors.Is(Eerr, mysql.ErrorUserExist)) {
		return Eerr
	}

	// 1.0 解析githubUrl正确性并得到githubID
	githubName, err := utools.ParseGithubUrl(githubUrl)
	if err != nil {
		return err
	}
	githubID, err := utools.GetGithubID(githubName, githubUrl, githubToken)
	if err != nil {
		return err
	}
	// 1.1 如果存在该用户，则更改githubID并发送邀请
	if errors.Is(Eerr, mysql.ErrorUserExist) {
		if err := mysql.UpdateGithubIDByQyvxID(qyvxID, githubName, githubID); err != nil {
			return err
		}

		if err := invite(githubID); err != nil {
			return err
		}
		return nil
	}

	// 1.2 如果该用户不存在，则插入一条新数据并发送邀请
	if errors.Is(Eerr, mysql.ErrorUserNotExist) {
		uid := snowflake.GenID()
		if err := mysql.InsertUserInfo(uid, githubID, qyvxID, githubName); err != nil {
			return err
		}

		if err := invite(githubID); err != nil {
			return err
		}
		return nil
	}
	return ErrorInvite
}

// UpdateAlias 将githubID以别名的形式更新到企业微信
func UpdateAlias(p *models.InviteParam) error {
	qyvxID := p.QyvxID
	zap.L().Info("UpdateAlias --", zap.Any("qyvxID: ", qyvxID))
	// 1. 判断数据库中是否有该用户
	Eerr := mysql.CheckUserExistByQyvxID(qyvxID)
	if !errors.Is(Eerr, mysql.ErrorUserExist) {
		return Eerr
	}

	// 2. 获取githubName -- 注：一切以数据库中的githubName为准
	githubName, err := mysql.GetGithubNameByQyvxID(qyvxID)
	if err != nil {
		return err
	}
	zap.L().Info("UpdateAlias --", zap.Any("githubName(Mysql): ", githubName))
	// 2. 更改信息
	if err := updateAlias(qyvxID, githubName); err != nil {
		return err
	}

	return nil
}

// ParseApproveInfo 解析审批获得的数据
func ParseApproveInfo(spNo string) (bool, *models.InviteParam, error) {
	// 1.获取详情信息
	body, err := getQyvxApproveInfo(spNo)
	if err != nil {
		return false, nil, err
	}

	// 2.判断并解析内容
	zap.L().Info("parse data...")
	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return false, nil, err
	}
	infoMap := jsonData["info"].(map[string]interface{})
	spName := infoMap["sp_name"].(string)
	spStatus := infoMap["sp_status"].(float64)
	if spName != "github账号绑定" || spStatus != 2 {
		zap.L().Info("spName or spStatus inaccuracy", zap.Any("spName", spName), zap.Any("spStatus", spStatus))
		return false, nil, nil
	}

	// 3.结果返回(需要处理)
	applyerMap := infoMap["applyer"].(map[string]interface{})
	qyvxID := applyerMap["userid"].(string)

	applyDataMap := infoMap["apply_data"].(map[string]interface{})
	contentList := applyDataMap["contents"].([]interface{})
	githubUrlInfo := contentList[0].(map[string]interface{})
	githubUrlValue := githubUrlInfo["value"].(map[string]interface{})
	githubUrl := githubUrlValue["text"].(string)

	param := &models.InviteParam{
		GithubUrl: githubUrl,
		QyvxID:    qyvxID,
	}

	return true, param, nil
}
