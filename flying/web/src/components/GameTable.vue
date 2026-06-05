<script setup lang="ts">
import { computed, ref } from 'vue'
import { useGameStore } from '../stores/game'
import type { Card, PlayResultData } from '../types/game'
import PokerCard from './PokerCard.vue'
import PlayerHand from './PlayerHand.vue'
import CardCounter from './CardCounter.vue'

const store = useGameStore()

// 出牌记录
const playRecords = ref<{playerId: string, playerName: string, cards: Card[], isPass: boolean}[]>([])

// 所有已出的牌（用于记牌器）
const allPlayedCards = ref<Card[]>([])

// 获取当前玩家索引
const currentPlayerIndex = computed(() => {
  return store.players.findIndex((p) => p.id === store.playerId)
})

// 左边玩家（下家）
const leftPlayer = computed(() => {
  if (currentPlayerIndex.value < 0) return null
  const idx = (currentPlayerIndex.value + 1) % 3
  return store.players[idx]
})

// 右边玩家（上家）
const rightPlayer = computed(() => {
  if (currentPlayerIndex.value < 0) return null
  const idx = (currentPlayerIndex.value + 2) % 3
  return store.players[idx]
})

// 我的信息
const myPlayer = computed(() => {
  if (currentPlayerIndex.value < 0) return null
  return store.players[currentPlayerIndex.value]
})

// 是否轮到我
const isMyTurn = computed(() => {
  if (!myPlayer.value) return false
  return store.currentTurn === currentPlayerIndex.value
})

// 是否可以出牌
const canPlay = computed(() => {
  return store.isPlaying && isMyTurn.value
})

// 是否可以叫地主
const canBid = computed(() => {
  return store.isBidding && isMyTurn.value
})

// 左边玩家是否轮到
const isLeftTurn = computed(() => {
  if (currentPlayerIndex.value < 0) return false
  return store.currentTurn === (currentPlayerIndex.value + 1) % 3
})

// 右边玩家是否轮到
const isRightTurn = computed(() => {
  if (currentPlayerIndex.value < 0) return false
  return store.currentTurn === (currentPlayerIndex.value + 2) % 3
})

// 右边玩家的出牌记录
const rightPlayerLastPlay = computed(() => {
  const record = playRecords.value.find(r => r.playerId === rightPlayer.value?.id)
  return record || null
})

// 左边玩家的出牌记录
const leftPlayerLastPlay = computed(() => {
  const record = playRecords.value.find(r => r.playerId === leftPlayer.value?.id)
  return record || null
})

// 我的出牌记录
const myLastPlay = computed(() => {
  const record = playRecords.value.find(r => r.playerId === store.playerId)
  return record || null
})

const handleCreateRoom = () => {
  store.createRoom()
}

const handleBid = (score: number) => {
  store.bid(score)
}

const handlePlay = () => {
  store.playCards()
}

const handlePass = () => {
  store.pass()
}

const handleHint = () => {
  store.getHint()
}

const handleCardSelect = (card: any) => {
  if (canPlay.value) {
    store.toggleCard(card)
  }
}

const handleRestart = () => {
  store.resetGame()
  playRecords.value = []
  allPlayedCards.value = []
}

// 监听出牌结果，更新出牌记录
const handlePlayResult = (data: PlayResultData) => {
  // 添加到出牌记录
  playRecords.value.unshift({
    playerId: data.player_id,
    playerName: data.player_name,
    cards: data.cards || [],
    isPass: data.is_pass
  })

  // 添加到所有已出的牌（用于记牌器）
  if (!data.is_pass && data.cards) {
    allPlayedCards.value.push(...data.cards)
  }

  // 只保留最近3条记录
  if (playRecords.value.length > 3) {
    playRecords.value = playRecords.value.slice(0, 3)
  }
}

// 监听游戏状态变化，清空出牌记录
const clearPlayRecords = () => {
  playRecords.value = []
  allPlayedCards.value = []
}

// 暴露给父组件
defineExpose({
  handlePlayResult,
  clearPlayRecords
})
</script>

<template>
  <div class="game-table">
    <!-- 等待界面 -->
    <div v-if="!store.gameStarted" class="waiting-screen">
      <div class="waiting-content">
        <div class="logo-container">
          <div class="logo-icon">🃏</div>
          <h1 class="title">斗地主</h1>
        </div>
        <p class="subtitle">经典单机版 · 智能AI对战</p>
        <button class="btn-start" @click="handleCreateRoom">
          <span class="btn-text">开始游戏</span>
          <span class="btn-icon">→</span>
        </button>
      </div>
      <div class="waiting-decoration">
        <div class="decoration-card card-1">♠</div>
        <div class="decoration-card card-2">♥</div>
        <div class="decoration-card card-3">♣</div>
        <div class="decoration-card card-4">♦</div>
      </div>
    </div>

    <!-- 游戏界面 -->
    <template v-else>
      <!-- 顶部信息栏 -->
      <div class="game-header">
        <div class="header-left">
          <div class="room-badge">
            <span class="badge-icon">🏠</span>
            <span class="badge-text">{{ store.roomId }}</span>
          </div>
        </div>
        <div class="header-center">
          <div class="game-state-badge" :class="{
            'state-bidding': store.isBidding,
            'state-playing': store.isPlaying,
            'state-gameover': store.isGameOver
          }">
            <template v-if="store.isBidding">🎯 叫地主阶段</template>
            <template v-else-if="store.isPlaying">🎮 出牌阶段</template>
            <template v-else-if="store.isGameOver">🏆 游戏结束</template>
          </div>
        </div>
        <div class="header-right">
          <div class="score-badge" v-if="store.bidScore > 0">
            <span class="badge-icon">💰</span>
            <span class="badge-text">底分: {{ store.bidScore }}</span>
          </div>
        </div>
      </div>

      <!-- 记牌器 - 正中上方 -->
      <div class="card-counter-container">
        <CardCounter :played-cards="allPlayedCards" />
      </div>

      <!-- 游戏桌面 -->
      <div class="table-content">
        <!-- 右边玩家（上家） -->
        <div class="player-area right-player" :class="{ 'is-turn': isRightTurn }">
          <div class="player-card">
            <div class="player-avatar">
              <div class="avatar-icon">🤖</div>
              <div class="turn-indicator" v-if="isRightTurn"></div>
            </div>
            <div class="player-info">
              <div class="player-name">{{ rightPlayer?.name || '等待中...' }}</div>
              <div class="player-role-badge" :class="{
                'role-landlord': rightPlayer?.status === 3,
                'role-farmer': rightPlayer?.status === 4
              }">
                <template v-if="rightPlayer?.status === 3">👑 地主</template>
                <template v-else-if="rightPlayer?.status === 4">🌾 农民</template>
              </div>
            </div>
            <div class="card-count-badge">
              <span class="count-number">{{ rightPlayer?.card_count || 0 }}</span>
              <span class="count-label">张</span>
            </div>
          </div>
          <!-- 出牌记录 -->
          <div class="play-record-area" v-if="rightPlayerLastPlay">
            <div class="record-content" v-if="!rightPlayerLastPlay.isPass">
              <PokerCard
                v-for="(card, index) in rightPlayerLastPlay.cards"
                :key="`${card.Suit}-${card.Rank}-${index}`"
                :card="card"
                :small="true"
              />
            </div>
            <div class="record-pass" v-else>不出</div>
          </div>
          <div class="thinking-bubble" v-if="isRightTurn && store.isBidding">
            <span>思考中...</span>
          </div>
        </div>

        <!-- 中央区域 -->
        <div class="center-area">
          <!-- 底牌展示区 -->
          <div class="bottom-cards-area" v-if="store.isPlaying || store.isGameOver">
            <div class="area-label">底牌</div>
            <div class="bottom-cards-row">
              <!-- 这里应该显示底牌 -->
            </div>
          </div>

          <!-- 等待提示 -->
          <div class="waiting-indicator" v-if="(store.isBidding && !isMyTurn) || (store.isPlaying && !isMyTurn)">
            <div class="indicator-dot"></div>
            <span v-if="store.isBidding">等待其他玩家叫地主...</span>
            <span v-else>等待其他玩家出牌...</span>
          </div>

          <!-- 游戏结束 -->
          <div class="game-over-panel" v-if="store.isGameOver && store.gameOverData">
            <div class="panel-header">
              <div class="trophy-icon">🏆</div>
              <h2>游戏结束</h2>
            </div>
            <div class="panel-body">
              <p class="winner-name">{{ store.gameOverData.winner_name }} 获胜！</p>
              <p class="winner-role">{{ store.gameOverData.is_landlord ? '👑 地主' : '🌾 农民' }} 胜利</p>
              <p class="winner-score">得分: <span class="score-value">{{ store.gameOverData.score }}</span></p>
            </div>
            <button class="btn-restart" @click="handleRestart">
              <span>再来一局</span>
            </button>
          </div>
        </div>

        <!-- 左边玩家（下家） -->
        <div class="player-area left-player" :class="{ 'is-turn': isLeftTurn }">
          <div class="player-card">
            <div class="player-avatar">
              <div class="avatar-icon">🤖</div>
              <div class="turn-indicator" v-if="isLeftTurn"></div>
            </div>
            <div class="player-info">
              <div class="player-name">{{ leftPlayer?.name || '等待中...' }}</div>
              <div class="player-role-badge" :class="{
                'role-landlord': leftPlayer?.status === 3,
                'role-farmer': leftPlayer?.status === 4
              }">
                <template v-if="leftPlayer?.status === 3">👑 地主</template>
                <template v-else-if="leftPlayer?.status === 4">🌾 农民</template>
              </div>
            </div>
            <div class="card-count-badge">
              <span class="count-number">{{ leftPlayer?.card_count || 0 }}</span>
              <span class="count-label">张</span>
            </div>
          </div>
          <!-- 出牌记录 -->
          <div class="play-record-area" v-if="leftPlayerLastPlay">
            <div class="record-content" v-if="!leftPlayerLastPlay.isPass">
              <PokerCard
                v-for="(card, index) in leftPlayerLastPlay.cards"
                :key="`${card.Suit}-${card.Rank}-${index}`"
                :card="card"
                :small="true"
              />
            </div>
            <div class="record-pass" v-else>不出</div>
          </div>
          <div class="thinking-bubble" v-if="isLeftTurn && store.isBidding">
            <span>思考中...</span>
          </div>
        </div>
      </div>

      <!-- 底部玩家区域 -->
      <div class="my-area">
        <!-- 我的信息 -->
        <div class="my-info-bar">
          <div class="my-avatar">
            <div class="avatar-icon player-avatar-icon">😊</div>
          </div>
          <div class="my-details">
            <div class="my-name">{{ myPlayer?.name || '你' }}</div>
            <div class="my-role-badge" :class="{
              'role-landlord': store.isLandlord,
              'role-farmer': !store.isLandlord && store.isPlaying
            }">
              <template v-if="store.isLandlord">👑 地主</template>
              <template v-else-if="store.isPlaying">🌾 农民</template>
            </div>
          </div>
          <div class="my-card-count">
            <span class="count-number">{{ store.myCards.length }}</span>
            <span class="count-label">张</span>
          </div>
        </div>

        <!-- 我的出牌记录 -->
        <div class="my-play-record" v-if="myLastPlay">
          <div class="record-content" v-if="!myLastPlay.isPass">
            <PokerCard
              v-for="(card, index) in myLastPlay.cards"
              :key="`${card.Suit}-${card.Rank}-${index}`"
              :card="card"
              :small="true"
            />
          </div>
          <div class="record-pass" v-else>不出</div>
        </div>

        <!-- 我的手牌 -->
        <PlayerHand
          :cards="store.myCards"
          :selected-cards="store.selectedCards"
          :is-me="true"
          :can-select="canPlay"
          @select="handleCardSelect"
        />

        <!-- 操作按钮 -->
        <div class="actions-bar">
          <!-- 叫地主按钮 -->
          <template v-if="canBid">
            <button class="btn-action btn-pass" @click="handleBid(0)">不叫</button>
            <button
              class="btn-action btn-bid"
              v-if="store.bidScore < 1"
              @click="handleBid(1)"
            >
              1分
            </button>
            <button
              class="btn-action btn-bid"
              v-if="store.bidScore < 2"
              @click="handleBid(2)"
            >
              2分
            </button>
            <button
              class="btn-action btn-bid"
              v-if="store.bidScore < 3"
              @click="handleBid(3)"
            >
              3分
            </button>
          </template>

          <!-- 出牌按钮 -->
          <template v-if="canPlay">
            <button class="btn-action btn-pass" @click="handlePass" :disabled="!store.canPass">
              不出
            </button>
            <button
              class="btn-action btn-play"
              @click="handlePlay"
              :disabled="store.selectedCards.length === 0"
            >
              出牌
            </button>
            <button class="btn-action btn-hint" @click="handleHint">提示</button>
          </template>
        </div>
      </div>

      <!-- 错误提示 -->
      <div class="error-toast" v-if="store.error">
        <div class="toast-icon">⚠️</div>
        <div class="toast-message">{{ store.error }}</div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.game-table {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #0c1220 0%, #1a2a3a 50%, #0c1220 100%);
  position: relative;
  overflow: hidden;
}

/* 等待界面 */
.waiting-screen {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
}

.waiting-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 30px;
  z-index: 2;
}

.logo-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.logo-icon {
  font-size: 80px;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-20px); }
}

.title {
  font-size: 64px;
  background: linear-gradient(135deg, #ffd700 0%, #ff8c00 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-shadow: none;
  font-weight: 800;
  letter-spacing: 8px;
}

.subtitle {
  font-size: 20px;
  color: rgba(255, 255, 255, 0.7);
  letter-spacing: 4px;
}

.btn-start {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 18px 48px;
  background: linear-gradient(135deg, #ffd700 0%, #ff8c00 100%);
  border: none;
  border-radius: 50px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 8px 30px rgba(255, 215, 0, 0.4);
}

.btn-start:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(255, 215, 0, 0.6);
}

.btn-text {
  font-size: 20px;
  font-weight: 700;
  color: #1a1a2e;
}

.btn-icon {
  font-size: 24px;
  color: #1a1a2e;
  transition: transform 0.3s ease;
}

.btn-start:hover .btn-icon {
  transform: translateX(4px);
}

.waiting-decoration {
  position: absolute;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.decoration-card {
  position: absolute;
  font-size: 120px;
  opacity: 0.1;
  animation: rotate 20s linear infinite;
}

.card-1 { top: 10%; left: 10%; color: #1e3a5f; animation-delay: 0s; }
.card-2 { top: 20%; right: 15%; color: #5f1e3a; animation-delay: 5s; }
.card-3 { bottom: 20%; left: 20%; color: #3a5f1e; animation-delay: 10s; }
.card-4 { bottom: 10%; right: 10%; color: #5f3a1e; animation-delay: 15s; }

@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 记牌器容器 */
.card-counter-container {
  display: flex;
  justify-content: center;
  padding: 8px 0;
  z-index: 10;
}

/* 游戏头部 */
.game-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 24px;
  background: linear-gradient(180deg, rgba(0, 0, 0, 0.6) 0%, rgba(0, 0, 0, 0.3) 100%);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.header-left, .header-right {
  flex: 1;
}

.header-center {
  flex: 2;
  display: flex;
  justify-content: center;
}

.room-badge, .score-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.badge-icon {
  font-size: 14px;
}

.badge-text {
  color: rgba(255, 255, 255, 0.8);
  font-size: 13px;
}

.game-state-badge {
  padding: 8px 24px;
  border-radius: 20px;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 2px;
}

.state-bidding {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 20px rgba(102, 126, 234, 0.4);
}

.state-playing {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: #000;
  box-shadow: 0 4px 20px rgba(79, 172, 254, 0.4);
}

.state-gameover {
  background: linear-gradient(135deg, #ffd700 0%, #ff8c00 100%);
  color: #000;
  box-shadow: 0 4px 20px rgba(255, 215, 0, 0.4);
}

/* 游戏桌面 */
.table-content {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 80px;
}

/* 玩家区域 */
.player-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  min-width: 160px;
}

.player-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.player-area.is-turn .player-card {
  border-color: rgba(255, 215, 0, 0.5);
  box-shadow: 0 0 30px rgba(255, 215, 0, 0.2);
}

.player-avatar {
  position: relative;
}

.avatar-icon {
  font-size: 48px;
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.1) 0%, rgba(255, 255, 255, 0.05) 100%);
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.2);
}

.player-avatar-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-color: rgba(102, 126, 234, 0.5);
}

.turn-indicator {
  position: absolute;
  bottom: -4px;
  right: -4px;
  width: 16px;
  height: 16px;
  background: #ffd700;
  border-radius: 50%;
  border: 3px solid #1a2a3a;
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.2); opacity: 0.8; }
}

.player-info {
  text-align: center;
}

.player-name {
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 6px;
}

.player-role-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.role-landlord {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
  color: white;
}

.role-farmer {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  color: #000;
}

.card-count-badge {
  display: flex;
  align-items: baseline;
  gap: 4px;
  padding: 6px 14px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
}

.count-number {
  color: #ffd700;
  font-size: 24px;
  font-weight: 700;
}

.count-label {
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
}

/* 出牌记录区域 */
.play-record-area {
  display: flex;
  justify-content: center;
  min-height: 80px;
  padding: 8px;
}

.record-content {
  display: flex;
  gap: 4px;
}

.record-pass {
  color: rgba(255, 255, 255, 0.5);
  font-size: 16px;
  font-style: italic;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
}

.thinking-bubble {
  padding: 8px 16px;
  background: rgba(255, 215, 0, 0.1);
  border: 1px solid rgba(255, 215, 0, 0.3);
  border-radius: 20px;
  color: #ffd700;
  font-size: 14px;
  animation: pulse 1.5s infinite;
}

/* 中央区域 */
.center-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  min-height: 200px;
}

.bottom-cards-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.area-label {
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
  letter-spacing: 2px;
}

.bottom-cards-row {
  display: flex;
  gap: 8px;
}

.waiting-indicator {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 24px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 16px;
}

.indicator-dot {
  width: 10px;
  height: 10px;
  background: #ffd700;
  border-radius: 50%;
  animation: pulse 1.5s infinite;
}

.game-over-panel {
  background: linear-gradient(135deg, rgba(15, 32, 39, 0.95) 0%, rgba(32, 58, 67, 0.95) 100%);
  padding: 40px;
  border-radius: 24px;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.panel-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.trophy-icon {
  font-size: 64px;
  animation: float 2s ease-in-out infinite;
}

.panel-header h2 {
  font-size: 32px;
  background: linear-gradient(135deg, #ffd700 0%, #ff8c00 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
}

.panel-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 32px;
}

.winner-name {
  font-size: 24px;
  color: #fff;
  font-weight: 600;
  margin: 0;
}

.winner-role {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
}

.winner-score {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
}

.score-value {
  color: #ffd700;
  font-weight: 700;
  font-size: 24px;
}

.btn-restart {
  padding: 14px 36px;
  background: linear-gradient(135deg, #ffd700 0%, #ff8c00 100%);
  border: none;
  border-radius: 50px;
  font-size: 18px;
  font-weight: 700;
  color: #1a1a2e;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 20px rgba(255, 215, 0, 0.4);
}

.btn-restart:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(255, 215, 0, 0.6);
}

/* 底部玩家区域 */
.my-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 20px 24px;
  background: linear-gradient(0deg, rgba(0, 0, 0, 0.6) 0%, rgba(0, 0, 0, 0.3) 100%);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.my-info-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 10px 24px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.my-avatar .avatar-icon {
  width: 48px;
  height: 48px;
  font-size: 32px;
}

.my-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.my-name {
  color: #fff;
  font-size: 18px;
  font-weight: 600;
}

.my-role-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.my-card-count {
  display: flex;
  align-items: baseline;
  gap: 4px;
  margin-left: auto;
}

.my-play-record {
  display: flex;
  justify-content: center;
  min-height: 60px;
  padding: 8px;
}

/* 操作按钮 */
.actions-bar {
  display: flex;
  gap: 16px;
  min-height: 52px;
}

.btn-action {
  padding: 12px 32px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  letter-spacing: 1px;
}

.btn-action:hover:not(:disabled) {
  transform: translateY(-2px);
}

.btn-action:active:not(:disabled) {
  transform: translateY(0);
}

.btn-action:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-pass {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-pass:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.2);
}

.btn-bid {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.btn-bid:hover:not(:disabled) {
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

.btn-play {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  color: #000;
  box-shadow: 0 4px 15px rgba(67, 233, 123, 0.4);
}

.btn-play:hover:not(:disabled) {
  box-shadow: 0 6px 20px rgba(67, 233, 123, 0.6);
}

.btn-hint {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: #000;
  box-shadow: 0 4px 15px rgba(79, 172, 254, 0.4);
}

.btn-hint:hover:not(:disabled) {
  box-shadow: 0 6px 20px rgba(79, 172, 254, 0.6);
}

/* 错误提示 */
.error-toast {
  position: fixed;
  top: 80px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 24px;
  background: linear-gradient(135deg, #ff4757 0%, #ff6b81 100%);
  border-radius: 12px;
  animation: slideDown 0.3s ease;
  box-shadow: 0 8px 30px rgba(255, 71, 87, 0.4);
  z-index: 100;
}

.toast-icon {
  font-size: 20px;
}

.toast-message {
  color: white;
  font-size: 16px;
  font-weight: 500;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}
</style>
