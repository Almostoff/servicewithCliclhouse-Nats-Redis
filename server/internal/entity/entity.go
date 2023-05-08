package entity

import (
	"time"
)

type Item struct {
	ID          int       `json:"id"`
	CampaignID  int       `json:"campaign_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"created_at"`
}

type DelR struct {
	ID         int  `json:"id"`
	CampaignId int  `json:"campaign_id"`
	Removed    bool `json:"removed"`
}

type LogData struct {
	ID          int       `json:"id"`
	CampaignId  int       `json:"campaign_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     uint8     `json:"removed"`
	EventTime   time.Time `json:"created_at"`
}
