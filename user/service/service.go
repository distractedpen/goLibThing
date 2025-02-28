package service

import (
	"context"
	"gamelib.cloud/user/models"
	"github.com/jackc/pgx/v5"
	"log"
)

type Service struct {
	Db *pgx.Conn
}

func (s *Service) GetUserByIdService(ctx context.Context, id int64) (models.User, error) {
	var user models.User

	result := s.Db.QueryRow(ctx, "SELECT id, name, developer FROM user WHERE id=$1", id)
    err := result.Scan(&user)
    if nil != err {
        log.Printf("%s\n", err.Error())
        return user, err
    }
    return user, nil
}

func (s *Service) GetUserByName(ctx context.Context, name string) (models.User, error) {
	var user models.User

	result := s.Db.QueryRow(ctx, "SELECT id, name FROM user WHERE name=$1", name)
    err := result.Scan(&user)
    if nil != err {
        log.Printf("%s\n", err.Error())
        return user, err
    }
    return user, nil
}

func (s *Service) AddUserService(ctx context.Context, data models.NewUserData) (models.User, error) {
	var user models.User
	_, err := s.Db.Exec(ctx,
		"INSERT INTO users(name) VALUES ($1)",
		data.Name)
	if nil != err {
		log.Printf("%s\n", err.Error())
		return user, err
	}

    user, err = s.GetUserByName(ctx, data.Name)
	if nil != err {
		log.Printf("%s\n", err.Error())
		return user, err
	}
	return user, nil
}


func (s *Service) UpdateUserService(ctx context.Context, id int64, userData models.UpdateUserData) (models.User, error) {
    var user models.User
    _, err := s.Db.Exec(ctx, "UPDATE users SET name=$1 WHERE id=$2", userData.Name, id)
    if nil != err {
        log.Printf("%s\n", err.Error())
        return user, err
    }
    user, err = s.GetUserByIdService(ctx, id)
    if nil != err {
        log.Printf("%s\n", err.Error())
        return user, err
    }
    return user, nil
}


func (s *Service) DeleteUserService(ctx context.Context, id int64) (models.User, error) {
    var user models.User
    user, err := s.GetUserByIdService(ctx, id)
    if nil != err {
        log.Printf("%s\n", err.Error())
        return user, err
    }

    _, err = s.Db.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
    if nil != err { 
        log.Printf("%s\n", err.Error())
        return user, err
    }
    return user, nil
}
