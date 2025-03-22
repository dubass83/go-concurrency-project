package main

import (
	"fmt"
	"strings"
	"time"

	goalone "github.com/bwmarrin/go-alone"
)

// GenerateTokenFromString generates a signed token
func (app *Server) GenerateTokenFromString(data string) string {
	secretKey := []byte(app.Config.TokenSecret)
	var urlToSign string

	s := goalone.New(secretKey, goalone.Timestamp)
	if strings.Contains(data, "?") {
		urlToSign = fmt.Sprintf("%s&hash=", data)
	} else {
		urlToSign = fmt.Sprintf("%s?hash=", data)
	}

	tokenBytes := s.Sign([]byte(urlToSign))
	token := string(tokenBytes)

	return token
}

// VerifyToken verifies a signed token
func (app *Server) VerifyToken(token string) bool {
	secretKey := []byte(app.Config.TokenSecret)
	s := goalone.New(secretKey, goalone.Timestamp)
	_, err := s.Unsign([]byte(token))

	if err != nil {
		// signature is not valid. Token was tampered with, forged, or maybe it's
		// not even a token at all! Either way, it's not safe to use it.
		return false
	}
	// valid hash
	return true

}

// Expired checks to see if a token has expired
func (app *Server) TokenExpired(token string, minutesUntilExpire int) bool {
	secretKey := []byte(app.Config.TokenSecret)
	s := goalone.New(secretKey, goalone.Timestamp)
	ts := s.Parse([]byte(token))

	// time.Duration(seconds)*time.Second
	return time.Since(ts.Timestamp) > time.Duration(minutesUntilExpire)*time.Minute
}
