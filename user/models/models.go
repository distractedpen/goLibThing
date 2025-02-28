package models

import (
	"gamelib.cloud/game/models"
)

type NewUserData struct {
	Name string
}
type UpdateUserData = NewUserData

type User struct {
	Id      int64
	Name    string
	Library []models.Game
}

func (nu *NewUserData) MapToUser(id int64) User {
	return User{
		Id:      id,
		Name:    nu.Name,
		Library: make([]models.Game, 0),
	}
}
