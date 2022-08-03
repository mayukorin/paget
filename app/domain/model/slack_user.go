package model

type SlackUser struct {
	ID                 int64
	SlackID            string
	ChannelID          string
	LatestMatchedPaper string
}
