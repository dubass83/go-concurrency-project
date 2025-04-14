package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func getCtx(req *http.Request) context.Context {
	ctx, err := testApp.Session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Error().Err(err)
		return context.Background()
	}

	return ctx
}

func TestAddDefaultData(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	testApp.Session.Put(ctx, "flash", "flash")
	testApp.Session.Put(ctx, "warning", "warning")
	testApp.Session.Put(ctx, "error", "error")

	td := testApp.AddDefaultData(&TemplateData{}, req)

	require.Equal(t, "flash", td.Flash)
	require.Equal(t, "warning", td.Warning)
	require.Equal(t, "error", td.Error)
}

func TestIsAuthenticated(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	auth := testApp.IsAuthenticated(req)
	require.False(t, auth, "returns true for non authenticated")

	testApp.Session.Put(ctx, "userID", 1)
	auth = testApp.IsAuthenticated(req)
	require.True(t, auth, "returns false for authenticated")

}

func TestRenderTemplate(t *testing.T) {
	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	testApp.render(rr, req, "home.page.gohtml", &TemplateData{})

	require.Equal(t, 200, rr.Code, "status code should be 200")

}
