package repository

import (
	"database/sql"

	"github.com/mayukorin/paget/app/domain/model"
	"github.com/pkg/errors"
)

func InsertSlackUser(db *sql.DB, slack_user *model.SlackUser) (string, error) {

	stmt, err := db.Prepare("INSERT INTO slack_user(slack_id, slack_channel_id) values ($1, $2) RETURNING slack_id")
	if err != nil {
		return "", errors.Wrap(err, "failed to prepare statement")
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	var slackId string
	err = stmt.QueryRow(slack_user.SlackID, slack_user.ChannelID).Scan(&slackId)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute insert slack_user")
	}
	return slackId, nil
}

func FindSlackUser(db *sql.DB, slack_id string) (*model.SlackUser, error) {
	stmt, err := db.Prepare("SELECT id, slack_id, slack_channel_id FROM keyword WHERE slack_id = $1")
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

func UpdateSlackUser(db *sql.DB, slack_user *model.SlackUser) error {
	stmt, err := db.Prepare("UPDATE slack_user SET latest_matched_paper = $1 WHERE slack_id = $2")
	if err != nil {
		return errors.Wrap(err, "failed to set prepared statement")
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	_, err = stmt.Exec(slack_user.LatestMatchedPaper, slack_user.SlackID)
	if err != nil {
		return errors.Wrap(err, "failed to execute update slack user")
	}
	return nil
}
