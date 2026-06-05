<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useWebSocket } from './composables/useWebSocket'
import { useGameStore } from './stores/game'
import type { GameMessage, PlayResultData } from './types/game'
import GameTable from './components/GameTable.vue'

const { send } = useWebSocket()
const store = useGameStore()
const gameTableRef = ref<InstanceType<typeof GameTable> | null>(null)

// 设置发送函数
store.setSendFn(send)

// 监听游戏消息
const handleGameMessage = (event: Event) => {
  const msg = (event as CustomEvent).detail as GameMessage
  store.handleMessage(msg)

  // 如果是出牌结果，更新出牌记录
  if (msg.type === 'play_result' && gameTableRef.value) {
    gameTableRef.value.handlePlayResult(msg.data as PlayResultData)
  }

  // 如果是游戏开始或结束，清空出牌记录
  if ((msg.type === 'deal_cards' || msg.type === 'game_over') && gameTableRef.value) {
    gameTableRef.value.clearPlayRecords()
  }
}

onMounted(() => {
  window.addEventListener('game-message', handleGameMessage)
})

onUnmounted(() => {
  window.removeEventListener('game-message', handleGameMessage)
})
</script>

<template>
  <div class="app">
    <GameTable ref="gameTableRef" />
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background: #0c1220;
  min-height: 100vh;
  overflow: hidden;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.app {
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #0c1220 0%, #1a2a3a 50%, #0c1220 100%);
}

/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}

/* 选中文本样式 */
::selection {
  background: rgba(102, 126, 234, 0.5);
  color: #fff;
}
</style>
