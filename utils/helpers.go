package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	data "github.com/dubass83/go-concurrency-project/data/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

const alphabet = "qazwsxedcrfvtgbyhnujmikolp"

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// RandomInt return random int64 between min and max values
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString return random string of given length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for range n {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generate a random Owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generate a random amount of Money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@example.com", RandomString(7))
}

// Generate random user with password
func RandomUser(password string) (user data.User) {

	hash, _ := HashPassword(password)

	return data.User{
		ID: int32(RandomInt(1, 100)),
		Email: pgtype.Text{
			String: RandomEmail(),
			Valid:  true,
		},
		FirstName: pgtype.Text{
			String: RandomString(7),
			Valid:  true,
		},
		LastName: pgtype.Text{
			String: RandomString(8),
			Valid:  true,
		},
		Password: pgtype.Text{
			String: hash,
			Valid:  true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	}
}

// Generate random plan
func RandomPlan() data.Plan {
	return data.Plan{
		ID: int32(RandomInt(1, 100)),
		PlanName: pgtype.Text{
			String: RandomString(16),
		},
		PlanAmount: pgtype.Int4{
			Int32: int32(RandomMoney()),
			Valid: true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	}
}

func TestUserPlan(userID, planID int32) data.UserPlan {
	return data.UserPlan{
		ID: int32(RandomInt(1, 100)),
		UserID: pgtype.Int4{
			Int32: userID,
			Valid: true,
		},
		PlanID: pgtype.Int4{
			Int32: planID,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	}
}
