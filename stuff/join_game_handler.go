package stuff

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	socketio "github.com/googollee/go-socket.io"
)

//
type myCustomClaims struct {
	Name string `json:"name"`
	Exp  int64  `json:"exp"`
	Nbf  int64  `json:"nbf"`
	jwt.StandardClaims
}

func (server *Server) vetify(name, userToken string) bool {
	token2, err := jwt.ParseWithClaims(userToken, &myCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return server.Key, nil
	})

	if claims, ok := token2.Claims.(*myCustomClaims); ok && token2.Valid {
		if claims.Name == name {
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
	fmt.Println(token2.Claims, err)
	return false
}

func (server *Server) FindPlayerFor(name string) *socketio.Conn {
	fmt.Println("finding", name)
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

func (server *Server) JoinGameHandler(s socketio.Conn, msg string) {
	fmt.Println("JoinGameHandler:", msg)

	var dat map[string]string
	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		panic(err)
	}
	//fmt.Println(dat)

	if !server.vetify(dat["player"], dat["userToken"]) {
		return
	}

	name := dat["player"]
	//fmt.Println("The player is", name, s.Context(), name == s.Context())

	otherPlayer := server.FindPlayerFor(name)
	if otherPlayer == nil {
		fmt.Println("Not found for ", name)
		return
	}

	game, err := NewBattleShips(name, (*otherPlayer).Context().(string))
	if err != nil {
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  game.ID,
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp": time.Now().Unix() + 36000,
	})
	tokenString, err := token.SignedString(server.Key)
	if err != nil {
		log.Fatal(err)
	}

	m1 := make(map[string]string)
	m2 := make(map[string]string)

	m1["gameToken"] = tokenString
	m2["gameToken"] = tokenString

	m1["gameId"] = game.ID
	m2["gameId"] = game.ID

	m1["otherPlayer"] = (*otherPlayer).Context().(string)
	m2["otherPlayer"] = name

	server.games[game.ID] = game

	//fmt.Println(">>", m1, "<<")
	b1, err := json.Marshal(m1)
	if err != nil {
		log.Fatal(err)
	}
	s.Emit("startGame", string(b1))

	//fmt.Println(">>", m2, "<<")
	b2, err := json.Marshal(m2)
	if err != nil {
		log.Fatal(err)
	}
	(*otherPlayer).Emit("startGame", string(b2))
}
