function onOpen(evt) {
	alert("ws connection established");
	console.log(evt);
}

function onClose(evt) {
	alert("closing");
	console.log(evt);
}

function onMessage(evt) {
	alert(evt.data);
	console.log(evt);
}

function onError(evt) {
	alert("error");
	console.log(evt);
}

var wsUri = "ws://127.0.0.1:3333";
ws = new WebSocket(wsUri);
ws.onopen = function(evt) { onOpen(evt) };
ws.onclose = function(evt) { onClose(evt) };
ws.onmessage = function(evt) { onMessage(evt) };
ws.onerror = function(evt) { onError(evt) };

ws.send("hello from browser");