$(document).ready(function () {

  let hostname = window.location.hostname;
  if (hostname === "localhost") {
    hostname += ":" + window.location.port;
  }

  // eslint-disable-next-line no-undef
  let socket = io.connect(hostname);
  let orginal = socket.emit.bind(socket);

  socket.emit = function (msg, data) {
    if (data != undefined) {
      data.userToken = userToken;
      data.gameToken = gameToken;
      orginal(msg, JSON.stringify(data));
    } else {
      orginal(msg);
    }
  };

  let username = 'Not Choosen';

  $('#globalLoading').hide();
  $('#namePrompt').hide();
  $('#joinGame').hide();
  $('#choosePlacement').hide();
  $('#board').hide();
  $('#gameOver').hide();

  let userToken = window.localStorage.getItem('userToken');
  let gameToken = window.localStorage.getItem('gameToken');
  let gameId = null;

  function deleteElement(id) {
    let toDelete = document.getElementById(id);
    let iskaParent = toDelete.parentNode;
    iskaParent.removeChild(toDelete);
  }

  //===== NAME

  if (window.localStorage.getItem('username')) {
    username = window.localStorage.getItem('username');
    socket.emit('updateSocket', { player: username });
  } else {
    $('#namePrompt').show();
    $('#errorName').text('.');
  }

  socket.on('updateFailed', function (_data) {
    username = null;
    userToken = null;
    window.localStorage.removeItem("userToken");
    window.localStorage.removeItem("username");
    $('#namePrompt').show();
    $('#errorName').text('.');
  });

  socket.on('updateSuccess', function (_data) {
    deleteElement('namePrompt');
    $('#joinGame').show();
  });


  let lockName = false;

  $('#btnSubmitName').on('click', function () {
    $('#errorName').text('.');
    if (lockName) {
      $('#errorName').text("Please Wait");
      return;
    }
    username = $('#inptName').val();
    let result = validateName(username);
    if (result != 'OK') {
      $('#errorName').text(result);
      return;
    }
    lockName = true;
    socket.emit('addUser', { name: username });
    $('#globalLoading').show();
  });

  socket.on('userAdded', function (data) {
    data = JSON.parse(data);
    console.log("userAdded", data);
    $('#globalLoading').hide();
    console.log("userAdded", data.msg);
    if (data.msg != 'OK') {
      username = null;
      lockName = false;
      $('#errorName').text(data.msg);
      return;
    }
    deleteElement('namePrompt');
    $('#joinGame').show();
    window.localStorage.setItem('username', data.name);
    window.localStorage.setItem('userToken', data.userToken);

    userToken = data.userToken;
  });

  function validateName(name) {
    if (name.length < 5) {
      return "Too Short. Minimum 5 characters";
    }
    if (name.length > 255) {
      return "Too Long. Maximum 255 characters";
    }
    if (/^\w+$/.test(name)) {
      return "OK";
    }
    return "Please Choose alphabets, numbers or '_'";
  }

  //========= JOIN

  let lockJoin = false;

  $('#btnJoin').click(function () {
    $('#errorJoin').text('.');
    if (lockJoin) {
      $('#errorJoin').text("Wait");
      return;
    }
    lockJoin = true;
    socket.emit('join', { player: username });
    $('#globalLoading').show();
  });

  socket.on('lockJoin', function (_data) {
    $('#errorJoin').text('Wait');
    lockJoin = true;
  });

  socket.on('startGame', function (data) {
    data = JSON.parse(data);
    lockJoin = false;
    deleteElement('joinGame');
    $('#globalLoading').hide();
    $('#choosePlacement').show();
    console.log('Player2 is' + data.otherPlayer);

    window.localStorage.setItem('gameToken', data.gameToken);
    gameToken = data.gameToken;
    gameId = data.gameId;
  });

  //========== BOARD INITIALIZATION

  let lockReady = false;
  let boardValid = false;

  $('#btnReady').click(function () {
    $('#errorReady').text('.');
    if (lockReady) {
      $('#errorReady').text("Wait");
      return;
    }
    boardValid = boardIsValid();
    if (!boardValid) {
      $('#errorReady').text('Invalid Board');
      return;
    }
    lockReady = true;
    for (let shipType in locked) {
      if (!Object.prototype.hasOwnProperty.call(locked, shipType)) {
        continue;
      }
      locked[shipType] = true;
    }
    let toSend = makeToSend();
    socket.emit('boardMade', { player: username, shipPlacement: toSend, gameId: gameId });
  });

  function makeToSend() {
    let arrToSend = {};
    for (let shipType in pointsOfShip) {
      if (!Object.prototype.hasOwnProperty.call(pointsOfShip, shipType)) {
        continue;
      }
      arrToSend[shipType] = {};
      let points = pointsOfShip[shipType];
      let z = points.keys();
      let i = 0;
      while (!z.done) {
        let point = z.next();
        point = point.value;
        if (!point) {
          break;
        }
        point = JSON.parse(point);
        arrToSend[shipType][i] = point;
        i++;
      }
    }
    return arrToSend;
  }

  function boardIsValid() {
    for (let shipType in pointsOfShip) {
      if (!Object.prototype.hasOwnProperty.call(pointsOfShip, shipType)) {
        continue;
      }
      if (pointsOfShip[shipType].size !== lengthOfType[shipType]) {
        return false;
      }
      if (!placedBefore[shipType]) {
        return false;
      }
    }
    return true;
  }

  let lengthOfType = { A: 5, B: 4, C: 3, D: 3, E: 2 };
  let arrOfI = ['1', '2', '3', '4', '5', '6', '7', '8', '9', '10'];
  let arrOfJ = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J'];

  function addShipClass(type, i, j, horizontal) {
    if (horizontal) {
      for (let y = j; y < j + lengthOfType[type]; y++) {
        console.log('#cell-' + i + "-" + y, $('#cell-' + i + "-" + y));
        $('#cell-' + i + "-" + y).addClass('ship' + type);
      }
    } else {
      for (let x = i; x < i + lengthOfType[type]; x++) {
        console.log('#cell-' + x + "-" + j, $('#cell-' + x + "-" + j));
        $('#cell-' + x + "-" + j).addClass('ship' + type);
      }
    }
  }

  let playerBoard = new Array(10);

  for (let i = 0; i < 10; i++) {
    playerBoard[i] = (new Array(10)).fill(0);
  }

  let pointsOfShip = {
    A: new Set(),
    B: new Set(),
    C: new Set(),
    D: new Set(),
    E: new Set(),
  };

  let hor = { A: false, B: false, C: false, D: false, E: false };
  let placedBefore = { A: false, B: false, C: false, D: false, E: false };
  let locked = { A: false, B: false, C: false, D: false, E: false };

  function addPointsToShip(type, i, j, horizontal) {
    let points = pointsOfShip[type];
    points.clear();
    if (horizontal) {
      for (let y = j; y < j + lengthOfType[type]; y++) {
        points.add(JSON.stringify({ 'x': i, 'y': y }));
      }
    } else {
      for (let x = i; x < i + lengthOfType[type]; x++) {
        points.add(JSON.stringify({ 'x': x, 'y': j }));
      }
    }
  }

  function removeShip(type) {
    let points = pointsOfShip[type];
    let z = points.keys();
    while (!z.done) {
      points = z.next();
      points = points.value;
      if (points) {
        points = JSON.parse(points);
        $('#cell-' + points.x + "-" + points.y).removeClass('ship' + type);
      } else {
        break;
      }
    }
    pointsOfShip[type] = new Set();
  }

  function choicesChanged(ship) {
    let valI = $('#' + ship + 'i').val();
    let valJ = $('#' + ship + 'j').val();
    valJ = valJ.toUpperCase();
    classInverser(ship, false);
    $('#errorPlaceShip' + ship).text(".");
    let possibleBounds = true;
    let possibleOverlap = true;
    if (valI && valJ) {
      if (arrOfI.indexOf(valI) > -1 && arrOfJ.indexOf(valJ) > -1) {
        possibleBounds = checkBounds(valI, valJ, ship, hor[ship]);
        if (possibleBounds) {
          possibleOverlap = checkOverlap(valI, valJ, ship, hor[ship]);
          if (possibleOverlap) {
            if (placedBefore[ship]) {
              removeShip(ship);
            }
            placedBefore[ship] = true;
            addPointsToShip(ship, arrOfI.indexOf(valI), arrOfJ.indexOf(valJ), hor[ship]);
            addShipClass(ship, arrOfI.indexOf(valI), arrOfJ.indexOf(valJ), hor[ship]);
          } else {
            classInverser(ship, true);
            $('#errorPlaceShip' + ship).text("Overlapping Ships");
          }
        } else {
          classInverser(ship, true);
          $('#errorPlaceShip' + ship).text("Out of bounds");
        }
      } else {
        classInverser(ship, true);
        $('#errorPlaceShip' + ship).text("Inavlid Entries");
      }
    }
  }

  function checkBounds(valI, valJ, ship, horizontal) {
    if (horizontal) {
      if (arrOfJ.indexOf(valJ) + lengthOfType[ship] > 10) {
        return false;
      }
    } else {
      if (arrOfI.indexOf(valI) + lengthOfType[ship] > 10) {
        return false;
      }
    }
    return true;
  }

  function checkOverlap(valI, valJ, ship, horizontal) {
    let j = arrOfJ.indexOf(valJ);
    let i = arrOfI.indexOf(valI);
    let tempPoints = new Set();
    if (horizontal) {
      for (let y = j; y < j + lengthOfType[ship]; y++) {
        tempPoints.add(JSON.stringify({ 'x': i, 'y': y }));
      }
    } else {
      for (let x = i; x < i + lengthOfType[ship]; x++) {
        tempPoints.add(JSON.stringify({ 'x': x, 'y': j }));
      }
    }
    for (let shipType in pointsOfShip) {
      if (!Object.prototype.hasOwnProperty.call(pointsOfShip, shipType)) {
        continue;
      }
      if (shipType != ship) {

        let intersection = new Set([...pointsOfShip[shipType]].filter(x => tempPoints.has(x)));
        if (intersection.size > 0) {
          return false;
        }

      }
    }
    return true;
  }

  $('.inptXY').change(function () {
    let ship = $(this).data("ship");
    if (locked[ship]) {
      classInverser(ship, true);
      $('#errorPlaceShip' + ship).text("Locked");
      return;
    }
    choicesChanged(ship);
  });

  $('.btnRot').click(function () {
    let ship = $(this).data("ship");
    if (locked[ship]) {
      classInverser(ship, true);
      $('#errorPlaceShip' + ship).text("Locked");
      return;
    }
    hor[ship] = !hor[ship];
    let dir = hor[ship] ? "Horizontal" : "Vertical";
    $('#btnRotIndic' + ship).text("Currently " + dir);
    choicesChanged(ship);
  });

  $('.btnDrop').click(function () {
    let ship = $(this).data("ship");
    if (locked[ship]) {
      classInverser(ship, true);
      $('#errorPlaceShip' + ship).text("Already Locked");
      return;
    }
    if (placedBefore[ship]) {
      this.classList.remove('btn-primary');
      this.classList.add('btn-danger');
      locked[ship] = true;
    } else {
      classInverser(ship, true);
      $('#errorPlaceShip' + ship).text("Please Place before locking");
    }
  });

  function classInverser(ship, errorOn) {
    let classToAdd = errorOn ? "label-danger" : "label-default";
    let classToRemove = errorOn ? "label-default" : "label-danger";

    $('#errorPlaceShip' + ship).removeClass(classToRemove);
    $('#errorPlaceShip' + ship).addClass(classToAdd);
  }

  socket.on('wait', function (data) {
    data = JSON.parse(data);
    if (data.status === "Error") {
      $('#errorReady').text(data.msg);
      return;
    }
    $('#globalLoading').show();
    cloneAndAppend();
    deleteElement('choosePlacement');
    $('#board').show();
  });

  function cloneAndAppend() {
    let toClone = document.getElementById("chooseBoard");
    let cloned = toClone.cloneNode(true);
    let clonedToChange = toClone.cloneNode(true);
    let targetParent = document.getElementById("battleBoard");
    changeId(clonedToChange);
    targetParent.appendChild(cloned);
    targetParent.appendChild(clonedToChange);
  }

  function changeId(x) {
    if (x.id) {
      x.id = "opp-" + x.id;
    }
    x.className = x.className.replace(/ ship[A-E]/, "");
    if (x.childElementCount > 0) {
      for (let i = 0; i < x.childElementCount; i++) {
        changeId(x.children[i]);
      }
    }
  }

  //let se = setInterval(function(){console.log(pointsOfShip);},2000);

  //=========== GAMEPLAY

  let myShips = {};
  let oppShips = {
    A: new Set(),
    B: new Set(),
    C: new Set(),
    D: new Set(),
    E: new Set(),
  };

  let otherPlayerBoard = new Array(10);

  for (let i = 0; i < 10; i++) {
    otherPlayerBoard[i] = (new Array(10)).fill(0);
  }

  let myTurn = false;

  let lastMove = {};

  socket.on('go', function (data) {
    data = JSON.parse(data);
    console.log("Go", data);
    if (data.start == "true") {
      $('#globalLoading').hide();
      myTurn = true;
    }
    for (let shipType in pointsOfShip) {
      if (!Object.prototype.hasOwnProperty.call(pointsOfShip, shipType)) {
        continue;
      }
      myShips[shipType] = new Set(pointsOfShip[shipType]);
    }
  });

  $('#btnShoot').click(function () {
    if (myTurn) {
      let x = $('#shoti').val();
      let y = $('#shotj').val();
      y = y.toUpperCase();
      classInverserShoot(false);
      if (arrOfI.indexOf(x) > -1 && arrOfJ.indexOf(y) > -1) {
        y = arrOfJ.indexOf(y);
        x = arrOfI.indexOf(x);
        if (otherPlayerBoard[x][y] === 0) {
          $('#globalLoading').show();
          socket.emit('makeMove', { player: username, move: { x: x, y: y }, gameId: gameId });
          lastMove.x = x;
          lastMove.y = y;
        } else {
          classInverserShoot(true);
          $('#errorShoot').text("Already");
        }
      } else {
        classInverserShoot(true);
        $('#errorShoot').text("Invalid Entry");
      }
    } else {
      classInverserShoot(true);
      $('#errorShoot').text("Please Wait For Your turn");
    }
  });

  function classInverserShoot(errorOn) {
    let classToAdd = errorOn ? "label-danger" : "label-default";
    let classToRemove = errorOn ? "label-default" : "label-danger";

    $('#errorShoot').addClass(classToAdd);
    $('#errorShoot').removeClass(classToRemove);
  }

  socket.on('yourMove', function (data) {
    data = JSON.parse(data);
    if (data.extra) {
      data.extra = JSON.parse(data.extra);
    }
    if (data.result === "Hit") {
      myTurn = !myTurn;
      otherPlayerBoard[lastMove.x][lastMove.y] = 1;
      $('#opp-cell-' + lastMove.x + "-" + lastMove.y).addClass("hit");
      oppShips[data.extra.partOf].add(JSON.stringify(lastMove));
      if (data.extra.shipDown) {
        markShipDown(data.extra.partOf);
        if (data.extra.gameOver) {
          //
          endGame(false);
        }
      }
    } else if (data.result === "Miss") {
      myTurn = !myTurn;
      otherPlayerBoard[lastMove.x][lastMove.y] = -1;
      $('#opp-cell-' + lastMove.x + "-" + lastMove.y).addClass("miss");
    } else {
      $('#errorShoot').text("Repeat");
    }
  });

  socket.on('oppMove', function (data) {
    data = JSON.parse(data);
    if (data.extra) {
      data.extra = JSON.parse(data.extra);
    }
    let cc = data.point.split(",");
    data.point = {
      x: parseInt(cc[0], 10),
      y: parseInt(cc[1], 10)
    };
    console.log(data.point);
    switch (data.result) {
      case "Hit":
        myTurn = !myTurn;
        $('#cell-' + data.point.x + "-" + data.point.y).addClass("hit");
        if (data.extra && data.extra.gameOver) {
          //
          endGame(true);
        }
        break;
      case "Miss":
        myTurn = !myTurn;
        $('#cell-' + data.point.x + "-" + data.point.y).addClass("miss");
        break;
    }
    $('#globalLoading').hide();
  });

  function markShipDown(type) {
    let points = oppShips[type];
    let z = points.keys();
    while (!z.done) {
      points = z.next();
      points = points.value;
      if (points) {
        points = JSON.parse(points);
        $('#opp-cell-' + points.x + "-" + points.y).addClass('ship' + type);
      } else {
        break;
      }
    }
  }

  function endGame(lost) {
    gameToken = null;
    gameId = null;
    window.localStorage.removeItem("gameToken");

    deleteElement('board');
    $('#globalLoading').hide();
    $('#gameOver').show();

    if (lost) {
      $('#gameOverWin').hide();
    } else {
      $('#gameOverLose').hide();
    }
  }

  //=============

  let dots = [
    document.getElementById('_dot1'),
    document.getElementById('_dot2'),
  ];

  const colorArr = [
    "#F44336", "#E91E63", "#9C27B0", "#673AB7",
    "#3F51B5", "#2196F3", "#03A9F4", "#00BCD4",
    "#009688", "#4CAF50", "#8BC34A", "#CDDC39",
    "#FFEB3B", "#FFC107", "#FF9800", "#FF5722",
  ];
  let colorArrLength = colorArr.length;

  function changeColors() {
    for (let dot of dots) {
      let z = parseInt(Math.random() * colorArrLength);
      dot.style.backgroundColor = colorArr[z];
    }
  }

  let _changeColorInterval = setInterval(changeColors, 1500);

});
