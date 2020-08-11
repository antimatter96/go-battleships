package stuff

import (
	"encoding/json"
	"time"

	"github.com/dgrijalva/jwt-go"
	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/zerolog/log"
)

func (server *Server) addUserHandler(s socketio.Conn, msg string) {
	sublogger := log.With().Str("service", "addUserHandler").Logger()

	//sublogger.Debug().Msgf("Data : %s", msg)

	var dat map[string]string
	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		sublogger.Error().AnErr("JSON unmarshalling error", err)
		panic(err)
	}
	sublogger.Debug().Msgf("Data : %+v", dat)

	m := make(map[string]string)
	defer func() {
		b, err := json.Marshal(m)
		if err != nil {
			sublogger.Error().AnErr("JSON marshalling error", err)
		}
		s.Emit("userAdded", string(b))
	}()

	name := dat["name"]

	if _, present := server.present[name]; present {
		m["msg"] = "Username taken"
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  name,
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp": time.Now().Unix() + 36000,
	})
	tokenString, err := token.SignedString(server.Key)
	if err != nil {
		sublogger.Error().AnErr("Token signing", err)
		panic(err)
	}

	m["msg"] = "OK"
	m["name"] = name
	m["userToken"] = tokenString

	s.SetContext(name)

	server.present[name] = true
	server.socketOf[name] = &s
}
