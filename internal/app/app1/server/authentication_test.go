package server

import (
	"encoding/json"
	"github.com/pepeunlimited/microservice-kit/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var secret = "secret"


func TestAuthorization_SignInOk(t *testing.T) {
	username := "user2@gmail.com"
	password := "p2ssw0rd"

	server := NewAuthenticationServer([]byte(secret))

	//request
	req,_ := http.NewRequest(http.MethodGet, SignInPath, nil)
	req.SetBasicAuth(username, password)

	// recorder
	recorder := httptest.NewRecorder()
	server.SignIn().ServeHTTP(recorder, req)

	if recorder.Code != 200 {
		t.Log(recorder.Code)
		t.Log(recorder.Body.String())
		t.FailNow()
	}
}

func TestAuthorization_VerifyOk(t *testing.T) {
	username := "user2@gmail.com"
	password := "p2ssw0rd"

	server := NewAuthenticationServer([]byte(secret))

	//request
	req,_ := http.NewRequest(http.MethodGet, SignInPath, nil)
	req.SetBasicAuth(username, password)

	recorder0 := httptest.NewRecorder()
	server.SignIn().ServeHTTP(recorder0, req)

	var auth Auth
	err := json.Unmarshal(recorder0.Body.Bytes(), &auth)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	recorder1 := httptest.NewRecorder()
	req,_ = http.NewRequest(http.MethodGet, VerifyPath, nil)
	jwt.SetBearer(auth.Token, req)
	server.Verify().ServeHTTP(recorder1, req)

	if recorder1.Code != http.StatusOK {
		t.Log(recorder1.Code)
		t.Log(recorder1.Body.String())
		t.FailNow()
	}
}

func TestAuthorization_Verify401(t *testing.T) {

	server := NewAuthenticationServer([]byte(secret))

	req,_ := http.NewRequest(http.MethodGet, VerifyPath, nil)
	jwt.SetBearer("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIyQGdtYWlsLmNvbSIsImVtYWlsIjpudWxsLCJyb2xlIjpudWxsLCJ1c2VyX2lkIjpudWxsLCJleHAiOjE1NzYxNjY3NDF9.v2wEmP8gdp5-_AYLAVCZaAqiGZmY6rYDQlGQCy5Mb-0", req)

	recorder := httptest.NewRecorder()
	server.Verify().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Log(recorder.Code)
		t.Log(recorder.Body.String())
		t.FailNow()
	}
}