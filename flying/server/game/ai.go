package game

import (
	"math/rand"
	"sort"
	"time"
)

// AIPlayer AI玩家
type AIPlayer struct {
	ID     string
	Name   string
	Cards  Cards
	Status PlayerStatus
	rng    *rand.Rand
}

// NewAIPlayer 创建AI玩家
func NewAIPlayer(id, name string) *AIPlayer {
	return &AIPlayer{
		ID:     id,
		Name:   name,
		Status: PlayerStatusWaiting,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (p *AIPlayer) GetID() string                 { return p.ID }
func (p *AIPlayer) GetName() string               { return p.Name }
func (p *AIPlayer) GetCards() Cards               { return p.Cards }
func (p *AIPlayer) SetCards(cards Cards)          { p.Cards = cards }
func (p *AIPlayer) GetStatus() PlayerStatus       { return p.Status }
func (p *AIPlayer) SetStatus(status PlayerStatus) { p.Status = status }
func (p *AIPlayer) IsAI() bool                    { return true }

// Bid 叫地主策略
func (p *AIPlayer) Bid(currentBid int) int {
	strength := p.evaluateHand()

	// 根据手牌强度决定叫分
	switch {
	case strength >= 10:
		return 3 // 强牌，叫3分
	case strength >= 7:
		if currentBid < 2 {
			return 2
		}
		return 0
	case strength >= 4:
		if currentBid < 1 {
			return 1
		}
		return 0
	default:
		// 弱牌也有一定概率叫1分
		if currentBid == 0 && p.rng.Float64() < 0.3 {
			return 1
		}
		return 0
	}
}

// evaluateHand 评估手牌强度
func (p *AIPlayer) evaluateHand() int {
	score := 0
	rankGroups := p.Cards.GetRankGroups()

	// 大小王各加4分
	for _, card := range p.Cards {
		if card.Rank == RankRJ {
			score += 4
		}
		if card.Rank == RankBJ {
			score += 3
		}
	}

	// 2的数量
	twoCount := rankGroups[Rank2]
	score += len(twoCount) * 2

	// A的数量
	aceCount := rankGroups[RankA]
	score += len(aceCount)

	// 炸弹
	for _, cards := range rankGroups {
		if len(cards) == 4 {
			score += 6
		}
	}

	// 三条
	for _, cards := range rankGroups {
		if len(cards) == 3 {
			score += 2
		}
	}

	return score
}

// Play 出牌策略
func (p *AIPlayer) Play(lastCards Cards, lastResult HandResult) Cards {
	if lastCards == nil || len(lastCards) == 0 {
		// 自由出牌
		return p.playFree()
	}

	// 需要压牌
	return p.playBeat(lastResult)
}

// playFree 自由出牌策略
func (p *AIPlayer) playFree() Cards {
	rankGroups := p.Cards.GetRankGroups()

	// 策略1：先出顺子
	straight := p.findStraight(rankGroups)
	if straight != nil {
		return straight
	}

	// 策略2：出连对
	pairStraight := p.findPairStraight(rankGroups)
	if pairStraight != nil {
		return pairStraight
	}

	// 策略3：出三条（带单或带双）
	triple := p.findTriple(rankGroups)
	if triple != nil {
		return triple
	}

	// 策略4：出对子
	pair := p.findPair(rankGroups)
	if pair != nil {
		return pair
	}

	// 策略5：出单张（从小到大）
	single := p.findSingle(rankGroups)
	if single != nil {
		return single
	}

	// 如果都没有，出最小的牌
	if len(p.Cards) > 0 {
		return Cards{p.Cards[len(p.Cards)-1]}
	}

	return nil
}

// findStraight 找顺子
func (p *AIPlayer) findStraight(rankGroups map[Rank]Cards) Cards {
	// 找5张以上的连续单牌
	straight := make(Cards, 0)

	for rank := Rank3; rank <= RankA; rank++ {
		if cards, exists := rankGroups[rank]; exists && len(cards) >= 1 {
			straight = append(straight, cards[0])
			if len(straight) >= 5 {
				return straight
			}
		} else {
			straight = make(Cards, 0)
		}
	}

	return nil
}

// findPairStraight 找连对
func (p *AIPlayer) findPairStraight(rankGroups map[Rank]Cards) Cards {
	// 找3对以上的连续对子
	pairs := make(Cards, 0)

	for rank := Rank3; rank <= RankA; rank++ {
		if cards, exists := rankGroups[rank]; exists && len(cards) >= 2 {
			pairs = append(pairs, cards[:2]...)
			if len(pairs) >= 6 {
				return pairs
			}
		} else {
			pairs = make(Cards, 0)
		}
	}

	return nil
}

// findTriple 找三条
func (p *AIPlayer) findTriple(rankGroups map[Rank]Cards) Cards {
	for rank := Rank3; rank <= Rank2; rank++ {
		if cards, exists := rankGroups[rank]; exists && len(cards) >= 3 {
			// 优先带单牌
			for otherRank, otherCards := range rankGroups {
				if otherRank != rank && len(otherCards) >= 1 {
					result := Cards{cards[0], cards[1], cards[2], otherCards[0]}
					return result
				}
			}
			// 没有单牌可带，只出三条
			return cards[:3]
		}
	}
	return nil
}

// findPair 找对子
func (p *AIPlayer) findPair(rankGroups map[Rank]Cards) Cards {
	// 从小到大找对子
	for rank := Rank3; rank <= Rank2; rank++ {
		if cards, exists := rankGroups[rank]; exists && len(cards) >= 2 {
			return cards[:2]
		}
	}
	return nil
}

// findSingle 找单张
func (p *AIPlayer) findSingle(rankGroups map[Rank]Cards) Cards {
	// 从小到大找单张
	for rank := Rank3; rank <= Rank2; rank++ {
		if cards, exists := rankGroups[rank]; exists && len(cards) >= 1 {
			return Cards{cards[0]}
		}
	}
	return nil
}

// playBeat 压牌策略
func (p *AIPlayer) playBeat(lastResult HandResult) Cards {
	possibleHands := FindAllHands(p.Cards, lastResult)

	if len(possibleHands) == 0 {
		// 没有能压的牌
		return nil
	}

	// 策略：选择最小的能压的牌
	// 如果是队友出的牌，考虑不出
	if p.isTeammatePlay(lastResult) {
		// 50%概率不出
		if p.rng.Float64() < 0.5 {
			return nil
		}
	}

	// 找最小的能压的牌
	var bestHand Cards
	bestScore := 999

	for _, hand := range possibleHands {
		score := p.evaluateHandScore(hand)
		if score < bestScore {
			bestScore = score
			bestHand = hand
		}
	}

	return bestHand
}

// evaluateHandScore 评估出牌的分数（越小越好）
func (p *AIPlayer) evaluateHandScore(hand Cards) int {
	score := 0

	for _, card := range hand {
		// 小牌分数低
		if card.Rank <= Rank10 {
			score += 1
		} else if card.Rank <= RankA {
			score += 2
		} else {
			score += 3
		}
	}

	// 炸弹和火箭加分
	result := AnalyzeHand(hand)
	if result.Type == CardTypeBomb {
		score += 10
	}
	if result.Type == CardTypeRocket {
		score += 20
	}

	return score
}

// isTeammatePlay 判断是否是队友出的牌
func (p *AIPlayer) isTeammatePlay(lastResult HandResult) bool {
	// 简化处理：如果自己是农民，上一手也是农民出的，则是队友
	// 实际应该根据玩家ID判断，这里简化
	return false
}

// SuggestPlay 提示出牌（用于人类玩家的提示功能）
func SuggestPlay(hand Cards, lastCards Cards, lastResult HandResult) Cards {
	ai := &AIPlayer{
		Cards: hand,
		rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return ai.Play(lastCards, lastResult)
}

// SortCardsForDisplay 用于显示的排序（从小到大）
func SortCardsForDisplay(cards Cards) Cards {
	sorted := make(Cards, len(cards))
	copy(sorted, cards)

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Rank != sorted[j].Rank {
			return sorted[i].Rank < sorted[j].Rank
		}
		return sorted[i].Suit < sorted[j].Suit
	})

	return sorted
}
