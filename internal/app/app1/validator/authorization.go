package validator

import (
	"github.com/pepeunlimited/microservice-kit/httpz"
	"github.com/pepeunlimited/microservice-kit/jwt"
	"net/http"
)

type AuthorizationServerValidator struct {
	jwt jwt.JWT
}

func NewAuthorizationServerValidator(jwt jwt.JWT) AuthorizationServerValidator {
	return AuthorizationServerValidator{jwt:jwt}
}

func (AuthorizationServerValidator) SignIn(r *http.Request) (*string, *string, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, nil, httpz.NewMsgError("request not include the basic auth", http.StatusBadRequest)
	}
	return &username, &password, nil
}

func (validator AuthorizationServerValidator) Verify(r *http.Request) (*jwt.CustomClaims, error) {
	authorization := r.Header.Get("Authorization")
	bearer, err := jwt.GetBearer(authorization)
	if err != nil {
		return nil, httpz.NewError(err, http.StatusUnauthorized)
	}
	claims, err := validator.jwt.VerifyCustomClaims(bearer)
	if err != nil {
		return nil, httpz.NewError(err, http.StatusUnauthorized)
	}

	return claims, nil
}