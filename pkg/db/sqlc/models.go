// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int64          `json:"id"`
	Email          string         `json:"email"`
	Phone          sql.NullString `json:"phone"`
	HashedPassword string         `json:"hashed_password"`
	Name           string         `json:"name"`
	CreatedAt      time.Time      `json:"created_at"`
}
