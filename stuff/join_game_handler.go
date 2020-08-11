package stuff

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	socketio "github.com/googollee/go-socket.io"
)

func (server *Server) joinGameHandler(s socketio.Conn, msg string) {
	//fmt.Println("JoinGameHandler:", msg)

	var dat map[string]string
	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		panic(err)
	}
	//fmt.Println(dat)

	if !server.verify(dat["player"], dat["userToken"]) {
		return
	}

	name := dat["player"]
	//fmt.Println("The player is", name, s.Context(), name == s.Context())

	otherPlayer := server.findPlayerFor(name)
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

	m1 := map[string]string{
		"gameToken":   tokenString,
		"gameId":      game.ID,
		"otherPlayer": (*otherPlayer).Context().(string),
	}

	m2 := map[string]string{
		"gameToken":   tokenString,
		"gameId":      game.ID,
		"otherPlayer": name,
	}

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
