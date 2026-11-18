package model

import "time"

type NotificationMethod struct {
	ProviderName string // Имя провайдера: telegram, email, push
	Target       string // Адрес получателя: email, ID чата и т.д.
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

// UserInfo — краткая информация о пользователе (то, что возвращает Whoami).
type UserInfo struct {
	UUID  string
	Login string
	Email string
}

// RegisterInput — данные для регистрации нового пользователя.
type RegisterInput struct {
	Login               string
	Password            string
	Email               string
	NotificationMethods []NotificationMethod
}
