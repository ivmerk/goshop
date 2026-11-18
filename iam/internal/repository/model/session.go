package model

import "time"

// Session хранится в Redis в виде JSON под ключом session:<uuid>.
type Session struct {
	UUID      string    `json:"uuid"`
	UserUUID  string    `json:"user_uuid"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
