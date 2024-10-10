package mysql

func CheckUserExistByQyvxID(qyvxID string) (err error) {
	sqlStr := `select count(name) from user where qyvx_id = ?`

	var count int
	if err = db.Get(&count, sqlStr, qyvxID); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return ErrorUserNotExist
}

func InsertUserInfo(githubID, qyvxID, githubName, qyvxName string) (err error) {
	sqlStr := `Insert into user (github_id,qyvx_id,github_name,name) values(?,?,?,?)`
	_, err = db.Exec(sqlStr, githubID, qyvxID, githubName, qyvxName)
	return err
}

func UpdateGithubIDByQyvxID(qyvxID, githubName, githubID string) (err error) {
	sqlStr := `UPDATE user SET github_id =?, github_name =? WHERE qyvx_id =?`
	_, err = db.Exec(sqlStr, githubID, githubName, qyvxID)
	return
}

func GetGithubNameByQyvxID(qyvxId string) (githubName string, err error) {
	sqlStr := `select github_name from user where qyvx_id = ?`
	err = db.Get(&githubName, sqlStr, qyvxId)
	return
}

func GetQyvxIDs() ([]string, error) {
	var ids []string
	sqlStr := `select qyvx_id from user`
	err := db.Select(&ids, sqlStr)
	if len(ids) == 0 {
		return nil, nil
	}
	return ids, err
}

func UpdateStatusByQyvxID(qyvxIds string) error {
	sqlStr := `update user set status = 1 where qyvx_id=?`
	_, err := db.Exec(sqlStr, qyvxIds)
	return err
}

func DeleteUserByQyvxID(qyvxID string) error {
	sqlStr := `update user set status = 0 where qyvx_id=?`
	_, err := db.Exec(sqlStr, qyvxID)
	return err
}

func GetQyvxIDByName(name string) (string, error) {
	var qyvxId string
	sqlStr := `select qyvx_id from user where name = ?`
	err := db.Get(&qyvxId, sqlStr, name)
	return qyvxId, err
}
