const { v4: uuid } = require('uuid');

/*
  5 ships

  shipA = 5
  shipB = 4
  shipC = 3
  shipD = 3
  shipE = 2
*/

const lengthOfType = { A: 5, B: 4, C: 3, D: 3, E: 2 };
const lettersArray = ["A", "B", "C", "D", "E"];
const letters = new Set();
for (let i = 0; i < lettersArray.length; i++) {
  letters.add(lettersArray[i]);
}

class Game {
  constructor(player1, player2) {
    if (typeof (player1) !== "string" || typeof (player2) !== "string") {
      throw "player name missing";
    }

    if (player1.trim() === "" || player2.trim() === "") {
      throw "player name missing";
    }

    this.id = uuid();
    this.p1 = player1;
    this.p2 = player2;
    this.p1BoardDone = { bool: false };
    this.p2BoardDone = { bool: false };
    this.turnOf = this.p1;

    this.p1Board = new Array(10);
    this.p2Board = new Array(10);

    for (let i = 0; i < 10; i++) {
      this.p1Board[i] = (new Array(10)).fill(0);
      this.p2Board[i] = (new Array(10)).fill(0);
    }

    this.p1Ship = { A: new Set(), B: new Set(), C: new Set(), D: new Set(), E: new Set() };
    this.p2Ship = { A: new Set(), B: new Set(), C: new Set(), D: new Set(), E: new Set() };
  }

  playerReady(player, shipPlacement) {
    // add bounds and checks
    // player is not in p1,p2, not valid
    // playerShipment valid ?
    let playerBoardDone = this.p1BoardDone;
    let playerShip = this.p1Ship;
    let playerBoard = this.p1Board;

    if (this.p2 === player) {
      playerBoardDone = this.p2BoardDone;
      playerShip = this.p2Ship;
      playerBoard = this.p2Board;
    }

    if (playerBoardDone.bool) {
      return {
        thisPlayer: [
          { message: "wait", data: { status: "Error", msg: "Already Choosen" } }
        ],
      };
    }

    for (let shipType in shipPlacement) {
      if (!shipPlacement.hasOwnProperty(shipType)) { continue; }
      let length = lengthOfType[shipType];
      for (let i = 0; i < length; i++) {
        let point = shipPlacement[shipType][i];
        playerShip[shipType].add(JSON.stringify(point));
        playerBoard[point.x][point.y] = 1;
      }
    }

    playerBoardDone.bool = true;

    if (this.bothReady()) {
      this.startGame(player);
      return {
        thisPlayer: [
          { message: "wait", data: { status: "OK", msg: "Done" } },
          { message: "go", data: { status: "OK", start: true } }
        ],
        otherPlayer: [
          { message: "go", data: { status: "OK", start: false } }
        ]
      };
    } else {
      return {
        thisPlayer: [
          { message: "wait", data: { status: "OK", msg: "Done" } },
        ]
      };
    }
  }

  bothReady() {
    return this.p1BoardDone.bool && this.p2BoardDone.bool;
  }

  otherPlayer(player) {
    // Add check on other player, throw error if neccessary
    if (player === this.p1) {
      return this.p2;
    } else {
      return this.p1;
    }
  }

  startGame(player) {
    // Add check on that
    this.turnOf = player;
  }

  makeMove(player, move) {
    if (this.turnOf !== player) {
      return {
        thisPlayer: [{ message: 'moveError', data: { status: "Error", msg: "Not your turn" } }]
      };
    }
    let x = move.x;
    let y = move.y;
    let point = { x: x, y: y };

    var otherPlayerBoard;

    var otherPlayerShip;

    if (this.p1 === player) {
      otherPlayerBoard = this.p2Board;
      otherPlayerShip = this.p2Ship;
    } else {
      otherPlayerBoard = this.p1Board;
      otherPlayerShip = this.p1Ship;
    }

    if (otherPlayerBoard[x][y] === 1) {
      otherPlayerBoard[x][y] = -1;
      let tempPoint = JSON.stringify(point);
      let countZero = 0;
      let extra = {};
      for (let shipType in otherPlayerShip) {
        if (!otherPlayerShip.hasOwnProperty(shipType)) { continue; }
        if (otherPlayerShip[shipType].has(tempPoint)) {
          otherPlayerShip[shipType].delete(tempPoint);
          extra.partOf = shipType;
          if (otherPlayerShip[shipType].size === 0) {
            extra.shipDown = true;
            countZero++;
          }
        } else if (otherPlayerShip[shipType].size === 0) {
          countZero++;
        }
        if (countZero === 5) {
          console.log("Over");
          extra.gameOver = true;
        }
      }
      this.turnOf = this.otherPlayer(player);
      return {
        thisPlayer: [{ message: "yourMove", data: { status: "OK", result: "Hit", extra: extra } }],
        otherPlayer: [{ message: "oppMove", data: { status: "OK", result: "Hit", point: move, extra: extra } }]
      };
    } else if (otherPlayerBoard[x][y] === 0) {
      otherPlayerBoard[x][y] = -1;
      this.turnOf = this.otherPlayer(player);
      return {
        thisPlayer: [{ message: "yourMove", data: { status: "OK", result: "Miss" } }],
        otherPlayer: [{ message: "oppMove", data: { status: "OK", result: "Miss", point: move } }]
      };
    } else {
      return {
        thisPlayer: [{ message: "yourMove", data: { status: "OK", result: "Repeat" } }],
      };
    }
  }

  /** 
   * Utility functions
  */

  toJSON () {
    return {
      id: this.id,
      p1: this.p1, p2: this.p2,
      turnOf: this.turnOf,
      p1BoardDone: this.p1BoardDone, p2BoardDone: this.p2BoardDone,
      p1Board: this.p1Board, p2Board: this.p2Board,
      p1Ship: JSON.stringify(this.p1Ship, Game.setToJson), p2Ship: JSON.stringify(this.p2Ship, Game.setToJson),
    };
  }

  static setToJson(_key, value) {
    if (typeof value === 'object' && value instanceof Set) {
      return [...value];
    }
    return value;
  }

  static jsonToSet(key, value) {
    if (key == "p1Ship" || key == "p2Ship") {
      return JSON.parse(value, Game.jsonToSet);
    }

    if (letters.has(key)) {
      let set = new Set();
      for(let i = 0; i < value.length; i++) {
        set.add(value[i]);
      }
      return set;
    }
    return value;
  }

  static gameFromString(string) {
    let gameJSON = JSON.parse(string, Game.jsonToSet);

    let game = new Game(gameJSON.p1, gameJSON.p2);

    game.id = gameJSON["id"];
    game.turnOf = gameJSON["turnOf"];
    game.p1BoardDone = gameJSON["p1BoardDone"];
    game.p2BoardDone = gameJSON["p2BoardDone"];

    game.p1Board = gameJSON["p1Board"];
    game.p2Board = gameJSON["p2Board"];

    game.p1Ship = gameJSON["p1Ship"];
    game.p2Ship = gameJSON["p2Ship"];

    return game;
  }

}
module.exports = Game;
