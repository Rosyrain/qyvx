package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"qyvx/pkg/ihttp"
	"sync"
)

var (
	cropID            = "" // qyvx的企业id
	agentID           = "" // 应用id
	org               = "" // github组织名
	githubSecret      = "" // github应用的secret
	addressBookSecret = "" //通讯录的secret
	githubToken       = "" // github可以操作组织的token
)

var (
	githubAccessTokenLock      sync.Mutex
	addressBookAccessTokenLock sync.Mutex
	qyvxGithubAccessToken      = ""
	qyvxAddressBookAccessToken = ""
)

func Init() bool {
	cropID = os.Getenv("CROP_ID")
	if cropID == "" {
		zap.L().Warn("Warning: CROP_ID environment variable is not set.")
		return false
	}

	agentID = os.Getenv("AGENT_ID")
	if agentID == "" {
		zap.L().Warn("Warning: AGENT_ID environment variable is not set.")
	}

	org = os.Getenv("ORG")
	if org == "" {
		zap.L().Warn("Warning: ORG environment variable is not set.")
		return false
	}

	githubSecret = os.Getenv("GITHUB_SECRET")
	if githubSecret == "" {
		zap.L().Warn("Warning: GITHUB_SECRET environment variable is not set.")
		return false
	}

	addressBookSecret = os.Getenv("ADDRESS_BOOK_SECRET")
	if addressBookSecret == "" {
		zap.L().Warn("Warning: ADDRESS_BOOK_SECRET environment variable is not set.")
		return false
	}

	githubToken = os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		zap.L().Warn("Warning: GITHUB_TOKEN environment variable is not set.")
		return false
	}
	return true
}

func GetQyvxGithubAccessToken() error {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", cropID, githubSecret)
	body, err := ihttp.Request("GET", url, "", "", nil)
	if err != nil {
		return err
	}
	data := make(map[string]interface{}, 0)
	if err = json.Unmarshal(body, &data); err != nil {
		return err
	}
	qyvxGithubAccessToken = data["access_token"].(string)
	return nil
}

func CheckQyvxGithubAccessToken() error {
	checkUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/agent/get?access_token=%s&agentid=%s", qyvxGithubAccessToken, agentID)
	body, err := ihttp.Request("GET", checkUrl, "", "", nil)
	if err != nil {
		return err
	}
	data := make(map[string]interface{}, 0)
	if err = json.Unmarshal(body, &data); err != nil {
		return err
	}
	errCode := data["errcode"]
	if errCode != float64(0) {
		zap.L().Warn("token expired", zap.Any("errcode", errCode))
		if errCode == float64(42001) || errCode == float64(41001) {
			githubAccessTokenLock.Lock()
			defer githubAccessTokenLock.Unlock()
			zap.L().Info(fmt.Sprintf("the qyvxGithubAccessToken expired or empty,errcode:%f", errCode))
			if err := GetQyvxGithubAccessToken(); err != nil {
				return ErrorRefreshToken
			}
			zap.L().Info("Refresh token success: qyvxGithubAccessToken")
			return nil
		}

		return ErrorGithubAccessToken
	}
	return nil
}

func GetQyvxAddressBookAccessToken() error {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", cropID, addressBookSecret)
	body, err := ihttp.Request("GET", url, "", "", nil)
	if err != nil {
		return err
	}
	data := make(map[string]interface{}, 0)
	if err = json.Unmarshal(body, &data); err != nil {
		return err
	}
	qyvxAddressBookAccessToken = data["access_token"].(string)
	return nil
}

func CheckQyvxAddressBookAccessToken() error {
	postBody := map[string]interface{}{
		"cursor": "",
		"limit":  10000,
	}
	jsonData, err := json.Marshal(postBody)
	if err != nil {
		return err
	}
	checkUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/list_id?access_token=%s", qyvxAddressBookAccessToken)
	body, err := ihttp.Request("POST", checkUrl, "application/json", "", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	data := make(map[string]interface{}, 0)
	if err = json.Unmarshal(body, &data); err != nil {
		return err
	}
	errCode := data["errcode"]
	if errCode != float64(0) {
		zap.L().Warn("token expired", zap.Any("errcode", errCode))
		if errCode == float64(42001) || errCode == float64(41001) {
			addressBookAccessTokenLock.Lock()
			defer addressBookAccessTokenLock.Unlock()
			zap.L().Info(fmt.Sprintf("the qyvxAddressBookAccessToken expired or empty,errcode:%f", data["errcode"]))
			if err := GetQyvxAddressBookAccessToken(); err != nil {
				return ErrorRefreshToken
			}
			zap.L().Info("Refresh token success: qyvxAddressBookAccessToken")
			return nil
		}
		return ErrorAddressAccessToken
	}
	return nil
}

// invite 向github账号发送邀请
func invite(githubID int64) error {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/invitations", org)
	data := map[string]interface{}{
		"invitee_id": githubID,
		"role":       "direct_member",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = ihttp.Request("POST", url, "application/json", githubToken, bytes.NewReader(jsonData))
	if err != nil {
		return ErrorInvite
	}
	return nil
}

func updateAlias(qyvxID, githubName string) error {
	if err := CheckQyvxAddressBookAccessToken(); err != nil {
		return err
	}
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/update?access_token=%s", qyvxAddressBookAccessToken)
	data := map[string]interface{}{
		"userid": qyvxID,
		"alias":  githubName,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = ihttp.Request("POST", url, "application/json", "", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	return nil
}

func getQyvxApproveInfo(spNo string) ([]byte, error) {
	if err := CheckQyvxGithubAccessToken(); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/oa/getapprovaldetail?access_token=%s", qyvxGithubAccessToken)
	data := map[string]string{
		"sp_no": spNo,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	body, err := ihttp.Request("POST", url, "application/json", "", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	return body, nil
}
