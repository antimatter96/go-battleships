package stuff

import (
	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
)

/*
  5 ships

  shipA = 5
  shipB = 4
  shipC = 3
  shipD = 3
  shipE = 2
*/

type Player struct {
	conn *socketio.Conn
	name string
}

// Command has the text command and the json stuff
type Command struct {
	commandType string
	data        map[string]string
}

const (
	shipACode = iota
	shipBCode
	shipCCode
	shipDCode
	shipECode
)

var lengthOfType map[string]int = map[string]int{
	"A": 5,
	"B": 4,
	"C": 3,
	"D": 3,
	"E": 2,
}

var arrOfI = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
var arrOfJ = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

type BattleShips struct {
	ID string
	p1 string
	p2 string

	p1BoardDone bool
	p2BoardDone bool

	p1Board [][]int
	p2Board [][]int

	p1Ship map[int]StringSet
	p2Ship map[int]StringSet

	turnOf string
}

type Game interface {
	init() error

	PlayerReady(string, map[string][]int) error
	BothReady() (bool, error)

	OtherPlayer(string) (string, error)

	StartGame(string) error

	MakeMove(string, []int) error
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

	g.p1Ship = map[int]StringSet{
		shipACode: StringSet{},
		shipBCode: StringSet{},
		shipCCode: StringSet{},
		shipDCode: StringSet{},
		shipECode: StringSet{},
	}
	g.p2Ship = map[int]StringSet{
		shipACode: StringSet{},
		shipBCode: StringSet{},
		shipCCode: StringSet{},
		shipDCode: StringSet{},
		shipECode: StringSet{},
	}

	return nil
}

// StartGame sets the turn to the current player
func (g *BattleShips) StartGame(player string) {
	g.turnOf = player
}

// BothReady returns true when both player are true
func (g *BattleShips) BothReady() bool {
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

// class Game {

//   playerReady(player, shipPlacement) {
//     let playerBoardDone = this.p1BoardDone;
//     let playerShip = this.p1Ship;
//     let playerBoard = this.p1Board;

//     if (this.p2 === player) {
//       playerBoardDone = this.p2BoardDone;
//       playerShip = this.p2Ship;
//       playerBoard = this.p2Board;
//     }

//     if (playerBoardDone.bool) {
//       return {
//         thisPlayer: [
//           { message: "wait", data: { status: "Error", msg: "Already Choosen" } }
//         ],
//       };
//     }

//     for (let shipType in shipPlacement) {
//       let length = this.lengthOfType[shipType];
//       for (let i = 0; i < length; i++) {
//         let point = shipPlacement[shipType][i];
//         playerShip[shipType].add(JSON.stringify(point));
//         playerBoard[point.x][point.y] = 1;
//       }
//     }

//     playerBoardDone.bool = true;

//     if (this.bothReady()) {
//       this.startGame(player);
//       return {
//         thisPlayer: [
//           { message: "wait", data: { status: "OK", msg: "Done" } },
//           { message: "go", data: { status: "OK", start: true } }
//         ],
//         otherPlayer: [
//           { message: "go", data: { status: "OK", start: false } }
//         ]
//       };
//     } else {
//       return {
//         thisPlayer: [
//           { message: "wait", data: { status: "OK", msg: "Done" } },
//         ]
//       };
//     }
//   }
//
//   makeMove(player, move) {
//     if (this.turnOf != player) {
//       return {
//         thisPlayer: [{ message: 'moveError', data: { status: "Error", msg: "Not your turn" } }]
//       };
//     }
//     let x = move.x;
//     let y = move.y;
//     let point = { x: x, y: y };

//     var otherPlayerBoard;

//     var otherPlayerShip;

//     if (this.p1 === player) {
//       otherPlayerBoard = this.p2Board;
//       otherPlayerShip = this.p2Ship;
//     } else {
//       otherPlayerBoard = this.p1Board;
//       otherPlayerShip = this.p1Ship;
//     }

//     if (otherPlayerBoard[x][y] === 1) {
//       otherPlayerBoard[x][y] = -1;
//       let tempPoint = JSON.stringify(point);
//       let countZero = 0;
//       let extra = {};
//       for (var shipType in otherPlayerShip) {
//         if (otherPlayerShip[shipType].has(tempPoint)) {
//           otherPlayerShip[shipType].delete(tempPoint);
//           extra.partOf = shipType;
//           if (otherPlayerShip[shipType].size === 0) {
//             extra.shipDown = true;
//             countZero++;
//           }
//         } else if (otherPlayerShip[shipType].size === 0) {
//           countZero++;
//         }
//         if (countZero === 5) {
//           console.log("Over");
//           extra.gameOver = true;
//         }
//       }
//       this.turnOf = this.otherPlayer(player);
//       return {
//         thisPlayer: [{ message: "yourMove", data: { status: "OK", result: "Hit", extra: extra } }],
//         otherPlayer: [{ message: "oppMove", data: { status: "OK", result: "Hit", point: move, extra: extra } }]
//       };
//     } else if (otherPlayerBoard[x][y] === 0) {
//       otherPlayerBoard[x][y] = -1;
//       this.turnOf = this.otherPlayer(player);
//       return {
//         thisPlayer: [{ message: "yourMove", data: { status: "OK", result: "Miss" } }],
//         otherPlayer: [{ message: "oppMove", data: { status: "OK", result: "Miss", point: move } }]
//       };
//     } else {
//       return {
//         thisPlayer: [{ message: "yourMove", data: { status: "OK", result: "Repeat" } }],
//       };
//     }
//   }

// }

// package main

// type Game interface {
// 	AddPlayer()
// 	EndGame()
// 	Update(GameCommand) bool
// }

// var (
// 	ShipLength []int = []int{5, 4, 3, 3, 2}
// )

// const (
// 	ShipUp = iota
// 	ShipDown
// )

type BoardPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}
