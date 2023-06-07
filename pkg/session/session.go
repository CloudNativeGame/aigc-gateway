package session

import (
	"github.com/gin-contrib/sessions"
)

type SessionStorage struct {
	Session sessions.Session
}

func (storage *SessionStorage) GetItem(key string) string {
	value := storage.Session.Get(key)
	if value == nil {
		return ""
	}
	return value.(string)
}

func (storage *SessionStorage) SetItem(key, value string) {
	storage.Session.Set(key, value)
	storage.Session.Save()
}
