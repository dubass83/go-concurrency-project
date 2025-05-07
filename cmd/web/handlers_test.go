package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	data "github.com/dubass83/go-concurrency-project/data/sqlc"
	"github.com/stretchr/testify/require"
)

var pageTests = []struct {
	name               string
	url                string
	expectedStatusCode int
	handler            http.HandlerFunc
	sessionData        map[string]any
	expectedHTML       string
}{
	{
		name:               "HomePage",
		url:                "/",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.HomePage,
	},
	{
		name:               "login",
		url:                "/login",
		expectedStatusCode: http.StatusOK,
		handler:            testApp.LoginPage,
		expectedHTML:       `<h1 class="mt-5">Login</h1>`,
	},
	{
		name:               "logout",
		url:                "/logout",
		expectedStatusCode: http.StatusSeeOther,
		handler:            testApp.Logout,
		sessionData: map[string]any{
			"userID": 1,
			"user":   data.User{},
		},
	},
}

func TestHandler(t *testing.T) {

	for _, pt := range pageTests {

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
