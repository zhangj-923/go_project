<script setup lang="ts">
import { computed } from 'vue'
import { useGameStore } from '../stores/gameStore'
import { useWebSocket } from '../composables/useWebSocket'

const store = useGameStore()
const ws = useWebSocket()

const canRoll = computed(() => {
  return store.currentTurn === store.myColor && !store.isRolling && store.diceValue === 0
})

const rollDice = () => {
  if (canRoll.value) {
    ws.send('roll_dice')
  }
}

// 3x3 Grid Dots mapping for standard 1-6 dice
const isDotActive = (index: number) => {
  const val = store.diceValue
  if (val === 1) return index === 5
  if (val === 2) return index === 1 || index === 9
  if (val === 3) return index === 1 || index === 5 || index === 9
  if (val === 4) return index === 1 || index === 3 || index === 7 || index === 9
  if (val === 5) return index === 1 || index === 3 || index === 5 || index === 7 || index === 9
  if (val === 6) return index === 1 || index === 3 || index === 4 || index === 6 || index === 7 || index === 9
  return false
}
</script>

<template>
  <div class="dice-panel">
    <div class="glass-card">
      <h3>骰子与操作</h3>
      
      <div class="dice-container" :class="{ rolling: store.isRolling }" @click="rollDice">
        <!-- Waiting/Rolling Placeholder -->
        <div v-if="store.diceValue === 0" class="dice placeholder">
          <span v-if="canRoll">点击掷骰</span>
          <span v-else>等待对方</span>
        </div>
        
        <!-- Standard Green Dice with White Dots -->
        <div v-else class="dice result green-dice">
          <div class="dice-grid">
            <div 
              v-for="index in 9" 
              :key="index" 
              class="dot-slot" 
              :class="{ active: isDotActive(index) }"
            ></div>
          </div>
        </div>
      </div>
      
      <div class="turn-info" :class="store.currentTurn">
        当前回合: {{ store.board.players.find(p => p.color === store.currentTurn)?.name }}
      </div>
      
      <div class="events" v-if="store.continuousSix > 0">
        连续抛出6点: {{ store.continuousSix }} 次! (3次将返航)
      </div>
    </div>
  </div>
</template>

<style scoped>
.glass-card {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px;
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

h3 {
  margin: 0;
  color: #f8fafc;
}

.dice-container {
  width: 100px;
  height: 100px;
  perspective: 1000px;
  cursor: pointer;
}

.dice {
  width: 100%;
  height: 100%;
  border-radius: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  transition: transform 0.6s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.dice.placeholder {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.5);
  font-size: 1rem;
  box-shadow: none;
  border: 2px dashed rgba(255, 255, 255, 0.2);
}

/* Premium Green Dice Styles */
.green-dice {
  background: radial-gradient(circle at 30% 30%, #2ecc71, #27ae60);
  box-shadow: inset 0 0 12px rgba(0, 0, 0, 0.4), 
              0 10px 20px rgba(0, 0, 0, 0.3), 
              0 2px 2px rgba(255, 255, 255, 0.2);
  border: 2px solid #2ecc71;
}

.dice-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(3, 1fr);
  width: 66px;
  height: 66px;
  gap: 5px;
  padding: 5px;
}

.dot-slot {
  display: flex;
  justify-content: center;
  align-items: center;
}

.dot-slot.active::after {
  content: '';
  width: 12px;
  height: 12px;
  background-color: white;
  border-radius: 50%;
  box-shadow: inset 1px 1px 2px rgba(0, 0, 0, 0.4), 
              0 1px 1px rgba(255, 255, 255, 0.5);
}

.dice-container.rolling .dice {
  animation: roll 0.6s ease-in-out infinite alternate;
}

@keyframes roll {
  0% { transform: rotateX(0deg) rotateY(0deg) scale(0.95); }
  100% { transform: rotateX(360deg) rotateY(360deg) scale(1.05); }
}

.turn-info {
  font-weight: bold;
  padding: 8px 16px;
  border-radius: 20px;
  background: rgba(0,0,0,0.3);
}

.turn-info.red { color: var(--red-color); }
.turn-info.yellow { color: var(--yellow-color); }
.turn-info.blue { color: var(--blue-color); }
.turn-info.green { color: var(--green-color); }

.events {
  color: var(--yellow-color);
  font-weight: bold;
  animation: pulse 1s infinite alternate;
}

@keyframes pulse {
  from { opacity: 0.7; transform: scale(1); }
  to { opacity: 1; transform: scale(1.05); }
}
</style>
