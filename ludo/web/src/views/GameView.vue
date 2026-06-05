<script setup lang="ts">
import { onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useGameStore } from '../stores/gameStore'
import { useWebSocket } from '../composables/useWebSocket'
import GameBoard from '../components/GameBoard.vue'
import DicePanel from '../components/DicePanel.vue'

const router = useRouter()
const store = useGameStore()
const ws = useWebSocket()

const audioDice = new Audio('/sounds/dice.mp3')
const audioMove = new Audio('/sounds/move.mp3')
const audioEat = new Audio('/sounds/eat.mp3')
const audioWin = new Audio('/sounds/win.mp3')

onMounted(() => {
  if (!store.roomId) {
    router.push('/')
    return
  }

  ws.onMessage('player_joined', onPlayerJoined)
  ws.onMessage('player_left', onPlayerLeft)
  ws.onMessage('game_started', onGameStarted)
  ws.onMessage('dice_rolled', onDiceRolled)
  ws.onMessage('piece_moved', onPieceMoved)
  ws.onMessage('turn_changed', onTurnChanged)
  ws.onMessage('game_over', onGameOver)
})

onUnmounted(() => {
  ws.removeHandler('player_joined', onPlayerJoined)
  ws.removeHandler('player_left', onPlayerLeft)
  ws.removeHandler('game_started', onGameStarted)
  ws.removeHandler('dice_rolled', onDiceRolled)
  ws.removeHandler('piece_moved', onPieceMoved)
  ws.removeHandler('turn_changed', onTurnChanged)
  ws.removeHandler('game_over', onGameOver)
})

const onPlayerJoined = (player: any) => {
  store.board.players.push(player)
}

const onPlayerLeft = (playerId: string) => {
  const p = store.board.players.find(p => p.id === playerId)
  if (p) p.connected = false
}

const onGameStarted = (payload: any) => {
  store.board = payload.boardState
  store.currentTurn = payload.currentTurn
  store.gameState = 'playing'
}

const onDiceRolled = (payload: any) => {
  audioDice.play().catch(() => {})
  store.diceValue = payload.value
  store.isRolling = true
  store.continuousSix = payload.continuousSix
  setTimeout(() => {
    store.isRolling = false
    store.movablePieces = payload.movablePieces
  }, 600)
}

const onPieceMoved = (payload: any) => {
  const isEat = payload.events.some((e: any) => e.type === 'eat')
  if (isEat) {
    audioEat.play().catch(() => {})
  } else {
    audioMove.play().catch(() => {})
  }
  // Simplified direct update
  const player = store.board.players.find(p => p.color === payload.color)
  if (player) {
    const piece = player.pieces.find(p => p.id === payload.pieceId)
    if (piece) piece.position = payload.to
  }
  // Also handle eaten pieces
  payload.events.forEach((e: any) => {
    if (e.type === 'eat') {
      store.board.players.forEach(p => {
        if (p.color !== payload.color) {
          p.pieces.forEach(pi => {
            if (pi.position === payload.to) pi.position = -1
          })
        }
      })
    }
  })
}

const onTurnChanged = (payload: any) => {
  store.currentTurn = payload.color
  store.diceValue = 0
  store.movablePieces = []
}

const onGameOver = (payload: any) => {
  audioWin.play().catch(() => {})
  store.gameState = 'game_over'
  store.winnerId = payload.winnerId
}

const startGame = () => {
  ws.send('start_game')
}

const resetStore = () => {
  store.roomId = ''
  store.myId = ''
  store.myColor = null
  store.gameState = 'lobby'
  store.board = { players: [] }
  store.currentTurn = null
  store.diceValue = 0
  store.isRolling = false
  store.movablePieces = []
  store.continuousSix = 0
  store.winnerId = ''
}

const leaveRoom = () => {
  ws.send('leave_room')
  resetStore()
  router.push('/')
}

const amIHost = computed(() => {
  return store.board && store.board.players && store.board.players.length > 0 && store.board.players[0].id === store.myId
})
</script>

<template>
  <div class="game-container">
    <div v-if="store.gameState === 'waiting'" class="waiting-room">
      <div class="glass-panel">
        <h2>房间: {{ store.roomId }}</h2>
        <div class="players-list">
          <div v-for="p in store.board.players" :key="p.id" class="player-item" :class="p.color">
            <span class="color-dot"></span>
            {{ p.name }} {{ p.id === store.myId ? '(你)' : '' }}
          </div>
        </div>
        
        <div class="actions">
          <button v-if="amIHost" class="btn-primary" @click="startGame">开始游戏</button>
          <div v-else class="waiting-text">等待房主开始...</div>
          <button class="btn-secondary" @click="leaveRoom">离开房间</button>
        </div>
      </div>
    </div>

    <div v-else-if="store.gameState === 'playing'" class="playing-area">
      <div class="left-panel">
        <div class="room-info">房间: {{ store.roomId }}</div>
        <!-- Player list / status -->
        <div class="players-status">
          <div v-for="p in store.board.players" :key="p.id" 
               class="status-card" 
               :class="{ active: store.currentTurn === p.color, [p.color]: true }">
            <div class="name">{{ p.name }}</div>
            <div class="pieces-left">剩余: {{ p.pieces.filter(x => x.position !== 105 && x.position !== 115 && x.position !== 125 && x.position !== 135).length }}</div>
          </div>
        </div>
      </div>
      
      <div class="board-container">
        <GameBoard />
      </div>
      
      <div class="right-panel">
        <DicePanel />
      </div>
    </div>
    
    <div v-else-if="store.gameState === 'game_over'" class="game-over-modal">
      <div class="glass-panel">
        <h1>游戏结束!</h1>
        <p>获胜者: {{ store.board.players.find(p => p.id === store.winnerId)?.name }}</p>
        <button class="btn-primary" @click="leaveRoom">返回大厅</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.game-container {
  width: 100%;
  height: 100%;
  background: var(--bg-color);
  color: white;
}

.waiting-room, .game-over-modal {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.glass-panel {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 40px;
  border-radius: 20px;
  width: 400px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.players-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.player-item {
  background: rgba(0, 0, 0, 0.2);
  padding: 12px 16px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.color-dot {
  width: 16px;
  height: 16px;
  border-radius: 50%;
}

.red .color-dot { background-color: var(--red-color); }
.yellow .color-dot { background-color: var(--yellow-color); }
.blue .color-dot { background-color: var(--blue-color); }
.green .color-dot { background-color: var(--green-color); }

.actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 20px;
}

.playing-area {
  display: flex;
  width: 100%;
  height: 100%;
  padding: 20px;
  box-sizing: border-box;
  gap: 20px;
  justify-content: center;
}

.board-container {
  flex: 0 0 auto;
  aspect-ratio: 1/1;
  height: calc(100vh - 40px);
  max-width: calc(100vh - 40px);
  background: rgba(255, 255, 255, 0.02);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.05);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.left-panel, .right-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.players-status {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.status-card {
  padding: 15px;
  border-radius: 12px;
  background: rgba(0, 0, 0, 0.3);
  border: 2px solid transparent;
  transition: all 0.3s;
}

.status-card.active {
  transform: scale(1.05);
  box-shadow: 0 0 20px rgba(255, 255, 255, 0.1);
}

.status-card.red.active { border-color: var(--red-color); box-shadow: 0 0 20px rgba(239, 68, 68, 0.3); }
.status-card.yellow.active { border-color: var(--yellow-color); box-shadow: 0 0 20px rgba(234, 179, 8, 0.3); }
.status-card.blue.active { border-color: var(--blue-color); box-shadow: 0 0 20px rgba(59, 130, 246, 0.3); }
.status-card.green.active { border-color: var(--green-color); box-shadow: 0 0 20px rgba(34, 197, 94, 0.3); }

.status-card .name { font-weight: bold; font-size: 1.1rem; }
.status-card .pieces-left { font-size: 0.9rem; color: #94a3b8; margin-top: 5px; }

</style>
