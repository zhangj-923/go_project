<script setup lang="ts">
import { computed } from 'vue'
import type { Piece } from '../types/game'

const props = defineProps<{
  piece: Piece
  movable: boolean
  x: number
  y: number
  isStacked?: boolean
  stackIndex?: number
  stackTotal?: number
}>()

const emit = defineEmits(['move'])

// Offset for stacked pieces
const offsetX = computed(() => {
  if (!props.isStacked || !props.stackTotal) return 0
  if (props.stackTotal === 1) return 0
  const angle = (Math.PI * 2 * props.stackIndex!) / props.stackTotal
  return Math.cos(angle) * 15
})

const offsetY = computed(() => {
  if (!props.isStacked || !props.stackTotal) return 0
  if (props.stackTotal === 1) return 0
  const angle = (Math.PI * 2 * props.stackIndex!) / props.stackTotal
  return Math.sin(angle) * 15
})

const onClick = () => {
  if (props.movable) {
    emit('move', props.piece.id)
  }
}
</script>

<template>
  <g 
    class="piece-token" 
    :class="[piece.color, { movable }]"
    :style="{ transform: `translate(${x + offsetX}px, ${y + offsetY}px)` }"
    @click="onClick"
  >
    <!-- Drop Shadow -->
    <circle cx="0" cy="0" r="16" fill="rgba(0, 0, 0, 0.3)" transform="translate(1.5, 2.5)" />
    <!-- Token Outer Border (White Rim) -->
    <circle cx="0" cy="0" r="16" fill="#ffffff" stroke="#333333" stroke-width="1" />
    <!-- Inner Colored Body -->
    <circle cx="0" cy="0" r="13.5" class="token-body" />
    <!-- White Airplane Silhouette -->
    <path 
      d="M 0 -9 C -1.2 -9 -2 -6.5 -2 -3.5 L -8.5 -0.5 L -8.5 1.5 L -2 0 L -2 5.5 L -5 7.5 L -5 8.5 L 0 7 L 5 8.5 L 5 7.5 L 2 5.5 L 2 0 L 8.5 1.5 L 8.5 -0.5 L 2 -3.5 C 2 -6.5 1.2 -9 0 -9 Z" 
      fill="white" 
    />
    
    <!-- Highlight ring when movable -->
    <circle v-if="movable" cx="0" cy="0" r="22" class="highlight" />
    
    <text v-if="piece.position >= 100 && piece.position % 10 === 5" 
          x="0" y="3" text-anchor="middle" font-size="10" fill="white" font-weight="bold">
      ★
    </text>
  </g>
</template>

<style scoped>
.piece-token {
  transition: transform 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  cursor: default;
}

.piece-token.movable {
  cursor: pointer;
}

.red .token-body { fill: var(--red-color); }
.yellow .token-body { fill: var(--yellow-color); }
.blue .token-body { fill: var(--blue-color); }
.green .token-body { fill: var(--green-color); }

.highlight {
  fill: none;
  stroke: white;
  stroke-width: 2.5;
  stroke-dasharray: 4 4;
  animation: rotate 6s linear infinite;
  opacity: 0.9;
}

@keyframes rotate {
  100% { transform: rotate(360deg); }
}

.piece-token:hover.movable {
  transform: scale(1.15);
  z-index: 10;
}
</style>
