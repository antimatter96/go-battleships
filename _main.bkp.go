package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/gorilla/websocket"

	"github.com/antimatter96/go-battleships/game"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type websocketMessage struct {
	Command string                 `json:"command"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

var port = flag.String("port", "8080", "http service address")
var json = jsoniter.ConfigCompatibleWithStandardLibrary

var connOf map[string]*websocket.Conn

var st map[string]*websocket.Conn
var allocator *game.Allocator

func main() {
	allocator = game.NewAllocator()
	connOf = make(map[string]*websocket.Conn)

	x, err := game.NewBattleShips("asd", "asd")
	fmt.Println(x, err)
	var fm = template.FuncMap{
		"decInt": decInt,
		"incInt": incInt,
	}
	shortnerTemplate = template.Must(template.New("index.html").Funcs(fm).ParseFiles("./index.html"))

	flag.Parse()

	if (*port)[0] != ':' {
		*port = ":" + *port
	}
	fmt.Println("starting at", *port)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/ws", handleWebSocket)

	http.HandleFunc("/", serveHome)
	err = http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		log.Println(r.URL)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	shortnerTemplate.Execute(w, ts)
}

type daddyWebSocket struct {
	*websocket.Conn
	username string
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err, "error")
		http.Error(w, "some error", http.StatusInternalServerError)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("IsUnexpectedCloseError error: %v\n", err)
			} else {
				fmt.Printf("error: %v\n", err)
			}
			break
		}

		conn.SetCloseHandler(func(a int, v string) error {
			fmt.Println(a, v, message)
			return nil
		})

		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		var received websocketMessage
		err = json.Unmarshal(message, &received)

		fmt.Println("data =>", received, err)

		switch received.Command {
		case "addUser":
			fmt.Println("AddUser")
			fmt.Println("name =>", received.Data["name"])
			x := websocketMessage{"userAdded", map[string]interface{}{
				"status":   "OK",
				"username": received.Data["name"],
			}}
			yy, _ := json.Marshal(x)
			fmt.Println("marshalled => ", yy)
			err := writeToSocket(conn, yy)
			fmt.Println("writeToSocket err", err)
			userName, _ := received.Data["name"].(string)
			connOf[userName] = conn
		case "join":
			userName, _ := received.Data["name"].(string)
			connOf[userName] = conn
			output := make(chan string)
			allocator.Find(userName, output)
			timedOut := time.NewTimer(30 * time.Second)
			select {
			case <-timedOut.C:
				fmt.Println("Shit is over")
				allocator.IDontNeedAnyMore(userName)
				otherPlayer := <-output
				if otherPlayer == "" {
					fmt.Println("Other Player Not Found")
				} else {
					fmt.Println(otherPlayer, "NONONONON")
				}
			case otherPLayer := <-output:
				var err error
				err = writeToSocket(connOf[otherPLayer], []byte("AsAS"))
				fmt.Println(err)
				err = writeToSocket(conn, []byte("AsAS"))
				fmt.Println(err)
			}
			fmt.Println("connOf", connOf)
			fmt.Println("data.Data", received.Data)
		default:
			fmt.Println("Asd")
		}
	}

	fmt.Println("exiting loop")

}

func writeToSocket(conn *websocket.Conn, message []byte) error {
	conn.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))

	w, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	w.Write(message)

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func closeSocket(conn *websocket.Conn) {
	conn.WriteMessage(websocket.CloseMessage, []byte{})
}
