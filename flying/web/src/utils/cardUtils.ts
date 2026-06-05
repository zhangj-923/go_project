import type { Card } from '../types/game'
import { Suit } from '../types/game'

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

// 牌型结果
export interface HandResult {
  type: CardType
  mainRank: number
  length: number
}

// 分析出牌的牌型
export function analyzeHand(cards: Card[]): HandResult {
  if (cards.length === 0) {
    return { type: CardType.None, mainRank: 0, length: 0 }
  }

  // 按点数统计
  const rankCount: Record<number, number> = {}
  cards.forEach(card => {
    rankCount[card.Rank] = (rankCount[card.Rank] || 0) + 1
  })

  const ranks = Object.keys(rankCount).map(Number).sort((a, b) => a - b)

  // 火箭：大小王
  if (cards.length === 2 && cards[0].Suit === Suit.Joker && cards[1].Suit === Suit.Joker) {
    return { type: CardType.Rocket, mainRank: 17, length: 1 }
  }

  // 炸弹：4张相同
  if (cards.length === 4 && ranks.length === 1) {
    return { type: CardType.Bomb, mainRank: ranks[0], length: 1 }
  }

  // 单张
  if (cards.length === 1) {
    return { type: CardType.Single, mainRank: cards[0].Rank, length: 1 }
  }

  // 对子
  if (cards.length === 2 && ranks.length === 1) {
    return { type: CardType.Pair, mainRank: ranks[0], length: 1 }
  }

  // 三条
  if (cards.length === 3 && ranks.length === 1) {
    return { type: CardType.Triple, mainRank: ranks[0], length: 1 }
  }

  // 三带一
  if (cards.length === 4 && ranks.length === 2) {
    for (const rank of ranks) {
      if (rankCount[rank] === 3) {
        return { type: CardType.TripleOne, mainRank: rank, length: 1 }
      }
    }
  }

  // 三带二
  if (cards.length === 5 && ranks.length === 2) {
    for (const rank of ranks) {
      if (rankCount[rank] === 3) {
        return { type: CardType.TripleTwo, mainRank: rank, length: 1 }
      }
    }
  }

  // 顺子：5张以上连续单牌（不含2和王）
  if (cards.length >= 5 && ranks.length === cards.length) {
    if (isStraight(ranks) && ranks[ranks.length - 1] < 15) {
      return { type: CardType.Straight, mainRank: ranks[ranks.length - 1], length: cards.length }
    }
  }

  // 连对：3对以上连续对子（不含2和王）
  if (cards.length >= 6 && cards.length % 2 === 0) {
    const pairRanks: number[] = []
    let allPairs = true
    for (const rank of ranks) {
      if (rankCount[rank] !== 2) {
        allPairs = false
        break
      }
      pairRanks.push(rank)
    }
    if (allPairs && pairRanks.length === cards.length / 2) {
      if (isStraight(pairRanks) && pairRanks[pairRanks.length - 1] < 15) {
        return { type: CardType.PairStraight, mainRank: pairRanks[pairRanks.length - 1], length: pairRanks.length }
      }
    }
  }

  // 飞机不带：2组以上连续三条
  const tripleRanks: number[] = []
  for (const rank of ranks) {
    if (rankCount[rank] === 3) {
      tripleRanks.push(rank)
    }
  }
  if (tripleRanks.length >= 2 && cards.length === tripleRanks.length * 3) {
    if (isStraight(tripleRanks) && tripleRanks[tripleRanks.length - 1] < 15) {
      return { type: CardType.Plane, mainRank: tripleRanks[tripleRanks.length - 1], length: tripleRanks.length }
    }
  }

  // 飞机带单
  if (tripleRanks.length >= 2) {
    if (isStraight(tripleRanks) && tripleRanks[tripleRanks.length - 1] < 15) {
      const singleCount = cards.length - tripleRanks.length * 3
      if (singleCount === tripleRanks.length) {
        return { type: CardType.PlaneOne, mainRank: tripleRanks[tripleRanks.length - 1], length: tripleRanks.length }
      }
    }
  }

  // 飞机带双
  if (tripleRanks.length >= 2) {
    if (isStraight(tripleRanks) && tripleRanks[tripleRanks.length - 1] < 15) {
      let pairCount = 0
      for (const count of Object.values(rankCount)) {
        if (count === 2) pairCount++
      }
      if (pairCount === tripleRanks.length && cards.length === tripleRanks.length * 5) {
        return { type: CardType.PlaneTwo, mainRank: tripleRanks[tripleRanks.length - 1], length: tripleRanks.length }
      }
    }
  }

  // 四带二
  if (cards.length === 6) {
    for (const rank of ranks) {
      if (rankCount[rank] === 4) {
        return { type: CardType.FourTwo, mainRank: rank, length: 1 }
      }
    }
  }

  // 四带两对
  if (cards.length === 8) {
    for (const rank of ranks) {
      if (rankCount[rank] === 4) {
        let pairCount = 0
        for (const count of Object.values(rankCount)) {
          if (count === 2) pairCount++
        }
        if (pairCount === 2) {
          return { type: CardType.FourFour, mainRank: rank, length: 1 }
        }
      }
    }
  }

  return { type: CardType.None, mainRank: 0, length: 0 }
}

// 判断是否为连续序列
function isStraight(ranks: number[]): boolean {
  if (ranks.length < 2) return true

  const sorted = [...ranks].sort((a, b) => a - b)
  for (let i = 1; i < sorted.length; i++) {
    if (sorted[i] - sorted[i - 1] !== 1) {
      return false
    }
  }
  return true
}

// 判断当前出牌是否能压过上一手
export function canBeat(current: Card[], lastResult: HandResult): boolean {
  const currentResult = analyzeHand(current)

  if (currentResult.type === CardType.None) {
    return false
  }

  // 火箭最大
  if (currentResult.type === CardType.Rocket) {
    return true
  }

  // 炸弹可以压其他牌型
  if (currentResult.type === CardType.Bomb) {
    if (lastResult.type === CardType.Rocket) {
      return false
    }
    if (lastResult.type === CardType.Bomb) {
      return currentResult.mainRank > lastResult.mainRank
    }
    return true
  }

  // 同类型比较
  if (currentResult.type === lastResult.type) {
    if (currentResult.type === CardType.Straight || currentResult.type === CardType.PairStraight) {
      return currentResult.length === lastResult.length && currentResult.mainRank > lastResult.mainRank
    }
    return currentResult.mainRank > lastResult.mainRank
  }

  return false
}

// 获取牌型名称
export function getCardTypeName(type: CardType): string {
  const names: Record<CardType, string> = {
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
  return names[type] || '未知'
}
