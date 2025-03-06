package utils

import (
	"context"
	"fmt"
	"time"

	data "github.com/dubass83/go-concurrency-project/data/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generate password hash or return error
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error: can not generate hash from password - %v", err)
	}
	return string(hash), nil
}

// CheckPassword check if provided password correct or not.
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// ResetPassword is the method we will use to change a user's password.
func ResetPassword(password string, user_id int32, store data.Store) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	// Update user password
	_, err = store.UpdateUser(ctx, data.UpdateUserParams{
		Password: pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		},
		ID: user_id,
	})
	if err != nil {
		return err
	}

	return nil
}
