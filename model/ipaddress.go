package model

import (
	"blog-backend/logger"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (repo *Repository) CheckIPAddressAndReturnUserID(ctx context.Context, ip string) (uuid.UUID, error) {
	logger.Println("ip address is", ip)
	user, err := repo.GetUserByIpAddress(ctx, ip)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return uuid.Nil, err
	}
	if err == nil {
		return user.ID, nil
	}
	newUser := User{
		ID:        uuid.New(),
		Email:     sql.NullString{String: "", Valid: false},
		IpAddress: sql.NullString{String: ip, Valid: true},
		Username:  sql.NullString{String: "", Valid: false},
		CreatedAt: time.Now(),
		IsAdmin:   false,
	}
	err = repo.CreateUser(ctx, newUser)
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

func (repo *Repository) CheckIPAddressAndReturnUserIDWithUserName(ctx context.Context, ip, username string) (uuid.UUID, error) {
	logger.Println("ip address is", ip)
	user, err := repo.GetUserByIpAddress(ctx, ip)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Println("error getting user by ip address", err)
		return uuid.Nil, err
	}
	if err == nil {
		if user.Username.String != username {
			user.Username.String = username
			err := repo.UpdateUser(ctx, user)
			if err != nil {
				logger.Println("error updating user", err)
				return uuid.Nil, err
			}
		}
		return user.ID, nil
	}
	newUser := User{
		ID:        uuid.New(),
		Email:     sql.NullString{String: "", Valid: false},
		IpAddress: sql.NullString{String: ip, Valid: true},
		Username:  sql.NullString{String: username, Valid: true},
		CreatedAt: time.Now(),
		IsAdmin:   false,
	}
	err = repo.CreateUser(ctx, newUser)
	if err != nil {
		logger.Println("error creating user", err)
		return uuid.Nil, err
	}
	return user.ID, nil
}

func (repo *Repository) CheckIPAddressAndReturnUserIDWithUserNameAsAdmin(ctx context.Context, ip, username string, alreadyCheckedNoAdmin bool) (uuid.UUID, error) {
	if !alreadyCheckedNoAdmin {
		users, err := repo.GetAdminUsers(ctx)
		if err != nil {
			logger.Println("error getting admin users", err)
			return uuid.Nil, err
		}
		if len(users) > 0 {
			return uuid.Nil, errors.New("admin user already exists")
		}
		if len(users) == 0 {
			return repo.CheckIPAddressAndReturnUserIDWithUserNameAsAdmin(ctx, ip, username, true)
		}
		logger.Fatal("unexpected error")
	}
	logger.Println("ip address is", ip)
	user, err := repo.GetUserByIpAddress(ctx, ip)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logger.Println("error getting user by ip address", err)
		return uuid.Nil, err
	}
	if err == nil {
		// ipアドレスが同じユーザーが存在するけど、adminユーザーでない場合
		if user.Username.String != username {
			user.Username.String = username
			user.IsAdmin = true
			err := repo.UpdateUser(ctx, user)
			if err != nil {
				logger.Println("error updating user", err)
				return uuid.Nil, err
			}
		}
		return user.ID, nil
	}
	// ipアドレスが同じユーザーも存在しない場合
	newUser := User{
		ID:        uuid.New(),
		Email:     sql.NullString{String: "", Valid: false},
		IpAddress: sql.NullString{String: ip, Valid: true},
		Username:  sql.NullString{String: username, Valid: true},
		CreatedAt: time.Now(),
		IsAdmin:   true,
	}
	err = repo.CreateUser(ctx, newUser)
	if err != nil {
		logger.Println("error creating user", err)
		return uuid.Nil, err
	}
	return user.ID, nil
}
