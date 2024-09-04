package models

import "encoding/xml"

type MsgContent struct {
	XMLName      xml.Name `xml:"xml"`
	ToUsername   string   `xml:"ToUserName"`
	FromUsername string   `xml:"FromUserName"`
	CreateTime   uint32   `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	Msgid        string   `xml:"MsgId"`
	Agentid      uint32   `xml:"AgentId"`
	Event        string   `xml:"Event"`
	ApprovalInfo struct {
		SpNo       string `xml:"SpNo"`
		SpName     string `xml:"SpName"`
		SpStatus   int    `xml:"SpStatus"`
		TemplateId string `xml:"TemplateId"`
		ApplyTime  uint32 `xml:"ApplyTime"`
		Applyer    struct {
			UserId string `xml:"UserId"`
			Party  string `xml:"Party"`
		} `xml:"Applyer"`
		SpRecord struct {
			SpStatus     int `xml:"SpStatus"`
			ApproverAttr int `xml:"ApproverAttr"`
			Details      struct {
				Approver struct {
					UserId string `xml:"UserId"`
				} `xml:"Approver"`
				Speech   string `xml:"Speech"`
				SpStatus int    `xml:"SpStatus"`
				SpTime   int    `xml:"SpTime"`
			} `xml:"Details"`
		} `xml:"SpRecord"`
		Notifyer struct {
			UserId string `xml:"UserId"`
		} `xml:"Notifyer"`
		StatuChangeEvent int `xml:"StatuChangeEvent"`
		ProcessList      struct {
			NodeList []struct {
				NodeType    int `xml:"NodeType"`
				SpStatus    int `xml:"SpStatus"`
				ApvRel      int `xml:"ApvRel"`
				SubNodeList struct {
					UserInfo struct {
						UserId string `xml:"UserId"`
					} `xml:"UserInfo"`
					Speech string `xml:"Speech"`
					SpYj   int    `xml:"SpYj"`
					Sptime int    `xml:"Sptime"`
				} `xml:"SubNodeList"`
			} `xml:"NodeList"`
		} `xml:"ProcessList"`
	} `xml:"ApprovalInfo"`
}
