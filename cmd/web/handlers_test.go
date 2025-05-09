package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	mockdb "github.com/dubass83/go-concurrency-project/data/mock"
	data "github.com/dubass83/go-concurrency-project/data/sqlc"
	"github.com/dubass83/go-concurrency-project/utils"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {

	pageGetTests := []struct {
		name               string
		url                string
		expectedStatusCode int
		handler            http.HandlerFunc
		sessionData        map[string]any
		expectedHTML       string
	}{
		{
			name:               "homePage",
			url:                "/",
			expectedStatusCode: http.StatusOK,
			handler:            testApp.HomePage,
		},
		{
			name:               "loginPage",
			url:                "/login",
			expectedStatusCode: http.StatusOK,
			handler:            testApp.LoginPage,
			expectedHTML:       `<h1 class="mt-5">Login</h1>`,
		},
		{
			name:               "logoutPage",
			url:                "/logout",
			expectedStatusCode: http.StatusSeeOther,
			handler:            testApp.Logout,
			sessionData: map[string]any{
				"userID": 1,
				"user":   data.User{},
			},
		},
	}

	for _, pt := range pageGetTests {

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", pt.url, nil)

		ctx := getCtx(req)
		req = req.WithContext(ctx)

		if len(pt.sessionData) > 0 {
			for key, value := range pt.sessionData {
				testApp.Session.Put(ctx, key, value)
			}
		}

		pt.handler.ServeHTTP(rr, req)

		require.Equal(t, pt.expectedStatusCode, rr.Code, fmt.Sprintf("test name: %s", pt.name))

		if len(pt.expectedHTML) > 0 {
			html := rr.Body.String()
			require.Contains(t, html, pt.expectedHTML, fmt.Sprintf("test name: %s", pt.name))
		}
	}
}

func TestPostHandler(t *testing.T) {

	hash, err := utils.HashPassword("Qw12345678!")
	require.NoError(t, err)

	user := data.User{
		ID: 1,
		Email: pgtype.Text{
			String: "admin@example.com",
			Valid:  true,
		},
		FirstName: pgtype.Text{
			String: "Fake",
			Valid:  true,
		},
		LastName: pgtype.Text{
			String: "Admin",
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

	// plan := data.Plan{
	// 	ID: 1,
	// 	PlanName: pgtype.Text{
	// 		String: "fakePlan",
	// 	},
	// 	PlanAmount: pgtype.Int4{
	// 		Int32: 100,
	// 		Valid: true,
	// 	},
	// 	CreatedAt: pgtype.Timestamp{
	// 		Time:  time.Now(),
	// 		Valid: true,
	// 	},
	// 	UpdatedAt: pgtype.Timestamp{
	// 		Time:  time.Now(),
	// 		Valid: true,
	// 	},
	// }

	userPlan := data.UserPlan{
		ID: 1,
		UserID: pgtype.Int4{
			Int32: 1,
			Valid: true,
		},
		PlanID: pgtype.Int4{
			Int32: 1,
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

	pagePostTests := []struct {
		name               string
		url                string
		expectedStatusCode int
		handler            http.HandlerFunc
		postedData         url.Values
		sessionData        map[string]any
		expectedHTML       string
		buildStubs         func(store *mockdb.MockStore)
	}{
		{
			name:               "loginPage",
			url:                "/login",
			expectedStatusCode: http.StatusSeeOther,
			handler:            testApp.PostLoginPage,
			postedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"Qw12345678!"},
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					GetOneUserPlan(gomock.Any(), gomock.Any()).
					Times(1).
					Return(userPlan, nil)
			},
		},
	}

	for _, pt := range pagePostTests {

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", pt.url, strings.NewReader(pt.postedData.Encode()))

		ctx := getCtx(req)
		req = req.WithContext(ctx)

		if len(pt.sessionData) > 0 {
			for key, value := range pt.sessionData {
				testApp.Session.Put(ctx, key, value)
			}
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		pt.buildStubs(store)

		testApp.Store = store

		pt.handler.ServeHTTP(rr, req)

		require.Equal(t, pt.expectedStatusCode, rr.Code, fmt.Sprintf("test name: %s", pt.name))

		if len(pt.expectedHTML) > 0 {
			html := rr.Body.String()
			require.Contains(t, html, pt.expectedHTML, fmt.Sprintf("test name: %s", pt.name))
		}
	}
}
