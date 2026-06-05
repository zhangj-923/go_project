const WebSocket = require('ws');
const ws = new WebSocket('ws://localhost:18189/ws');

ws.on('open', () => {
  console.log('Connected');
  ws.send(JSON.stringify({
    type: 'create_room',
    payload: {
      playerName: 'TestPlayer',
      gameMode: 'single',
      maxPlayers: 4,
      aiCount: 3
    }
  }));
});

ws.on('message', (data) => {
  console.log('Received:', data.toString());
  process.exit(0);
});

ws.on('error', (err) => {
  console.error('Error:', err);
  process.exit(1);
});

setTimeout(() => {
  console.log('Timeout');
  process.exit(1);
}, 2000);
