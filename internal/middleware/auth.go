package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

const authorizationHeader = "Authorization"

type AuthMiddleware struct {
	signingKey string
}

type tokenClaims struct {
	jwt.StandardClaims
}

func NewAuthMiddleware(signingKey string) *AuthMiddleware {
	/*AtuhMiddleware constructor function.*/
	return &AuthMiddleware{signingKey: signingKey}
}

func (m *AuthMiddleware) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	/*HTTP middleware wrapper function.*/
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var headerParts []string
		var header string = r.Header.Get(authorizationHeader)

		// check auth header
		if header == "" {
			log.Println("AuthMiddleware: Got empty auth header.")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		headerParts = strings.Split(header, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			log.Println("AuthMiddleware: bad auth header string.")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// check empty token
		if headerParts[1] == "" {
			log.Println("AuthMiddleware: empty token string.")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// parse token and check is it valid
		err = m.parseToken(headerParts[1])
		if err != nil {
			log.Printf("m.parseToken failed: %s.", err.Error())
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// token is valid can run next handler
		next(w, r)
	}
}

func (m *AuthMiddleware) parseToken(accessToken string) error {
	/*Perform parsing JWT token.*/
	var err error
	var token *jwt.Token

	token, err = jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// define parsing key function
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(m.signingKey), nil
	})

	if err != nil {
		return err
	}

	// parse token claims to struct
	_, ok := token.Claims.(*tokenClaims)
	if !ok {
		return errors.New("token claims are not of type *tokenClaims")
	}

	return nil
}
