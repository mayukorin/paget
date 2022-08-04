package usecase

import (
	"database/sql"
	"fmt"

	"github.com/mayukorin/paget/app/domain/model"
	"github.com/mayukorin/paget/app/domain/repository"
)

type SlackUserUseCase struct {
	db *sql.DB
}

func NewSlackUserUseCase(db *sql.DB) *SlackUserUseCase {
	return &SlackUserUseCase{
		db: db,
	}
}

func (s *SlackUserUseCase) Show(showSlackUserId string) (*model.SlackUser, error) {
	return repository.FindSlackUser(s.db, showSlackUserId)
}

func (s *SlackUserUseCase) Create(slackId string, slackChannelId string) (string, error) {
	newSlackUser := &model.SlackUser{
		SlackID:   slackId,
		ChannelID: slackChannelId,
	}

	var createdSlackId string
	slackId, err := repository.InsertSlackUser(s.db, newSlackUser)
	if err != nil {
		return "", fmt.Errorf("failed slack user insert: %w", err)
	}
	createdSlackId = slackId
	return createdSlackId, nil
}
