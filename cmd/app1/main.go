package main

import (
	"github.com/pepeunlimited/authentication/internal/app/app1/server"
	"github.com/pepeunlimited/microservice-kit/jwt"
	"github.com/pepeunlimited/microservice-kit/misc"
	"log"
	"net/http"
)

const (
	Version = "0.1"
)

func main() {
	log.Printf("Starting the AuthenticationServer... version=[%v]", Version)
	secret := misc.GetEnv(jwt.SECRET_KEY, "v3ry-s3cr3t-k3y")
	s := server.NewAuthenticationServer([]byte(secret))
	mux := http.NewServeMux()
	mux.Handle(server.SignInPath, s.SignIn())
	mux.Handle(server.RefreshPath, s.Refresh())
	mux.Handle(server.VerifyPath, s.Verify())
	mux.Handle("/", s.NotFound())

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panic(err)
	}

}