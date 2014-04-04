var fs = require('fs');
var WebSocket = require('ws');

var dbPath = process.env.DB_PATH || "msgs.json";

var msgs = [];

if (fs.existsSync(dbPath)) {
	msgs = JSON.parse(fs.readFileSync(dbPath));
}

var ws = new WebSocket(process.env.HOST_URL || "ws://localhost:8080");

ws.on('message', function(data) {
	msgs.push(JSON.parse(data));
	fs.writeFileSync(dbPath, JSON.stringify(msgs));
});
