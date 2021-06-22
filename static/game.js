let socket = new WebSocket("ws://localhost:3000/echo?id={{ .ID }}");
socket.onopen = function (e) {};

socket.onclose = function (event) {
  if (!event.wasClean) {
    alert("[close] Соединение прервано");
  }
};

socket.onmessage = function (event) {
  alert(event.data);
};

socket.onerror = function (error) {
  alert(`[error] ${error.message}`);
};

function send() {
  let radios = document.getElementsByName("choice");
  for (let i = 0, len = radios.length; i < len; i++) {
    if (radios[i].checked) {
      socket.send(radios[i].value);
      break;
    }
  }
}
