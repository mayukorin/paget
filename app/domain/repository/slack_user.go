package repository

import (
	"database/sql"

	"github.com/mayukorin/paget/app/domain/model"
	"github.com/pkg/errors"
)

func InsertSlackUser(db *sql.DB, slack_user *model.SlackUser) (int64, error) {

	stmt, err := db.Prepare("INSERT INTO slack_user(slack_id, slack_channel_id) values ($1, $2) RETURNING id")
	if err != nil {
		return 0, errors.Wrap(err, "failed to prepare statement")
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	var id int64
	err = stmt.QueryRow(slack_user.SlackID, slack_user.ChannelID).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to execute insert slack_user")
	}
	return id, nil
}

func FindSlackUser(db *sql.DB, slack_id string) (*model.SlackUser, error) {
	stmt, err := db.Prepare("SELECT id, slack_id, slack_channel_id FROM keyword WHERE content = $1")
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare statement")
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	slack_user := &model.SlackUser{}
	err = stmt.QueryRow(slack_id).Scan(&slack_user.ID, &slack_user.SlackID, &slack_user.ChannelID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute select slack_user")
	}
	return slack_user, nil
}
