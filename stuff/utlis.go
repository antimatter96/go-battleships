package stuff

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	socketio "github.com/googollee/go-socket.io"
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
			fmt.Println("Invalid Token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			fmt.Println("Expired Token")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}
	fmt.Println(token.Claims, err)
	return false
}

func (server *Server) findPlayerFor(name string) *socketio.Conn {
	//fmt.Println("finding", name)
	server.queueLock.Lock()
	defer server.queueLock.Unlock()

	if len(server.Queue) == 0 {
		server.Queue = append(server.Queue, name)
		return nil
	}

	for opp := server.Queue[0]; len(server.Queue) != 0; server.Queue = server.Queue[1:] {
		socket, present := server.socketOf[opp]
		if !present {
			continue
		}
		server.Queue = server.Queue[1:]
		return socket
	}

	return nil
}

type gameResponse struct {
	otherPlayerRes []response
	thisPlayerRes  []response
}

type response struct {
	message string
	data    map[string]string
}

type commonData struct {
	Player    string
	UserToken string
	GameToken string
	GameID    string
}
