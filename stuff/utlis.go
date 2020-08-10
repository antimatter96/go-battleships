package stuff

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type commonClaims struct {
	ID  string `json:"id"`
	Exp int64  `json:"exp"`
	Nbf int64  `json:"nbf"`
	jwt.StandardClaims
}

func (server *Server) verify(claim, encryptedToken string) bool {
	token, err := jwt.ParseWithClaims(encryptedToken, &commonClaims{}, func(token *jwt.Token) (interface{}, error) {
		return server.Key, nil
	})

	if claims, ok := token.Claims.(*commonClaims); ok && token.Valid {
		if claims.ID == claim {
			return true
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}
	fmt.Println(token.Claims, err)
	return false
}
