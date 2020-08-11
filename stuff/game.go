package stuff

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var lengthOfType = map[string]int{
	"A": 5,
	"B": 4,
	"C": 3,
	"D": 3,
	"E": 2,
}

type BattleShips struct {
	sync.Mutex

	ID string
	p1 string
	p2 string

	p1BoardDone bool
	p2BoardDone bool

	p1Board [][]int
	p2Board [][]int

	p1Ship map[string]*stringSet
	p2Ship map[string]*stringSet

	turnOf string
}

func NewBattleShips(p1, p2 string) (*BattleShips, error) {
	game := &BattleShips{p1: p1, p2: p2}
	err := game.init()
	return game, err
}

func (g *BattleShips) init() error {
	gameID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	g.ID = gameID.String()

	// Default Value FTW
	// g.p1BoardDone = false
	// g.p2BoardDone = false

	g.turnOf = g.p1

	g.p1Board = make([][]int, 10)
	g.p2Board = make([][]int, 10)

	for i := 0; i < 10; i++ {
		g.p1Board[i] = make([]int, 10)
		g.p2Board[i] = make([]int, 10)
	}

	g.p1Ship = make(map[string]*stringSet)
	g.p2Ship = make(map[string]*stringSet)

	for shipCode := range lengthOfType {
		g.p1Ship[shipCode] = &stringSet{}
		g.p2Ship[shipCode] = &stringSet{}
	}

	return nil
}

// StartGame sets the turn to the current player
func (g *BattleShips) StartGame(player string) {
	g.turnOf = player
}

// BothReady returns true when both player are true
func (g *BattleShips) BothReady() bool {
	g.Lock()
	defer g.Unlock()
	return g.p1BoardDone && g.p2BoardDone
}

// OtherPlayer returns the other player thean the given player
func (g *BattleShips) OtherPlayer(player string) string {
	// Add error if none
	if player == g.p1 {
		return g.p2
	}
	return g.p1
}

// PlayerReady is
func (g *BattleShips) PlayerReady(player string, sp shipPlacement) gameResponse {
	g.Lock()
	defer g.Unlock()

	pd := &g.p1BoardDone
	ps := &g.p1Ship
	pb := &g.p1Board

	if player == g.p2 {
		pd = &g.p2BoardDone
		ps = &g.p2Ship
		pb = &g.p2Board
	}

	thisPlayer := []response{{message: "", data: map[string]string{}}}
	thisPlayer[0].message = "wait"

	if *pd {
		thisPlayer[0].data["msg"] = "Already Choosen"
		thisPlayer[0].data["status"] = "Error"

		//return fmt.Errorf("%s", "Already Choosen")

		return gameResponse{
			thisPlayerRes: thisPlayer,
		}

	}

	for k, v := range sp {
		for _, vv := range v {
			(*ps)[k].add(vv.String())
			(*pb)[vv.X][vv.Y] = 1
		}
	}

	*pd = true

	thisPlayer[0].data["msg"] = "Done"
	thisPlayer[0].data["status"] = "OK"

	if g.BothReady() {
		thisPlayer = append(thisPlayer, response{"go", map[string]string{
			"status": "OK",
			"start":  "true",
		}})

		otherPlayer := response{"go", map[string]string{
			"status": "OK",
			"start":  "false",
		}}

		return gameResponse{
			thisPlayerRes:  thisPlayer,
			otherPlayerRes: []response{otherPlayer},
		}
	}

	return gameResponse{
		thisPlayerRes: thisPlayer,
	}
}

// MakeMove is used to
func (g *BattleShips) MakeMove(player string, point boardPoint) gameResponse {
	g.Lock()
	defer g.Unlock()

	thisPlayer := []response{{message: "", data: map[string]string{}}}

	if g.turnOf != player {
		thisPlayer[0].message = "moveError"
		thisPlayer[0].data["msg"] = "Not your turn"
		thisPlayer[0].data["status"] = "Error"

		return gameResponse{
			thisPlayerRes: thisPlayer,
		}
	}

	x, y := point.X, point.Y

	ps := &g.p2Ship
	pb := &g.p2Board

	// Add error if none
	if player == g.p2 {
		ps = &g.p1Ship
		pb = &g.p1Board
	}

	thisPlayer[0].message = "yourMove"
	thisPlayer[0].data["status"] = "OK"

	if (*pb)[x][y] == 1 {
		(*pb)[x][y] = -1
		g.turnOf = g.OtherPlayer(player)

		tempPoint := point.String()
		countZero := 0
		extra := &extra{}

		for k, v := range *ps {
			// TODO : Simply delete and check result, no need to check and delete
			if v.has(tempPoint) {
				(*v).delete(tempPoint)
				extra.ShipType = k
				if v.size() == 0 {
					extra.ShipDown = true
				}
			}
		}

		for _, v := range *ps {
			if v.size() == 0 {
				countZero++
			}
		}

		if countZero == 5 {
			extra.GameOver = true
		}

		b, err := json.Marshal(extra)
		if err != nil {
			panic(err)
		}

		thisPlayer[0].data["result"] = "Hit"
		thisPlayer[0].data["extra"] = string(b)

		otherPlayer := response{"oppMove", map[string]string{
			"status": "OK",
			"result": "Miss",
			"point":  point.String(),
			"extra":  string(b),
		}}

		return gameResponse{
			thisPlayerRes:  thisPlayer,
			otherPlayerRes: []response{otherPlayer},
		}
	} else if (*pb)[x][y] == 0 {
		(*pb)[x][y] = -1
		g.turnOf = g.OtherPlayer(player)

		thisPlayer[0].data["result"] = "Miss"

		otherPlayer := response{"oppMove", map[string]string{
			"status": "OK",
			"result": "Hit",
			"point":  point.String(),
		}}

		return gameResponse{
			thisPlayerRes:  thisPlayer,
			otherPlayerRes: []response{otherPlayer},
		}

	} else {
		thisPlayer[0].data["result"] = "Repeat"

		return gameResponse{
			thisPlayerRes: thisPlayer,
		}
	}
}

type boardPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p *boardPoint) String() string {
	return fmt.Sprintf("%02d,%02d", p.X, p.Y)
}

type extra struct {
	ShipDown bool   `json:"shipDown,omitempty"`
	GameOver bool   `json:"gameOver,omitempty"`
	ShipType string `json:"partOf,omitempty"`
}
