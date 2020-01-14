package server

import (
	"github.com/pepeunlimited/authentication/internal/app/app1/validator"
	"github.com/pepeunlimited/microservice-kit/headers"
	"github.com/pepeunlimited/microservice-kit/httpz"
	"github.com/pepeunlimited/microservice-kit/jwt"
	"net/http"
	"time"
)

/*
 * READ MORE:
 * http://www.svlada.com/jwt-token-authentication-with-spring-boot/
 * https://jwt.io/introduction/
 * https://blog.questionable.services/article/testing-http-handlers-go/
 */

const (
	VerifyPath  = "/verify"
	SignInPath  = "/sign-in"
	RefreshPath = "/refresh"
)

type Authentication struct {
	validator validator.AuthenticationServerValidator
	jwt jwt.JWT
}

type Auth struct {
	Token string `json:"token"`
}

// write the authorization logic
func (server Authentication) isAuthValid(username string, password string) error {
	return nil
}

func (server Authentication) SignIn() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// is the request even valid?
		username, password, err := server.validator.SignIn(r)
		if err != nil {
			httpz.WriteError(w, err)
			return
		}

		// validate from external service or what ever..
		err = server.isAuthValid(*username, *password)
		if err != nil {
			httpz.WriteError(w, err)
			return
		}

		// generate token
		token, err := server.jwt.SignIn(30*time.Minute, *username, nil, nil, nil)
		if err != nil {
			httpz.WriteError(w, err)
			return
		}
		httpz.WriteOk(w, Auth{Token: string(token)})
	})
}

func (server Authentication) Refresh() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO...
	})
}

func (server Authentication) Verify() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := server.validator.Verify(r)
		if err != nil {
			httpz.WriteError(w, err)
			return
		}
		//TODO:
		// - blacklist

		// add the headers for the microservices..
		w.Header().Add(headers.XJwtUsername, claims.Username)
		//w.Header().Add(headers.XJwtEmail, *claims.Email)
		//w.Header().Add(headers.XJwtRoles, *claims.Role)
		//w.Header().Add(headers.XJwtUserId, strconv.FormatInt(*claims.UserId, 10))
		w.WriteHeader(http.StatusOK)

	})
}


func (server Authentication) NotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

	})
}

func NewAuthenticationServer(secret []byte) Authentication {
	jwt := jwt.NewJWT(secret)
	return Authentication{validator: validator.NewAuthenticationServerValidator(jwt),
						 jwt: jwt}
}