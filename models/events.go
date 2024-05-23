package models

import (
	"time"
)

type Events struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	UserId       string    `json:"organizer"`
	AudienceList []string  `json:"audience_list"`
}
