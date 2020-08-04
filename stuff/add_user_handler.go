package stuff

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	socketio "github.com/googollee/go-socket.io"
)

func (server *Server) AddUserHandler(s socketio.Conn, msg string) {
	fmt.Println("AddUserHandler", msg)

	var dat map[string]string
	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)

	m := make(map[string]string)
	defer func() {
		fmt.Println(">>", m, "<<")
		b, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}
		s.Emit("userAdded", string(b))
	}()

	name := dat["name"]

	if _, present := server.present[name]; present {
		m["msg"] = "Username taken"
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": name,
		"nbf":  time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp":  time.Now().Unix() + 36000,
	})
	tokenString, err := token.SignedString(server.Key)

	if err != nil {
		log.Fatal(err)
	}

	m["msg"] = "OK"
	m["name"] = name
	m["userToken"] = tokenString

	s.SetContext(name)

	server.present[name] = true
}
