package utils

import (
	"fmt"

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
// func (u *User) ResetPassword(password string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
// 	defer cancel()

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
// 	if err != nil {
// 		return err
// 	}

// 	stmt := `update users set password = $1 where id = $2`
// 	_, err = db.ExecContext(ctx, stmt, hashedPassword, u.ID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
