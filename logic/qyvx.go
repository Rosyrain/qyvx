package logic

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"qyvx/dao/mysql"
	"qyvx/models"
	"qyvx/pkg/utools"
)

func Invite(param *models.InviteParam) error {
	qyvxID := param.QyvxID
	githubUrl := param.GithubUrl
	teamName := param.TeamName
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

	// 1.1 如果存在该用户，则更改githubID
	if errors.Is(Eerr, mysql.ErrorUserExist) {
		if err := mysql.UpdateGithubIDByQyvxID(qyvxID, githubName, githubID); err != nil {
			return err
		}
	}

	// 1.2 如果该用户不存在，则插入一条新数据
	if errors.Is(Eerr, mysql.ErrorUserNotExist) {
		// 根据qyvx_id获取name
		name, err := getqyvxName(qyvxID)
		if err != nil {
			return err
		}
		// 插入数据
		if err := mysql.InsertUserInfo(githubID, qyvxID, githubName, name); err != nil {
			return err
		}
	}
	// 2 更新别名
	if err := UpdateAlias(param); err != nil {
		return err
	}
	// 3 发送邀请
	if err := invite(teamName, githubName); err != nil {
		return err
	}
	return nil
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
	githubUrl := githubUrlInfo["value"].(map[string]interface{})["text"].(string)
	githubTeamInfo := contentList[1].(map[string]interface{})
	githubTeamName := githubTeamInfo["value"].(map[string]interface{})["selector"].(map[string]interface{})["options"].([]interface{})[0].(map[string]interface{})["value"].([]interface{})[0].(map[string]interface{})["text"].(string)
	param := &models.InviteParam{
		GithubUrl: githubUrl,
		QyvxID:    qyvxID,
		TeamName:  githubTeamName,
	}

	return true, param, nil
}

func UpdateUsers() error {
	// 1.调用接口，获取用户数据
	qyvxUserIds, err := getQyvxUsersIDs()
	if err != nil {
		return err
	}
	zap.L().Info("qyvxIDs", zap.Any("ids", qyvxUserIds))
	// 2.出现新用户进行添加/消失用户删除(更改status)
	// 2.1 数据库中取用户数据
	moUserIds, err := mysql.GetQyvxIDs()
	if err != nil {
		return err
	}
	zap.L().Info("moIDs", zap.Any("ids", moUserIds))

	// 如果数据库中没有数据，则直接添加
	if len(moUserIds) == 0 {
		for _, qid := range qyvxUserIds {
			// 根据qyvx_id获取name
			name, err := getqyvxName(qid)
			if err != nil {
				return err
			}
			if err := mysql.InsertUserInfo("0", qid, "", name); err != nil {
				return err
			}
		}
		return nil
	}

	// 2.2 组建两个map表
	var qyvxMap, moMap map[string]bool
	for _, v := range qyvxUserIds {
		qyvxMap[v] = true
	}
	for _, v := range moUserIds {
		moMap[v] = true
	}

	// 2.3 对比两个map表，进行更新
	//	   -- 在qyvxMap不在moMap 中，则添加
	//	   -- 不再qyvxMap在moMap 中，则删除(status=0)
	//	   -- 在qyvxMap中，在moMap中，但status=0更新0-->1
	for _, qid := range qyvxUserIds {
		if _, ok := moMap[qid]; !ok {
			// 在qyvxMap不在moMap 中，则添加
			// 根据qyvx_id获取name
			name, err := getqyvxName(qid)
			if err != nil {
				return err
			}
			if err := mysql.InsertUserInfo("0", qid, "", name); err != nil {
				return err
			}
		} else {
			// 在qyvxMap中，在moMap中，但status=0更新0-->1
			if err := mysql.UpdateStatusByQyvxID(qid); err != nil {
				return err
			}
		}
	}

	// 遍历moMap，查找不在qyvxMap中的用户并删除
	for u := range moMap {
		if _, ok := qyvxMap[u]; !ok {
			// 不在qyvxMap在moMap 中，则删除(status=0)
			if err := mysql.DeleteUserByQyvxID(u); err != nil {
				return err
			}
		}
	}

	// 3.返回结果
	return nil
}
