<script setup lang="ts">
import { computed } from 'vue'
import type { Card } from '../types/game'
import { Suit, Rank, suitSymbols, rankSymbols, getCardColor } from '../types/game'

const props = defineProps<{
  card: Card
  selected?: boolean
  disabled?: boolean
  small?: boolean
}>()

const emit = defineEmits<{
  click: []
}>()

const displayRank = computed(() => {
  if (props.card.Suit === Suit.Joker) {
    return props.card.Rank === Rank.RJ ? '大' : '小'
  }
  return rankSymbols[props.card.Rank as Rank]
})

const displaySuit = computed(() => {
  if (props.card.Suit === Suit.Joker) {
    return '王'
  }
  return suitSymbols[props.card.Suit as Suit]
})

const color = computed(() => getCardColor(props.card))

const isJoker = computed(() => props.card.Suit === Suit.Joker)

const handleClick = () => {
  if (!props.disabled) {
    emit('click')
  }
}
</script>

<template>
  <div
    class="poker-card"
    :class="{
      selected: selected,
      disabled: disabled,
      small: small,
      joker: isJoker,
      red: color === '#ff0000',
      black: color !== '#ff0000',
    }"
    @click="handleClick"
  >
    <div class="card-inner">
      <div class="card-corner top-left">
        <div class="corner-rank">{{ displayRank }}</div>
        <div class="corner-suit">{{ displaySuit }}</div>
      </div>
      <div class="card-center">
        <div class="center-suit">{{ displaySuit }}</div>
      </div>
      <div class="card-corner bottom-right">
        <div class="corner-rank">{{ displayRank }}</div>
        <div class="corner-suit">{{ displaySuit }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.poker-card {
  width: 70px;
  height: 100px;
  background: linear-gradient(135deg, #fff 0%, #f8f9fa 100%);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2), 0 1px 2px rgba(0, 0, 0, 0.1);
  user-select: none;
  border: 2px solid #e9ecef;
  overflow: hidden;
}

.poker-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.8) 0%, rgba(255, 255, 255, 0) 100%);
  pointer-events: none;
}

.poker-card:hover:not(.disabled) {
  transform: translateY(-12px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.3), 0 4px 8px rgba(0, 0, 0, 0.2);
  border-color: #adb5bd;
}

.poker-card.selected {
  transform: translateY(-24px);
  box-shadow: 0 16px 32px rgba(0, 150, 255, 0.5), 0 4px 8px rgba(0, 150, 255, 0.3);
  border-color: #4dabf7;
  background: linear-gradient(135deg, #e7f5ff 0%, #d0ebff 100%);
}

.poker-card.disabled {
  opacity: 0.7;
  cursor: default;
}

.poker-card.small {
  width: 50px;
  height: 72px;
  border-radius: 8px;
}

.poker-card.small .corner-rank {
  font-size: 12px;
}

.poker-card.small .corner-suit {
  font-size: 10px;
}

.poker-card.small .center-suit {
  font-size: 20px;
}

.poker-card.small .card-corner {
  padding: 4px;
}

.card-inner {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  position: relative;
  z-index: 1;
}

.card-corner {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 6px;
  line-height: 1;
}

.top-left {
  align-self: flex-start;
}

.bottom-right {
  align-self: flex-end;
  transform: rotate(180deg);
}

.corner-rank {
  font-size: 16px;
  font-weight: 700;
  line-height: 1;
}

.corner-suit {
  font-size: 12px;
  line-height: 1;
  margin-top: 2px;
}

.card-center {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.center-suit {
  font-size: 28px;
  opacity: 0.8;
}

.red {
  color: #e03131;
}

.red .card-inner {
  text-shadow: 0 1px 2px rgba(224, 49, 49, 0.2);
}

.black {
  color: #1a1a1a;
}

.joker {
  background: linear-gradient(135deg, #fff9db 0%, #fff3bf 100%);
  border-color: #ffd43b;
}

.joker.red {
  background: linear-gradient(135deg, #ffe8e8 0%, #ffc9c9 100%);
  border-color: #ff8787;
}

.joker.black {
  background: linear-gradient(135deg, #e7f5ff 0%, #d0ebff 100%);
  border-color: #74c0fc;
}

.joker .center-suit {
  font-size: 32px;
}
</style>
