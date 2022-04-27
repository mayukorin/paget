package paget

import (
	"database/sql"
	"fmt"
)

func FindUserId(db *sql.DB, slack_user_id string) (userId int64, err error) {
	userId = 0
	if err = db.QueryRow("SELECT id FROM slack_user WHERE slack_id = $1", slack_user_id).Scan(&userId); err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("error when select slack_user:%q\n", err)
			return
		}
		fmt.Printf("slack_user canot found")
		return
	}
	return
}
