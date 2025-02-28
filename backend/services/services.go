package services

import (
	"context"
	"gamelib.cloud/models"
	"github.com/jackc/pgx/v5"
	"log"
)

type Service struct {
	Db *pgx.Conn
}

func (s *Service) GetGamesService(ctx context.Context) ([]models.Game, error) {
	gamesResults := make([]models.Game, 0)
	// get the results for the database
	result, err := s.Db.Query(ctx, "SELECT id, name, developer FROM games;")
	if nil != err {
		return gamesResults, err
	}
	// parse the db results into objects
	for result.Next() {
		var game models.Game
		result.Scan(&game.Id, &game.Name, &game.Developer)
		gamesResults = append(gamesResults, game)
	}
	err = result.Err()
	if nil != err {
		return gamesResults, err
	}

	return gamesResults, nil
}

func (s *Service) AddGameService(ctx context.Context, data models.NewGameData) (models.Game, error) {
	var game models.Game
	_, err := s.Db.Exec(ctx,
		"INSERT INTO Games(name, developer) VALUES ($1, $2)",
		data.Name, data.Developer)
	if nil != err {
		log.Printf("%s\n", err.Error())
		return game, err
	}
	result := s.Db.QueryRow(ctx,
		"SELECT id, name, developer FROM Games WHERE name=$1 AND developer=$2",
		data.Name, data.Developer)
	err = result.Scan(&game.Id, &game.Name, &game.Developer)
	if nil != err {
		log.Printf("%s\n", err.Error())
		return game, err
	}
	return game, nil
}
