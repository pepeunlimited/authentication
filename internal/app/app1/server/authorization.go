package server

import (
	"github.com/pepeunlimited/authorization/internal/app/app1/validator"
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

type Authorization struct {
	validator validator.AuthorizationServerValidator
	jwt jwt.JWT
}

type Auth struct {
	Token string `json:"token"`
}

// write the authorization logic
func (server Authorization) isAuthValid(username string, password string) error {
	return nil
}

func (server Authorization) SignIn() http.Handler {
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

func (server Authorization) Refresh() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO...
	})
}

func (server Authorization) Verify() http.Handler {
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
		//w.Header().Add(headers.XJwtRole, *claims.Role)
		//w.Header().Add(headers.XJwtUserId, strconv.FormatInt(*claims.UserId, 10))
		w.WriteHeader(http.StatusOK)

	})
}


func (server Authorization) NotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

	})
}

func NewAuthorizationServer(secret []byte) Authorization {
	jwt := jwt.NewJWT(secret)
	return Authorization{validator: validator.NewAuthorizationServerValidator(jwt),
						 jwt: jwt}
}