package model

import "time"

type Session struct {
	UUID      string
	UserUUID  string
	CreatedAt time.Time
	ExpiresAt time.Time
}
