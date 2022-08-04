package usecase

import (
	"database/sql"
	"fmt"

	"github.com/mayukorin/paget/app/domain/model"
	"github.com/mayukorin/paget/app/domain/repository"
)

type KeywordUseCase struct {
	db *sql.DB
}

func NewKeywordUseCase(db *sql.DB) *KeywordUseCase {
	return &KeywordUseCase{
		db: db,
	}
}

func (k *KeywordUseCase) Show(showKeywordId int64) (*model.Keyword, error) {
	return repository.FindKeyword(k.db, showKeywordId)
}

func (k *KeywordUseCase) Create(content string) (int64, error) {
	newKeyword := &model.Keyword{
		Content: content,
	}

	var createdId int64
	id, err := repository.InsertKeyword(k.db, newKeyword)
	if err != nil {
		return 0, fmt.Errorf("failed keyword insert: %w", err)
	}
	createdId = id
	return createdId, nil
}
