import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { BoardState, Color } from '../types/game';

export const useGameStore = defineStore('game', () => {
  const roomId = ref('');
  const myId = ref('');
  const myColor = ref<Color | null>(null);
  const gameState = ref<'lobby' | 'waiting' | 'playing' | 'game_over'>('lobby');
  
  const board = ref<BoardState>({ players: [] });
  const currentTurn = ref<Color | null>(null);
  
  const diceValue = ref(0);
  const isRolling = ref(false);
  const movablePieces = ref<number[]>([]);
  const continuousSix = ref(0);

  const winnerId = ref('');

  return {
    roomId,
    myId,
    myColor,
    gameState,
    board,
    currentTurn,
    diceValue,
    isRolling,
    movablePieces,
    continuousSix,
    winnerId
  };
});
