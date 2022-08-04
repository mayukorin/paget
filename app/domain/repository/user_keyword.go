package repository

import (
	"database/sql"

	"github.com/mayukorin/paget/app/domain/model"
	"github.com/pkg/errors"
)

func InsertUserKeyword(db *sql.DB, user_keyword *model.UserKeyword) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO user_keyword(slack_user_id, keyword_id, user_slack_id) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, errors.Wrap(err, "failed to prepare statement")
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	var id int64
	err = stmt.QueryRow(user_keyword.SlackUserID, user_keyword.KeywordID, user_keyword.UserSlackID).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to insert user_keyword")
	}
	return id, nil
}

func DeleteUserKeyword(db *sql.DB, keywordId int64, userSlackId string) error {
	stmt, err := db.Prepare("DELETE FROM user_keyword WHERE keyword_id = $1 and user_slack_id = $2")
	if err != nil {
		return errors.Wrap(err, "failed to prepare statement")
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	_, err = stmt.Exec(keywordId, userSlackId)
	if err != nil {
		return errors.Wrap(err, "failed to delete user_keyword")
	}
	return nil
}
