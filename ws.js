function onOpen(evt) {
	console.log("ws connection established");
	console.log(evt);
}

function onClose(evt) {
	console.log("closed", evt);
}

function onMessage(evt) {
	console.log(evt.data);
	console.log(evt);
}

function onError(evt) {
	console.log("error");
	console.log(evt);
}

var wsUri = "wss://localhost:3333";
ws = new WebSocket(wsUri);
ws.onopen = function(evt) { onOpen(evt) };
ws.onclose = function(evt) { onClose(evt) };
ws.onmessage = function(evt) { onMessage(evt) };
ws.onerror = function(evt) { onError(evt) };

ws.send("hello from browser");
ws.close();