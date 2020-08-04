package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	socketio "github.com/googollee/go-socket.io"
)

type MyCustomClaims struct {
	Foo string `json:"foo"`
	ID  string `json:"id"`
	jwt.StandardClaims
}

func main() {

	privateKeyPEM, err := ioutil.ReadFile("./private_key.pem")

	ss := Server{key: privateKeyPEM}
	//_privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	// //fmt.Println(privateKey)

	// tokenString, err := token.SignedString(privateKeyPEM)

	// fmt.Println(">>", tokenString, "<<", err)

	// token2, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return privateKeyPEM, nil
	// })

	// if claims, ok := token2.Claims.(*MyCustomClaims); ok && token2.Valid {
	// 	fmt.Printf("%+v\n", claims)
	// } else {
	// 	fmt.Println(err)
	// }

	// if token2.Valid {
	// 	fmt.Println(token2.Claims)
	// } else if ve, ok := err.(*jwt.ValidationError); ok {
	// 	if ve.Errors&jwt.ValidationErrorMalformed != 0 {
	// 		fmt.Println("That's not even a token")
	// 	} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
	// 		// Token is either expired or not active yet
	// 		fmt.Println("Timing is everything")
	// 	} else {
	// 		fmt.Println("Couldn't handle this token:", err)
	// 	}
	// } else {
	// 	fmt.Println("Couldn't handle this token:", err)
	// }

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/", "updateSocket", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/", "join", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/", "addUser", ss.addUserHandler)

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	createFrontpage()

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func decInt(i int, by int) string {
	return fmt.Sprintf("%d", i-by)
}
func incInt(i int, by int) string {
	return fmt.Sprintf("%d", i+by)
}

type shipDesc struct {
	St   string
	Name string
}

type templateStruct struct {
	Letters []byte
	Numbers []int
	Names   []shipDesc
}

var ts = templateStruct{
	Letters: []byte{'/', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J'},
	Numbers: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	Names: []shipDesc{
		{"A", "Carrier (5)"},
		{"B", "Submarine (4)"},
		{"C", "Destroyer (3)"},
		{"D", "Cruiser (3)"},
		{"E", "Patrol (2)"},
	},
}

func createFrontpage() {
	fmt.Println("Creating file")
	var fm = template.FuncMap{
		"decInt": decInt,
		"incInt": incInt,
	}

	shortnerTemplate := template.Must(template.New("index.html").Funcs(fm).ParseFiles("./index.html"))

	fo, err := os.Create("static/index.html")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	w := bufio.NewWriter(fo)

	shortnerTemplate.Execute(w, ts)

	if err = w.Flush(); err != nil {
		panic(err)
	}

	fmt.Println("File created")
}

type Server struct {
	key []byte

	//addUserHandler func()
}

func (server *Server) addUserHandler(s socketio.Conn, msg string) {
	fmt.Println("addUser", msg)

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		panic(err)
	}

	name := dat["name"].(string)
	fmt.Println(dat)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": name,
		"nbf":  time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp":  time.Now().Unix() + 36000,
	})
	tokenString, err := token.SignedString(server.key)

	if err != nil {
		log.Fatal(err)
	}

	m := map[string]string{
		"msg":       "OK",
		"name":      name,
		"userToken": tokenString,
	}
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

	s.Emit("userAdded", string(b))
}
