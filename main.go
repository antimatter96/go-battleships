package main

import (
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"

	socketio "github.com/googollee/go-socket.io"

	stuff "github.com/antimatter96/go-battleships/stuff"
)

func main() {
	privateKeyPEM, err := ioutil.ReadFile("./private_key.pem")
	if err != nil {
		log.Warn().AnErr("Error getting private key", err)
	}

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal().AnErr("Error initialising socket server", err)

	}

	ss := stuff.Server{Key: privateKeyPEM, Server: server}
	ss.Init()

	stuff.CreateFrontpage()

	go func() {
		if errSocket := server.Serve(); errSocket != nil {
			log.Warn().AnErr("Error in socket server receiving connections", errSocket)
		}
	}()

	defer server.Close()

	http.Handle("/socket.io/", server)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Info().Msg("Serving at localhost:8000...")
	log.Fatal().AnErr("Error http server", http.ListenAndServe(":8000", nil))
}
