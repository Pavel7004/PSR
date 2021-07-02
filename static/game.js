function getParameterByName(name, url = window.location.href) {
  name = name.replace(/[\[\]]/g, "\\$&");
  var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
    results = regex.exec(url);
  if (!results) return null;
  if (!results[2]) return "";
  return decodeURIComponent(results[2].replace(/\+/g, " "));
}

let id = getParameterByName("id");
let socket = new WebSocket(`ws://localhost:3000/echo?id=${id}`);
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
