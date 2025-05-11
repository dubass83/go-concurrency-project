package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

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

	pass := "Qw12345678!"
	user := utils.RandomUser(pass)
	plan := utils.RandomPlan()
	userPlan := utils.TestUserPlan(user.ID, plan.ID)

	pagePostTests := []struct {
		name               string
		url                string
		expectedStatusCode int
		handler            http.HandlerFunc
		postedData         url.Values
		sessionData        map[string]any
		expectedHTML       string
		expectedSessionKey string
		buildStubs         func(store *mockdb.MockStore)
	}{
		{
			name:               "loginPage",
			url:                "/login",
			expectedStatusCode: http.StatusSeeOther,
			expectedSessionKey: "userID",
			handler:            testApp.PostLoginPage,
			postedData: url.Values{
				"email":    {user.Email.String},
				"password": {pass},
			},
			buildStubs: func(store *mockdb.MockStore) {
				argUser := pgtype.Text{
					String: user.Email.String,
					Valid:  true,
				}
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(argUser)).
					Times(1).
					Return(user, nil)

				argPlan := pgtype.Int4{
					Int32: user.ID,
					Valid: true,
				}
				store.EXPECT().
					GetOneUserPlan(gomock.Any(), gomock.Eq(argPlan)).
					Times(1).
					Return(userPlan, nil)
			},
		},
	}

	for _, pt := range pagePostTests {

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", pt.url, strings.NewReader(pt.postedData.Encode()))
		// Add this line to set the correct Content-Type header for form data
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
		if len(pt.expectedSessionKey) > 0 {
			require.True(t, testApp.Session.Exists(ctx, pt.expectedSessionKey), fmt.Sprintf("test name: %s", pt.name))
		}
	}
}

func TestGetHandlerSubscribeToPlan(t *testing.T) {

	pass := "Qw12345678!"
	user := utils.RandomUser(pass)
	plan := utils.RandomPlan()
	userPlan := utils.TestUserPlan(user.ID, plan.ID)

	pageGetTests := []struct {
		name               string
		url                string
		expectedStatusCode int
		handler            http.HandlerFunc
		sessionData        map[string]any
		expectedHTML       string
		expectedSessionKey string
		buildStubs         func(store *mockdb.MockStore)
	}{
		{
			name:    "SubscribeToPlan",
			url:     fmt.Sprintf("/subscribe?id=%d", user.ID),
			handler: testApp.SubscribeToPlan,
			sessionData: map[string]any{
				"userID": user.ID,
				"user":   user,
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedSessionKey: "user-plan",
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetOnePlan(gomock.Any(), gomock.Eq(plan.ID)).
					Times(1).
					Return(plan, nil)

				arg := data.SubscribeUserToPlanParams{
					UserID: user.ID,
					PlanID: plan.ID,
				}
				store.EXPECT().
					SubscribeUserToPlan(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(userPlan, nil)
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
		if len(pt.expectedSessionKey) > 0 {
			require.True(t, testApp.Session.Exists(ctx, pt.expectedSessionKey), fmt.Sprintf("test name: %s", pt.name))
		}
	}
}
