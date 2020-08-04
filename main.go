package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	socketio "github.com/googollee/go-socket.io"

	stuff "github.com/antimatter96/go-battleships/stuff"
)

type MyCustomClaims struct {
	Foo string `json:"foo"`
	ID  string `json:"id"`
	jwt.StandardClaims
}

func main() {

	privateKeyPEM, err := ioutil.ReadFile("./private_key.pem")

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

	ss := stuff.Server{Key: privateKeyPEM, Server: server}
	ss.Init()

	stuff.CreateFrontpage()

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
