package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID      `db:"id"`
	Email        sql.NullString `db:"email"`
	IpAddress    sql.NullString `db:"ipaddress"`
	Username     sql.NullString `db:"username"`
	PasswordHash sql.NullString `db:"password_hash"`
	CreatedAt    time.Time      `db:"created_at"`
	IsAdmin      bool           `db:"is_admin"`
}

func (repo *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	var user User
	err := repo.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", id)
	return user, err
}

func (repo *Repository) GetAdminUsers(ctx context.Context) ([]User, error) {
	var users []User
	err := repo.db.SelectContext(ctx, &users, "SELECT * FROM users WHERE is_admin = true")
	return users, err
}

func (repo *Repository) GetUsers(ctx context.Context, limitNumber *int) ([]User, error) {
	var users []User
	if limitNumber != nil {
		err := repo.db.SelectContext(ctx, &users, "SELECT * FROM users LIMIT ?", limitNumber)
		return users, err
	} else {
		err := repo.db.SelectContext(ctx, &users, "SELECT * FROM users")
		return users, err
	}
}

func (repo *Repository) CreateUser(ctx context.Context, user User) error {
	_, err := repo.db.NamedExecContext(ctx, "INSERT INTO users (id, email, ipaddress, username, password_hash, created_at) VALUES (:id, :email, :ipaddress, :username, :password_hash, :created_at)", user)
	return err
}

func (repo *Repository) UpdateUser(ctx context.Context, user User) error {
	_, err := repo.db.NamedExecContext(ctx, "UPDATE users SET email = :email, ipaddress = :ipaddress, username = :username, password_hash = :password_hash, created_at = :created_at, is_admin = :is_admin WHERE id = :id", user)
	return err
}

func (repo *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	return err
}

func (repo *Repository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := repo.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = ?", email)
	return user, err
}

func (repo *Repository) GetUserByIpAddress(ctx context.Context, ipaddress string) (User, error) {
	var user User
	err := repo.db.GetContext(ctx, &user, "SELECT * FROM users WHERE ipaddress = ?", ipaddress)
	return user, err
}

func (repo *Repository) GetUserByUsername(ctx context.Context, username string) (User, error) {
	var user User
	err := repo.db.GetContext(ctx, &user, "SELECT * FROM users WHERE username = ?", username)
	return user, err
}

func (repo *Repository) GetUserByEmailAndPassword(ctx context.Context, email string, password string) (User, error) {
	var user User
	err := repo.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = ? AND password_hash = ?", email, password)
	return user, err
}

func (repo *Repository) GetUserByUsernameAndPassword(ctx context.Context, username string, password string) (User, error) {
	var user User
	err := repo.db.GetContext(ctx, &user, "SELECT * FROM users WHERE username = ? AND password_hash = ?", username, password)
	return user, err
}

func (repo *Repository) GetUserNameById(ctx context.Context, id uuid.UUID) (User, error) {
	var user User
	err := repo.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", id)
	return user, err
}
