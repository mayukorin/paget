package repository

import (
	"database/sql"

	"github.com/mayukorin/paget/app/domain/model"
	"github.com/pkg/errors"
)

func InsertKeyword(db *sql.DB, keyword *model.Keyword) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO keyword(content) VALUES($1) RETURNING id")
	if err != nil {
		return 0, errors.Wrap(err, "failed to prepare statement")
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	var id int64
	err = stmt.QueryRow(keyword.Content).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to execute insert keyword")
	}
	return id, nil
}

func FindKeyword(db *sql.DB, id int64) (*model.Keyword, error) {
	stmt, err := db.Prepare("SELECT id, content FROM keyword WHERE id = $1")
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare statement")
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	keyword := &model.Keyword{}
	err = stmt.QueryRow(id).Scan(&keyword.ID, &keyword.Content)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute select keyword")
	}
	return keyword, nil
}

func AllKeyword(db *sql.DB, userSlackId string) ([]model.Keyword, error) {
	stmt, err := db.Prepare("SELECT content FROM keyword JOIN user_keyword on (keyword.id = user_keyword.keyword_id) WHERE user_keyword.user_slack_id = $1")
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare statment")
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	keywords := []model.Keyword{}
	rows, err := stmt.Query(userSlackId)
	for rows.Next() {
		keyword := model.Keyword{}
		err := rows.Scan(keyword.Content)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan")
		}
		keywords = append(keywords, keyword)
	}
	return keywords, nil
}
