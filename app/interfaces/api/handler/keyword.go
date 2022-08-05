package handler

import (
	"net/http"

	"github.com/mayukorin/paget/app/domain/usecase"
	"github.com/slack-go/slack"
)

type KeywordHandler struct {
	keywordUseCase *usecase.KeywordUseCase
}

func NewKeywordHandler(keywordUseCase *usecase.KeywordUseCase) *KeywordHandler {
	return &KeywordHandler{
		keywordUseCase: keywordUseCase,
	}
}

func (h *KeywordHandler) Index(_ http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	keywords, err := h.keywordUseCase.Index(r.Context()) // slack の User ID が欲しい
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	keywords_message := ""
	for _, k := range keywords {
		keywords_message += k.Content + ", "
	}

	return http.StatusOK, &slack.Msg{Text: "keyword の一覧： " + keywords_message}, nil
}
