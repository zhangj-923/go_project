import WebSocket from 'ws';
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
  const msg = JSON.parse(data.toString());
  console.log('Received type:', msg.type);
  if (msg.type === 'room_joined') {
    ws.send(JSON.stringify({ type: 'start_game' }));
  } else if (msg.type === 'game_started') {
    console.log('Game started successfully!');
    process.exit(0);
  }
});
