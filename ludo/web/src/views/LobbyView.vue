<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGameStore } from '../stores/gameStore'
import { useWebSocket } from '../composables/useWebSocket'
import type { RoomJoinedPayload } from '../types/game'

const router = useRouter()
const store = useGameStore()
const { isConnected, connect, send, onMessage, removeHandler } = useWebSocket()

const playerName = ref(`Player_${Math.floor(Math.random() * 1000)}`)
const roomIdToJoin = ref('')

onMounted(() => {
  connect()
  onMessage('room_joined', onRoomJoined)
})

onUnmounted(() => {
  removeHandler('room_joined', onRoomJoined)
})

const onRoomJoined = (payload: RoomJoinedPayload) => {
  store.roomId = payload.roomId
  store.myId = payload.myId
  store.myColor = payload.myColor
  store.board.players = payload.players
  store.gameState = 'waiting'
  router.push('/game')
}

const createSinglePlayer = () => {
  send('create_room', {
    playerName: playerName.value,
    gameMode: 'single',
    maxPlayers: 4,
    aiCount: 3
  })
}

const createLanGame = () => {
  send('create_room', {
    playerName: playerName.value,
    gameMode: 'lan',
    maxPlayers: 4,
    aiCount: 0 // Optional AI fill later
  })
}

const joinGame = () => {
  if (!roomIdToJoin.value) return
  send('join_room', {
    roomId: roomIdToJoin.value,
    playerName: playerName.value
  })
}
</script>

<template>
  <div class="lobby-container">
    <div class="glass-panel">
      <h1 class="title">✈️ 飞行棋 Ludo</h1>
      
      <div class="form-group">
        <label>玩家昵称</label>
        <input v-model="playerName" type="text" class="input-field" />
      </div>

      <div class="actions">
        <button class="btn-primary" @click="createSinglePlayer" :disabled="!isConnected">单机游戏 (1v3 AI)</button>
        <button class="btn-secondary" @click="createLanGame" :disabled="!isConnected">创建联机房间</button>
      </div>

      <div class="divider">
        <span>或</span>
      </div>

      <div class="join-group">
        <input v-model="roomIdToJoin" type="text" placeholder="输入房间码" class="input-field" />
        <button class="btn-join" @click="joinGame" :disabled="!roomIdToJoin || !isConnected">加入房间</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.lobby-container {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  background: radial-gradient(circle at center, #1e293b 0%, #0f172a 100%);
}

.glass-panel {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 40px;
  border-radius: 20px;
  width: 350px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

.title {
  text-align: center;
  margin: 0 0 10px 0;
  font-size: 2rem;
  font-weight: 800;
  background: -webkit-linear-gradient(45deg, #a855f7, #3b82f6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

label {
  font-size: 0.9rem;
  color: #94a3b8;
}

.input-field {
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: white;
  padding: 12px;
  border-radius: 8px;
  font-size: 1rem;
  outline: none;
  transition: all 0.3s;
}

.input-field:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(168, 85, 247, 0.2);
}

.actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.btn-primary {
  padding: 14px;
}

.btn-secondary {
  background-color: transparent;
  border: 1px solid var(--primary-color);
  color: var(--primary-color);
  padding: 14px;
}

.btn-secondary:hover {
  background-color: rgba(168, 85, 247, 0.1);
}

.divider {
  display: flex;
  align-items: center;
  text-align: center;
  color: #64748b;
  font-size: 0.9rem;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.divider span {
  padding: 0 10px;
}

.join-group {
  display: flex;
  gap: 10px;
}

.join-group .input-field {
  flex: 1;
}

.btn-join {
  background-color: var(--blue-color);
}

.btn-join:hover {
  background-color: #2563eb;
}
</style>
