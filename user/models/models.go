package models

import (
	"context"

	gameModels "gamelib.cloud/game/models"
	gameService "gamelib.cloud/game/services"
)

type NewUserData struct {
	Name string
}
type UpdateUserData = NewUserData

func (nu *NewUserData) MapToUser(id int64) User {
	return User{
		Id:      id,
		Name:    nu.Name,
		Library: make([]int64, 0),
	}
}

type UpdateUserLibraryData struct {
	Library []int64
}

// DB Data Object
type User struct {
	Id      int64
	Name    string
	Library []int64
}

// Create User Dto Object from User DB object
func (u *User) MapToUserDto(context context.Context, gameService *gameService.Service) UserDto {
	games := make([]gameModels.Game, 0)
	for _, id := range u.Library {
		game, _ := gameService.GetGameByIdService(context, id)
		games = append(games, game)
	}

	return UserDto{
		Id:      u.Id,
		Name:    u.Name,
		Library: games,
	}
}

// DTO sent to FE
type UserDto struct {
	Id      int64
	Name    string
	Library []gameModels.Game
}
