package stuff

import (
	"encoding/json"
	"fmt"
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func (server *Server) JoinGameHandler(s socketio.Conn, msg string) {
	fmt.Println("JoinGameHandler:", msg)
	fmt.Println("The player is", s.Context())

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

	// fmt.Println("notice:", msg)
	// s.Emit("reply", "have "+msg)

	// fmt.Println("addUser", msg)

	// var dat map[string]interface{}
	// if err := json.Unmarshal([]byte(msg), &dat); err != nil {
	// 	panic(err)
	// }

	// name := dat["name"].(string)
	// fmt.Println(dat)

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"name": name,
	// 	"nbf":  time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	// 	"exp":  time.Now().Unix() + 36000,
	// })
	// tokenString, err := token.SignedString(server.Key)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// m := map[string]string{
	// 	"msg":       "OK",
	// 	"name":      name,
	// 	"userToken": tokenString,
	// }
	// b, err := json.Marshal(m)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(b))

	// s.Emit("userAdded", string(b))
}
