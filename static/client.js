$(document).ready(function () {

  const lengthOfType = { A: 5, B: 4, C: 3, D: 3, E: 2 };
  const arrOfI = ['1', '2', '3', '4', '5', '6', '7', '8', '9', '10'];
  const arrOfJ = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J'];

  var hostname = window.location.hostname;
  if (hostname === "localhost") {
    hostname = "localhost" + ":" + window.location.port;
  }

  var socket = new WebSocket("ws://" + hostname + "/ws");

  socket.emit = function() {};
  socket.onmessage = function(event) {
    let message = undefined;
    console.log(event.data);
    try {
      message = JSON.parse(event.data);
    } catch(e) {
      console.log(e);
      return;
    }

    console.log(message);
    if (!message) {
      return;
    }

    switch(message.command) {
      case "userAdded":
        userAdded(message.data);
        break;
    }

  };

  var username = 'Not Choosen';

  $('#globalLoading').hide();
  $('#namePrompt').hide();
  $('#joinGame').hide();
  $('#choosePlacement').hide();
  $('#board').hide();
  $('#gameOver').hide();

  function deleteElement(id) {
    let toDelete = document.getElementById(id);
    let iskaParent = toDelete.parentNode;
    iskaParent.removeChild(toDelete);
  }

  //===== NAME

  var send = sender(socket);

  if (window.localStorage.getItem('username')) {
    username = window.localStorage.getItem('username');
    send( JSON.stringify( {command: 'updateSocket', data: { player: username } }));
    deleteElement('namePrompt');
    $('#joinGame').show();
  }
  else {
    $('#namePrompt').show();
    $('#errorName').text('.');
  }

  var lockName = false;

  $('#btnSubmitName').on('click', function () {
    $('#errorName').text('.');
    if (lockName) {
      $('#errorName').text("Please Wait");
      return;
    }
    let result = validateName($('#inptName').val());
    if (result != 'OK') {
      $('#errorName').text(result);
      return;
    }
    lockName = true;
    socket.send( JSON.stringify( {command: 'addUser', data: { name: $('#inptName').val() }} ));
    $('#globalLoading').show();
  });

  function userAdded (data) {
    $('#globalLoading').hide();
    if (data.status != 'OK') {
      lockName = false;
      $('#errorName').text(data.msg);
      return;
    }
    username = data.username;
    deleteElement('namePrompt');
    $('#joinGame').show();
    window.localStorage.setItem('username', data.name);
  }

  function validateName(name) {
    if (name.length < 5) {
      return "Too Short. Minimum 5 characters";
    }
    if (name.length > 25) {
      return "Too Long. Maximum 25 characters";
    }
    if (/^\w+$/.test(name)) {
      return "OK";
    }
    return "Please Choose alphabets, numbers or '_'";
  }

  //========= JOIN

  var lockJoin = false;

  $('#btnJoin').click(function () {
    $('#errorJoin').text('.');
    if (lockJoin) {
      $('#errorJoin').text("Wait");
      return;
    }
    lockJoin = true;
    socket.send( JSON.stringify({ command:'join', data:{ player: username } }));
    $('#globalLoading').show();
  });

  function lockJoin(data) {
    $('#errorJoin').text('Wait');
    lockJoin = true;
  }

  function startGame(data) {
    lockJoin = false;
    deleteElement('joinGame');
    $('#globalLoading').hide();
    $('#choosePlacement').show();
    console.log('Player2 is' + data.otherPlayer);
  }

  //========== BOARD INITIALIZATION

  var lockReady = false;
  var boardValid = false;

  $('#btnReady').click(function () {
    $('#errorReady').text('.');
    if (lockReady) {
      $('#errorReady').text("Wait");
      return;
    }
    boardValid = boardIsValid();
    if (!boardValid) {
      $('#errorReady').text('Invlid Board');
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
    socket.send( JSON.stringify({ command:'boardMade', data: { player: username, shipPlacement: toSend }}));
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

  function addShipClass(type, i, j, horizontal) {
    if (horizontal) {
      for (let y = j; y < j + lengthOfType[type]; y++) {
        $('#cell-' + i + y).addClass('ship' + type);
      }
    }
    else {
      for (let x = i; x < i + lengthOfType[type]; x++) {
        $('#cell-' + x + j).addClass('ship' + type);
      }
    }
  }

  let playerBoard = new Array(10);

  for (let i = 0; i < 10; i++) {
    playerBoard[i] = (new Array(10)).fill(0);
  }

  var pointsOfShip = {
    A: new Set(),
    B: new Set(),
    C: new Set(),
    D: new Set(),
    E: new Set(),
  };

  var hor = { A: false, B: false, C: false, D: false, E: false };
  var placedBefore = { A: false, B: false, C: false, D: false, E: false };
  var locked = { A: false, B: false, C: false, D: false, E: false };

  function addPointsToShip(type, i, j, horizontal) {
    let points = pointsOfShip[type];
    points.clear();
    if (horizontal) {
      for (let y = j; y < j + lengthOfType[type]; y++) {
        points.add(JSON.stringify({ 'x': i, 'y': y }));
      }
    }
    else {
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
        $('#cell-' + points.x + points.y).removeClass('ship' + type);
      }
      else {
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
          }
          else {
            classInverser(ship, true);
            $('#errorPlaceShip' + ship).text("Overlapping Ships");
          }
        }
        else {
          classInverser(ship, true);
          $('#errorPlaceShip' + ship).text("Out of bounds");
        }
      }
      else {
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
    }
    else {
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
    }
    else {
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
    }
    else {
      classInverser(ship, true);
      $('#errorPlaceShip' + ship).text("Please Place before locking");
    }
  });

  function classInverser(ship, errorOn) {
    let classToAdd = "label-default";
    let classToRemove = "label-danger";

    if (errorOn) {
      classToAdd = "label-danger";
      classToRemove = "label-default";
    }

    $('#errorPlaceShip' + ship).removeClass(classToRemove);
    $('#errorPlaceShip' + ship).addClass(classToAdd);
  }

  function wait(data) {
    if (data.status === "Error") {
      $('#errorReady').text(data.msg);
      return;
    }
    $('#globalLoading').show();
    cloneAndAppend();
    deleteElement('choosePlacement');
    $('#board').show();
  }

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

  //var se = setInterval(function(){console.log(pointsOfShip);},2000);

  //=========== GAMEPLAY

  var myShips = {};
  var oppShips = {
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

  var myTurn = false;

  var lastMove = {};

  function go(data) {
    if (data.start) {
      $('#globalLoading').hide();
      myTurn = true;
    }
    for (let shipType in pointsOfShip) {
      if (!Object.prototype.hasOwnProperty.call(pointsOfShip, shipType)) {
        continue;
      }
      myShips[shipType] = new Set(pointsOfShip[shipType]);
    }
  }

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
          socket.send( JSON.stringify({ command: 'makeMove', data: { player: username, move: { x: x, y: y } }}));
          lastMove.x = x;
          lastMove.y = y;
        }
        else {
          classInverserShoot(true);
          $('#errorShoot').text("Already");
        }
      }
      else {
        classInverserShoot(true);
        $('#errorShoot').text("Invlid Entry");
      }
    }
    else {
      classInverserShoot(true);
      $('#errorShoot').text("Please Wait For Your turn");
    }
  });

  function classInverserShoot(errorOn) {
    if (errorOn) {
      $('#errorShoot').addClass("label-danger");
      $('#errorShoot').removeClass("label-default");
    }
    else {
      $('#errorShoot').removeClass("label-danger");
      $('#errorShoot').addClass("label-default");
    }
  }

  function yourMove(data) {
    if (data.result === "Hit") {
      myTurn = !myTurn;
      otherPlayerBoard[lastMove.x][lastMove.y] = 1;
      $('#opp-cell-' + lastMove.x + lastMove.y).addClass("hit");
      oppShips[data.extra.partOf].add(JSON.stringify(lastMove));
      if (data.extra.shipDown) {
        markShipDown(data.extra.partOf);
        if (data.extra.gameOver) {
          //
          deleteElement('board');
          $('#gameOver').show();
          $('#gameOverLose').hide();
        }
      }
    }
    else if (data.result === "Miss") {
      myTurn = !myTurn;
      otherPlayerBoard[lastMove.x][lastMove.y] = -1;
      $('#opp-cell-' + lastMove.x + lastMove.y).addClass("miss");
    }
    else {
      $('#errorShoot').text("Repeat");
    }
  }

  function oppMove(data) {
    switch (data.result) {
      case "Hit":
        myTurn = !myTurn;
        $('#cell-' + data.point.x + data.point.y).addClass("hit");
        if (data.extra.gameOver) {
          //
          deleteElement('board');
          $('#globalLoading').hide();
          $('#gameOver').show();
          $('#gameOverWin').hide();
        }
        break;
      case "Miss":
        myTurn = !myTurn;
        $('#cell-' + data.point.x + data.point.y).addClass("miss");
        break;
    }
    $('#globalLoading').hide();
  }

  function markShipDown(type) {
    let points = oppShips[type];
    let z = points.keys();
    while (!z.done) {
      points = z.next();
      points = points.value;
      if (points) {
        points = JSON.parse(points);
        $('#opp-cell-' + points.x + points.y).addClass('ship' + type);
      }
      else {
        break;
      }
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

  let changeColorInterval = setInterval(changeColors, 1500);

  function sender(socket) {
    let arr = [];

    function send(data) {
      for(let i = 0; i < arr.length; i++) {
        socket.send(arr[i]);
      }
      socket.send(data);
    }

    return function(data) {
      if (socket.readyState != 1) {
        arr.push(data);
      } else {
        send(data);
      }
    };
  }

});
