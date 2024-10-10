package logic

import (
	"go.uber.org/zap"
	"qyvx/dao/mysql"
	"strings"
	"time"
)

func UpdateOncallers(whitelist string) {
	// 1.获取当前日期
	currentDate := time.Now()
	zap.L().Info("Start a scheduled task about shift")
	// 2.判断当前日期是否为更换值班日
	need, err := mysql.IsUpdateOncallDay(currentDate)
	if err != nil {
		zap.L().Error("mysql.IsUpdateOncallDay failed", zap.Error(err))
		return
	}
	if !need {
		zap.L().Info("today is not updateOncallDay", zap.Any("date", currentDate))
		return
	}

	// 3.获取当前轮次值班人员以及上一轮值班人员
	newOncallers, oldOncallers, err := mysql.GetOncallers(currentDate)
	if err != nil {
		zap.L().Info("get oncallers filed", zap.Any("err", err))
		return
	}
	if newOncallers == "" {
		zap.L().Warn("newOncaller is empty")
	}
	if oldOncallers == "" {
		zap.L().Warn("oldOncaller is empty")
	}

	newSlice := strings.Split(newOncallers, ",")
	oldSlice := strings.Split(oldOncallers, ",")

	// 加载白名单(防止白名单成员被踢出)
	whitelistSlice := strings.Split(whitelist, ",")
	zap.L().Info("whitelist info", zap.Any("whitelist", whitelistSlice))

	whitelistMap := make(map[string]bool, 0)
	for _, name := range whitelistSlice {
		whitelistMap[name] = true
	}

	var addOncallers []string
	var deleteOncallers []string
	for _, nName := range newSlice {
		if !whitelistMap[nName] {
			addOncallers = append(addOncallers, nName)
		}
	}
	for _, oName := range oldSlice {
		if !whitelistMap[oName] {
			deleteOncallers = append(deleteOncallers, oName)
		}
	}

	// 获取qyvxIDs
	var addOncallersIDs []string
	var deleteOncallersIDs []string
	for _, name := range addOncallers {
		qid, err := mysql.GetQyvxIDByName(name)
		if err != nil {
			zap.L().Error("addOncallers mysql.GetQyvxIDByName failed", zap.Any("name", name), zap.Error(err))
			continue
		}
		addOncallersIDs = append(addOncallersIDs, qid)
	}
	for _, name := range deleteOncallers {
		qid, err := mysql.GetQyvxIDByName(name)
		if err != nil {
			zap.L().Error("deleteOncallers mysql.GetQyvxIDByName failed", zap.Error(err))
			continue
		}
		deleteOncallersIDs = append(deleteOncallersIDs, qid)
	}

	// 4.企业微信 拉人/踢人
	if err := inviteAndDeleteQyvxUser(addOncallersIDs, deleteOncallersIDs); err != nil {
		zap.L().Error("InviteAndDeleteQyvxUser failed", zap.Error(err))
		return
	}
	zap.L().Info("Successfully completed the shift change task")
}
