package main

import (
	"io/ioutil"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"

	stuff "github.com/antimatter96/go-battleships/stuff"
)

func main() {
	privateKeyPEM, err := ioutil.ReadFile("./private_key.pem")

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
