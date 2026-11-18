package model

import "time"

type NotificationMethod struct {
	ProviderName string `json:"provider_name"`
	Target       string `json:"target"`
}

type User struct {
	UUID                string
	Login               string
	Email               string
	PasswordHash        string
	NotificationMethods []NotificationMethod
	CreatedAt           time.Time
	UpdatedAt           *time.Time
}
