package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // 导入MySQL驱动
	"time"
)

func IsUpdateOncallDay(cTime time.Time) (bool, error) {
	var count int
	sqlStr := `select count(oncallers) from shifts where start_date=?`
	err := db.Get(&count, sqlStr, cTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func GetOncallers(cTime time.Time) (string, string, error) {
	newSqlStr := `select oncallers from shifts where start_date=?`
	oldSqlStr := `select oncallers from shifts where start_date<? order by start_date desc limit 1`

	var newOncallers string
	var oldOncallers string
	err := db.Get(&newOncallers, newSqlStr, cTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", nil
		}
		return "", "", err
	}
	err = db.Get(&oldOncallers, oldSqlStr, cTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return newOncallers, "", nil
		}
		return "", "", err
	}
	return newOncallers, oldOncallers, nil
}
