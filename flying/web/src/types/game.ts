// 花色定义
export enum Suit {
  Spade = 0,   // 黑桃
  Heart = 1,   // 红桃
  Club = 2,    // 梅花
  Diamond = 3, // 方块
  Joker = 4,   // 王
}

// 牌面值定义
export enum Rank {
  R3 = 3,
  R4 = 4,
  R5 = 5,
  R6 = 6,
  R7 = 7,
  R8 = 8,
  R9 = 9,
  R10 = 10,
  J = 11,
  Q = 12,
  K = 13,
  A = 14,
  R2 = 15,
  BJ = 16, // 小王
  RJ = 17, // 大王
}

// 牌型定义
export enum CardType {
  None = 0,
  Single = 1,
  Pair = 2,
  Triple = 3,
  TripleOne = 4,
  TripleTwo = 5,
  Straight = 6,
  PairStraight = 7,
  Plane = 8,
  PlaneOne = 9,
  PlaneTwo = 10,
  FourTwo = 11,
  FourFour = 12,
  Bomb = 13,
  Rocket = 14,
}

// 扑克牌
export interface Card {
  Suit: number
  Rank: number
}

// 玩家状态
export enum PlayerStatus {
  Waiting = 0,
  Playing = 1,
  Ready = 2,
  Landlord = 3,
  Farmer = 4,
}

// 玩家信息
export interface Player {
  id: string
  name: string
  card_count: number
  status: PlayerStatus
  is_ai: boolean
}

// 游戏状态
export enum GameState {
  Waiting = 0,
  Bidding = 1,
  Playing = 2,
  GameOver = 3,
}

// 房间信息
export interface RoomInfo {
  id: string
  state: GameState
  players: Player[]
  bid_score: number
  turn: number
}

// WebSocket 消息
export interface GameMessage {
  type: string
  data: any
  player_id?: string
}

// 发牌数据
export interface DealCardsData {
  cards: Card[]
  is_landlord: boolean
}

// 叫地主数据
export interface BidData {
  current_bid: number
  min_bid: number
}

// 叫地主结果
export interface BidResultData {
  player_id: string
  player_name: string
  bid: number
  is_landlord: boolean
}

// 出牌数据
export interface PlayData {
  last_cards: Card[]
  last_player: string
  can_pass: boolean
}

// 出牌结果
export interface PlayResultData {
  player_id: string
  player_name: string
  cards: Card[]
  is_pass: boolean
  card_type: CardType
}

// 游戏结束数据
export interface GameOverData {
  winner: string
  winner_name: string
  is_landlord: boolean
  score: number
}

// 轮次数据
export interface TurnData {
  player_id: string
  player_name: string
  last_cards: Card[]
  last_player: string
  can_pass: boolean
}

// 花色符号映射
export const suitSymbols: Record<Suit, string> = {
  [Suit.Spade]: '♠',
  [Suit.Heart]: '♥',
  [Suit.Club]: '♣',
  [Suit.Diamond]: '♦',
  [Suit.Joker]: '🃏',
}

// 牌面值符号映射
export const rankSymbols: Record<Rank, string> = {
  [Rank.R3]: '3',
  [Rank.R4]: '4',
  [Rank.R5]: '5',
  [Rank.R6]: '6',
  [Rank.R7]: '7',
  [Rank.R8]: '8',
  [Rank.R9]: '9',
  [Rank.R10]: '10',
  [Rank.J]: 'J',
  [Rank.Q]: 'Q',
  [Rank.K]: 'K',
  [Rank.A]: 'A',
  [Rank.R2]: '2',
  [Rank.BJ]: '小王',
  [Rank.RJ]: '大王',
}

// 获取牌的显示文本
export function getCardDisplay(card: Card): string {
  if (card.Suit === Suit.Joker) {
    return rankSymbols[card.Rank as Rank]
  }
  return `${suitSymbols[card.Suit as Suit]}${rankSymbols[card.Rank as Rank]}`
}

// 获取牌的颜色
export function getCardColor(card: Card): string {
  if (card.Suit === Suit.Joker) {
    return card.Rank === Rank.RJ ? '#ff0000' : '#000000'
  }
  return card.Suit === Suit.Heart || card.Suit === Suit.Diamond ? '#ff0000' : '#000000'
}

// 牌型名称映射
export const cardTypeNames: Record<CardType, string> = {
  [CardType.None]: '无效',
  [CardType.Single]: '单张',
  [CardType.Pair]: '对子',
  [CardType.Triple]: '三条',
  [CardType.TripleOne]: '三带一',
  [CardType.TripleTwo]: '三带二',
  [CardType.Straight]: '顺子',
  [CardType.PairStraight]: '连对',
  [CardType.Plane]: '飞机',
  [CardType.PlaneOne]: '飞机带单',
  [CardType.PlaneTwo]: '飞机带双',
  [CardType.FourTwo]: '四带二',
  [CardType.FourFour]: '四带两对',
  [CardType.Bomb]: '炸弹',
  [CardType.Rocket]: '火箭',
}
