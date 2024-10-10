package utools

import (
	"encoding/json"
	"errors"
	"fmt"
	"qyvx/pkg/ihttp"
	"strconv"
	"strings"
)

var (
	ErrorUrl         = errors.New("wrong url")
	ErrorGithubUser  = errors.New("invalid github account")
	ErrorGetGithubID = errors.New("get githubID failed")
)

// parseGithubUrl 解析提供的github个人主页并获取githubID
func ParseGithubUrl(url string) (string, error) {
	//访问GithubUrl易超时，暂时废弃
	_, err := ihttp.Request("GET", url, "", "", nil)
	if err != nil {
		return "", ErrorUrl
	}
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		githubName := parts[len(parts)-1]
		return githubName, nil
	} else {
		return "", fmt.Errorf("get githubName failed")
	}
}

func GetGithubID(githubName, githubUrl, token string) (githubID string, err error) {
	url := fmt.Sprintf("https://api.github.com/search/users?q=%s", githubName)
	body, err := ihttp.Request("GET", url, "", token, nil)
	if err != nil {
		return "0", err
	}
	data := make(map[string]interface{}, 0)
	if err = json.Unmarshal(body, &data); err != nil {
		return "0", err
	}

	items, ok := data["items"].([]interface{})
	if !ok {
		return "0", err
	}
	if len(items) == 0 {
		return "0", ErrorGithubUser
	}

	// 假设 items 数组中的每个元素都是一个 map
	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return "0", err
		}

		// 1.先判断html_url是否为对应的githubUrl
		htmlUrl, ok := itemMap["html_url"].(string)
		if !ok || htmlUrl != githubUrl {
			continue
		}

		// 获取 id 字段并断言为 float64（JSON 数字被解码为 float64）
		id, ok := itemMap["id"].(float64)
		if !ok {
			return "0", ErrorGetGithubID
		}
		return strconv.Itoa(int(id)), nil
	}
	return "0", err
}
