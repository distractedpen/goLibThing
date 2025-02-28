package models


// Dto
type NewGameData struct {
	Name      string `json:"name"`
	Developer string `json:"developer"`
}
type UpdateGameData = NewGameData

// Database Object
type Game struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Developer string `json:"developer"`
}

// Mapping Helpers
func (g *NewGameData) MapToGame(id int64) Game {
	return Game{
		Id:        id,
		Name:      g.Name,
		Developer: g.Developer,
	}
}

