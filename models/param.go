package models

type InviteParam struct {
	QyvxID    string `json:"qyvx_id" binding:"required"`
	GithubUrl string `json:"github_url" binding:"required"`
	TeamName  string `json:"team_name" binding:"required"`
}
