import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type {
  Card,
  RoomInfo,
  GameMessage,
  DealCardsData,
  BidResultData,
  PlayResultData,
  GameOverData,
  TurnData,
} from '../types/game'
import { GameState } from '../types/game'
import { analyzeHand, canBeat, getCardTypeName } from '../utils/cardUtils'

export const useGameStore = defineStore('game', () => {
  // 状态
  const playerId = ref<string>('')
  const roomId = ref<string>('')
  const roomInfo = ref<RoomInfo | null>(null)
  const myCards = ref<Card[]>([])
  const selectedCards = ref<Card[]>([])
  const isLandlord = ref(false)
  const lastPlayedCards = ref<Card[]>([])
  const lastPlayer = ref<string>('')
  const canPass = ref(false)
  const isMyTurn = ref(false)
  const gameOverData = ref<GameOverData | null>(null)
  const error = ref<string>('')
  const bidScore = ref(0)
  const gameStarted = ref(false)
  const lastHandResult = ref<{ type: number, mainRank: number, length: number } | null>(null)

  // 计算属性
  const isWaiting = computed(() => roomInfo.value?.state === GameState.Waiting)
  const isBidding = computed(() => roomInfo.value?.state === GameState.Bidding)
  const isPlaying = computed(() => roomInfo.value?.state === GameState.Playing)
  const isGameOver = computed(() => roomInfo.value?.state === GameState.GameOver)

  const players = computed(() => roomInfo.value?.players || [])
  const currentTurn = computed(() => roomInfo.value?.turn ?? -1)

  // WebSocket 发送函数
  let sendFn: ((msg: GameMessage) => void) | null = null

  const setSendFn = (fn: (msg: GameMessage) => void) => {
    sendFn = fn
  }

  const send = (msg: GameMessage) => {
    if (sendFn) {
      sendFn(msg)
    }
  }

  // 创建房间
  const createRoom = () => {
    resetGame()
    send({ type: 'create_room', data: {} })
  }

  // 开始游戏
  const startGame = () => {
    send({ type: 'start_game', data: {} })
  }

  // 叫地主
  const bid = (score: number) => {
    send({ type: 'bid', data: { score } })
  }

  // 出牌
  const playCards = () => {
    if (selectedCards.value.length === 0) {
      error.value = '请选择要出的牌'
      setTimeout(() => { error.value = '' }, 2000)
      return
    }

    // 检查牌型是否有效
    const result = analyzeHand(selectedCards.value)
    if (result.type === 0) {
      error.value = '无效的牌型'
      setTimeout(() => { error.value = '' }, 2000)
      return
    }

    // 如果需要压牌，检查是否能压过
    if (lastHandResult.value && lastPlayer.value !== playerId.value) {
      if (!canBeat(selectedCards.value, lastHandResult.value)) {
        const typeName = getCardTypeName(result.type)
        error.value = `${typeName}压不过，请选择更大的牌或不出`
        setTimeout(() => { error.value = '' }, 2000)
        return
      }
    }

    send({ type: 'play', data: selectedCards.value })
    selectedCards.value = []
  }

  // 不出
  const pass = () => {
    send({ type: 'pass', data: {} })
    selectedCards.value = []
  }

  // 获取提示
  const getHint = () => {
    send({ type: 'get_hint', data: {} })
  }

  // 选择/取消选择牌
  const toggleCard = (card: Card) => {
    const index = selectedCards.value.findIndex(
      (c) => c.Suit === card.Suit && c.Rank === card.Rank
    )
    if (index >= 0) {
      selectedCards.value.splice(index, 1)
    } else {
      selectedCards.value.push(card)
    }
  }

  // 处理消息
  const handleMessage = (msg: GameMessage) => {
    console.log('收到消息:', msg.type, msg.data)

    switch (msg.type) {
      case 'connected':
        playerId.value = msg.data.id
        console.log('玩家ID:', playerId.value)
        break

      case 'room_created':
        roomId.value = msg.data.room_id
        console.log('房间创建:', roomId.value)
        // 自动开始游戏
        setTimeout(() => {
          startGame()
        }, 300)
        break

      case 'room_joined':
        roomId.value = msg.data.room_id
        break

      case 'room_info':
        roomInfo.value = msg.data as RoomInfo
        console.log('房间信息更新:', roomInfo.value)
        break

      case 'deal_cards':
        const dealData = msg.data as DealCardsData
        myCards.value = dealData.cards
        isLandlord.value = dealData.is_landlord
        gameStarted.value = true
        console.log('收到手牌:', myCards.value.length, '张')
        break

      case 'bid_result':
        const bidData = msg.data as BidResultData
        console.log('叫分结果:', bidData)
        if (bidData.is_landlord) {
          // 地主确定
          if (roomInfo.value) {
            roomInfo.value.players.forEach(p => {
              if (p.id === bidData.player_id) {
                p.status = 3 // Landlord
              } else {
                p.status = 4 // Farmer
              }
            })
          }
          bidScore.value = bidData.bid
        }
        break

      case 'turn':
        const turnData = msg.data as TurnData
        isMyTurn.value = turnData.player_id === playerId.value
        lastPlayedCards.value = turnData.last_cards || []
        lastPlayer.value = turnData.last_player
        canPass.value = turnData.can_pass
        // 更新牌型结果
        if (turnData.last_cards && turnData.last_cards.length > 0) {
          lastHandResult.value = analyzeHand(turnData.last_cards)
        } else {
          lastHandResult.value = null
        }
        console.log('轮到:', turnData.player_name, '是我?', isMyTurn.value)
        break

      case 'play_result':
        const playData = msg.data as PlayResultData
        console.log('出牌结果:', playData)
        if (playData.player_id === playerId.value) {
          // 移除已出的牌
          if (!playData.is_pass && playData.cards) {
            myCards.value = myCards.value.filter(
              (card) =>
                !playData.cards.some(
                  (c) => c.Suit === card.Suit && c.Rank === card.Rank
                )
            )
          }
        }
        // 更新最后出的牌
        if (!playData.is_pass) {
          lastPlayedCards.value = playData.cards || []
          lastPlayer.value = playData.player_id
          // 更新牌型结果
          if (playData.cards && playData.cards.length > 0) {
            lastHandResult.value = analyzeHand(playData.cards)
          }
        } else {
          // 不出时清空牌型结果
          lastHandResult.value = null
        }
        break

      case 'game_over':
        gameOverData.value = msg.data as GameOverData
        console.log('游戏结束:', gameOverData.value)
        break

      case 'hint':
        selectedCards.value = msg.data as Card[]
        break

      case 'error':
        error.value = msg.data as string
        console.error('错误:', error.value)
        setTimeout(() => {
          error.value = ''
        }, 3000)
        break
    }
  }

  // 重置游戏
  const resetGame = () => {
    myCards.value = []
    selectedCards.value = []
    isLandlord.value = false
    lastPlayedCards.value = []
    lastPlayer.value = ''
    canPass.value = false
    isMyTurn.value = false
    gameOverData.value = null
    error.value = ''
    bidScore.value = 0
    gameStarted.value = false
    roomInfo.value = null
    lastHandResult.value = null
  }

  return {
    // 状态
    playerId,
    roomId,
    roomInfo,
    myCards,
    selectedCards,
    isLandlord,
    lastPlayedCards,
    lastPlayer,
    canPass,
    isMyTurn,
    gameOverData,
    error,
    bidScore,
    gameStarted,
    lastHandResult,

    // 计算属性
    isWaiting,
    isBidding,
    isPlaying,
    isGameOver,
    players,
    currentTurn,

    // 方法
    setSendFn,
    createRoom,
    startGame,
    bid,
    playCards,
    pass,
    getHint,
    toggleCard,
    handleMessage,
    resetGame,
  }
})
