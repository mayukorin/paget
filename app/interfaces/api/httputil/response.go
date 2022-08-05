package httputil

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondJson(w http.ResponseWriter, status int, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr := w.Write([]byte(err.Error()))
		if writeErr != nil {
			log.Print(writeErr)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, writeErr := w.Write(res)
	if writeErr != nil {
		log.Print(writeErr)
	}
}

func RespondErrorJson(w http.ResponseWriter, code int, err error) {
	log.Printf("code=%d, err=%v", code, err)
	if e, ok := err.(*HTTPError); ok {
		RespondJson(w, code, e)
	} else if err != nil {
		he := HTTPError{
			Message: err.Error(),
		}
		RespondJson(w, code, he)
	}
}
