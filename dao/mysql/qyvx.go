package mysql

func CheckUserExistByQyvxID(qyvxID string) (err error) {
	sqlStr := `select count(user_id) from user where qyvx_id = ?`

	var count int
	if err = db.Get(&count, sqlStr, qyvxID); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return ErrorUserNotExist
}

func InsertUserInfo(userID, githubID int64, qyvxID, githubName string) (err error) {
	sqlStr := `Insert into user (user_id,github_id,qyvx_id,github_name) values(?,?,?,?)`
	_, err = db.Exec(sqlStr, userID, githubID, qyvxID, githubName)
	return err
}

func UpdateGithubIDByQyvxID(qyvxID, githubName string, githubID int64) (err error) {
	sqlStr := `UPDATE user SET github_id =?, github_name =? WHERE qyvx_id =?`
	_, err = db.Exec(sqlStr, githubID, githubName, qyvxID)
	return
}

func GetGithubNameByQyvxID(qyvxId string) (githubName string, err error) {
	sqlStr := `select github_name from user where qyvx_id = ?`
	err = db.Get(&githubName, sqlStr, qyvxId)
	return
}
