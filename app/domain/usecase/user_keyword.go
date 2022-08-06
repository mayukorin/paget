package usecase

import (
	"database/sql"
	"fmt"

	"github.com/mayukorin/paget/app/domain/model"
	"github.com/mayukorin/paget/app/domain/repository"
)

type UserKeywordUseCase struct {
	db *sql.DB
}

func NewUserKeywordUseCase(db *sql.DB) *UserKeywordUseCase {
	return &UserKeywordUseCase{
		db: db,
	}
}

func (u *UserKeywordUseCase) Create(slackUserId int64, keywordId int64, userSlackId string) (int64, error) {
	newUserKeyword := &model.UserKeyword{
		SlackUserID: slackUserId,
		KeywordID:   keywordId,
		UserSlackID: userSlackId,
	}

	var createdUserKeywordId int64
	id, err := repository.InsertUserKeyword(u.db, newUserKeyword)
	if err != nil {
		return 0, fmt.Errorf("failed user keyword insert: %w", err)
	}
	createdUserKeywordId = id
	return createdUserKeywordId, nil
}

func (u *UserKeywordUseCase) Destroy(keywordId int64, userSlackId string) error {
	err := repository.DeleteUserKeyword(u.db, keywordId, userSlackId)
	if err != nil {
		return fmt.Errorf("failed user keyword delete: %w", err)
	}
	return nil
}
