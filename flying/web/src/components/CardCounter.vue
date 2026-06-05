<script setup lang="ts">
import { computed } from 'vue'
import type { Card } from '../types/game'
import { Suit, Rank, rankSymbols } from '../types/game'

const props = defineProps<{
  playedCards: Card[]
}>()

// 统计每张牌出的数量
const cardCounts = computed(() => {
  const counts: Record<string, { total: number, played: number }> = {}

  // 初始化所有牌（除了大小王，按点数统计）
  const ranks = [Rank.R3, Rank.R4, Rank.R5, Rank.R6, Rank.R7, Rank.R8, Rank.R9, Rank.R10, Rank.J, Rank.Q, Rank.K, Rank.A, Rank.R2]
  ranks.forEach(rank => {
    counts[rank.toString()] = { total: 4, played: 0 }
  })

  // 大小王
  counts['bj'] = { total: 1, played: 0 }
  counts['rj'] = { total: 1, played: 0 }

  // 统计已出的牌
  props.playedCards.forEach(card => {
    if (card.Suit === Suit.Joker) {
      if (card.Rank === Rank.BJ) {
        counts['bj'].played++
      } else {
        counts['rj'].played++
      }
    } else {
      const key = card.Rank.toString()
      if (counts[key]) {
        counts[key].played++
      }
    }
  })

  return counts
})

// 获取牌的显示名称
const getRankName = (key: string): string => {
  if (key === 'bj') return '小王'
  if (key === 'rj') return '大王'
  const rank = parseInt(key) as Rank
  return rankSymbols[rank] || key
}

// 获取剩余数量
const getRemaining = (key: string): number => {
  const count = cardCounts.value[key]
  return count ? count.total - count.played : 0
}

// 获取显示数据
const displayData = computed(() => {
  const ranks = ['3', '4', '5', '6', '7', '8', '9', '10', 'j', 'q', 'k', 'a', '2', 'bj', 'rj']
  return ranks.map(key => ({
    key,
    name: getRankName(key),
    total: cardCounts.value[key]?.total || 0,
    played: cardCounts.value[key]?.played || 0,
    remaining: getRemaining(key)
  }))
})

// 获取剩余数量的颜色
const getCountColor = (remaining: number, total: number): string => {
  const ratio = remaining / total
  if (ratio === 0) return '#ff6b6b'
  if (ratio <= 0.25) return '#ffa502'
  if (ratio <= 0.5) return '#ffd700'
  return '#7bed9f'
}
</script>

<template>
  <div class="card-counter">
    <div class="counter-items">
      <div
        v-for="item in displayData"
        :key="item.key"
        class="counter-item"
        :class="{ 'is-zero': item.remaining === 0 }"
      >
        <div class="item-name">{{ item.name }}</div>
        <div class="item-count" :style="{ color: getCountColor(item.remaining, item.total) }">
          {{ item.remaining }}
        </div>
      </div>
    </div>
    <div class="counter-divider"></div>
    <div class="counter-total">
      <span class="total-label">已出</span>
      <span class="total-value">{{ playedCards.length }}</span>
      <span class="total-label">张</span>
    </div>
  </div>
</template>

<style scoped>
.card-counter {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 10px 20px;
  background: linear-gradient(135deg, rgba(15, 32, 39, 0.95) 0%, rgba(32, 58, 67, 0.95) 100%);
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3), inset 0 1px 0 rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

.counter-items {
  display: flex;
  align-items: center;
  gap: 4px;
}

.counter-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 6px 8px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  min-width: 36px;
  transition: all 0.2s ease;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.counter-item:hover {
  background: rgba(255, 255, 255, 0.1);
  transform: translateY(-2px);
}

.counter-item.is-zero {
  background: rgba(255, 107, 107, 0.15);
  border-color: rgba(255, 107, 107, 0.3);
}

.item-name {
  color: rgba(255, 255, 255, 0.7);
  font-size: 11px;
  font-weight: 500;
}

.item-count {
  font-size: 18px;
  font-weight: 700;
  text-shadow: 0 0 10px currentColor;
}

.counter-divider {
  width: 1px;
  height: 40px;
  background: rgba(255, 255, 255, 0.15);
  margin: 0 4px;
}

.counter-total {
  display: flex;
  align-items: baseline;
  gap: 4px;
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
}

.total-label {
  color: rgba(255, 255, 255, 0.6);
  font-size: 11px;
}

.total-value {
  color: #ffd700;
  font-size: 18px;
  font-weight: 700;
}
</style>
