package service

import (
	"context"
	"gamelib.cloud/game/models"
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

func (s *Service) GetGameByIdService(ctx context.Context, id int64) (models.Game, error) {
	var game models.Game

	result := s.Db.QueryRow(ctx, "SELECT id, name, developer FROM games WHERE id=$1", id)
    err := result.Scan(&game)
    if nil != err {
        log.Printf("%s\n", err.Error())
        return game, err
    }
    return game, nil
}

func (s *Service) GetGameByNameAndDeveloper(ctx context.Context, name string, developer string) (models.Game, error) {
	var game models.Game

	result := s.Db.QueryRow(ctx, "SELECT id, name, developer FROM games WHERE name=$1 AND developer=$developer", name, developer)
    err := result.Scan(&game)
    if nil != err {
        log.Printf("%s\n", err.Error())
        return game, err
    }
    return game, nil
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

    game, err = s.GetGameByNameAndDeveloper(ctx, data.Name, data.Developer)
	if nil != err {
		log.Printf("%s\n", err.Error())
		return game, err
	}
	return game, nil
}


func (s *Service) UpdateGameService(ctx context.Context, id int64, gameData models.UpdateGameData) (models.Game, error) {
    var game models.Game
    _, err := s.Db.Exec(ctx, "UPDATE games SET name=$1, developer=$2 WHERE id=$3", gameData.Name, gameData.Developer, id)
    if nil != err {
        log.Printf("%s\n", err.Error())
        return game, err
    }
    game, err = s.GetGameByIdService(ctx, id)
    if nil != err {
        log.Printf("%s\n", err.Error())
        return game, err
    }
    return game, nil
}


func (s *Service) DeleteGameService(ctx context.Context, id int64) (models.Game, error) {
    var game models.Game
    game, err := s.GetGameByIdService(ctx, id)
    if nil != err {
        log.Printf("%s\n", err.Error())
        return game, err
    }

    _, err = s.Db.Exec(ctx, "DELETE FROM games WHERE id=$1", id)
    if nil != err { 
        log.Printf("%s\n", err.Error())
        return game, err
    }
    return game, nil
}
