<script setup lang="ts">
import { computed } from 'vue'
import type { Card } from '../types/game'
import PokerCard from './PokerCard.vue'

const props = defineProps<{
  cards: Card[]
  selectedCards?: Card[]
  isMe?: boolean
  canSelect?: boolean
}>()

const emit = defineEmits<{
  select: [card: Card]
}>()

// 按点数排序（从小到大）
const sortedCards = computed(() => {
  return [...props.cards].sort((a, b) => {
    if (a.Rank !== b.Rank) {
      return a.Rank - b.Rank
    }
    return a.Suit - b.Suit
  })
})

const isSelected = (card: Card) => {
  if (!props.selectedCards) return false
  return props.selectedCards.some(
    (c) => c.Suit === card.Suit && c.Rank === card.Rank
  )
}

const handleCardClick = (card: Card) => {
  if (props.canSelect) {
    emit('select', card)
  }
}
</script>

<template>
  <div class="player-hand" :class="{ 'is-me': isMe }">
    <div class="cards-container">
      <PokerCard
        v-for="(card, index) in sortedCards"
        :key="`${card.Suit}-${card.Rank}-${index}`"
        :card="card"
        :selected="isSelected(card)"
        :disabled="!canSelect"
        :style="{ marginLeft: index > 0 ? '-20px' : '0', zIndex: index }"
        @click="handleCardClick(card)"
      />
    </div>
    <div class="card-count-badge" v-if="!isMe">
      <span class="count-number">{{ cards.length }}</span>
      <span class="count-label">张</span>
    </div>
  </div>
</template>

<style scoped>
.player-hand {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.cards-container {
  display: flex;
  justify-content: center;
  min-height: 110px;
  padding: 0 24px;
}

.is-me .cards-container {
  min-height: 130px;
  padding: 0 32px;
}

.card-count-badge {
  display: flex;
  align-items: baseline;
  gap: 4px;
  padding: 6px 16px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.count-number {
  color: #ffd700;
  font-size: 20px;
  font-weight: 700;
}

.count-label {
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
}
</style>
