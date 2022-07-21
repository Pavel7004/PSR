const params = new Proxy(new URLSearchParams(window.location.search), {
  get: (searchParams, prop) => searchParams.get(prop),
});

let id = params.id;
let roomID = params.roomID;
let socket = new WebSocket(`ws://localhost:3000/echo?id=${id}&roomID=${roomID}`);
let statusLabel = document.getElementById('status');

socket.onopen = function (e) {};

socket.onclose = function (event) {
  if (!event.wasClean) {
    statusLabel.innerText = '[close] Соединение прервано';
  }
};

socket.onmessage = function (event) {
  statusLabel.innerText = event.data;
};

socket.onerror = function (error) {
  statusLabel.innerText = `[error] ${error.message}`;
};

function send() {
  let radios = document.getElementsByName('choice');
  for (let i = 0, len = radios.length; i < len; i++) {
    if (radios[i].checked) {
      socket.send(radios[i].value);
      statusLabel.innerText = 'Ожидание хода другого игрока';
      break;
    }
  }
}
