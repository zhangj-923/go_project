package game

import "sort"

// HandResult 出牌结果
type HandResult struct {
	Type     CardType // 牌型
	MainRank Rank     // 主牌点数
	Length   int      // 顺子/连对长度
}

// AnalyzeHand 分析出牌的牌型
func AnalyzeHand(cards Cards) HandResult {
	if len(cards) == 0 {
		return HandResult{Type: CardTypeNone}
	}

	// 按点数统计
	rankCount := make(map[Rank]int)
	for _, card := range cards {
		rankCount[card.Rank]++
	}

	// 获取所有点数并排序
	ranks := make([]int, 0)
	for rank := range rankCount {
		ranks = append(ranks, int(rank))
	}
	sort.Ints(ranks)

	// 火箭：大小王
	if len(cards) == 2 && cards[0].Rank == RankBJ && cards[1].Rank == RankRJ {
		return HandResult{Type: CardTypeRocket}
	}

	// 炸弹：4张相同
	if len(cards) == 4 && len(rankCount) == 1 {
		return HandResult{Type: CardTypeBomb, MainRank: Rank(ranks[0])}
	}

	// 单张
	if len(cards) == 1 {
		return HandResult{Type: CardTypeSingle, MainRank: cards[0].Rank}
	}

	// 对子
	if len(cards) == 2 && len(rankCount) == 1 {
		return HandResult{Type: CardTypePair, MainRank: Rank(ranks[0])}
	}

	// 三条
	if len(cards) == 3 && len(rankCount) == 1 {
		return HandResult{Type: CardTypeTriple, MainRank: Rank(ranks[0])}
	}

	// 三带一
	if len(cards) == 4 && len(rankCount) == 2 {
		for rank, count := range rankCount {
			if count == 3 {
				return HandResult{Type: CardTypeTripleOne, MainRank: rank}
			}
		}
	}

	// 三带二
	if len(cards) == 5 && len(rankCount) == 2 {
		for rank, count := range rankCount {
			if count == 3 {
				return HandResult{Type: CardTypeTripleTwo, MainRank: rank}
			}
		}
	}

	// 顺子：5张以上连续单牌（不含2和王）
	if len(cards) >= 5 && len(rankCount) == len(cards) {
		if isStraight(ranks) && Rank(ranks[len(ranks)-1]) < Rank2 {
			return HandResult{Type: CardTypeStraight, MainRank: Rank(ranks[len(ranks)-1]), Length: len(cards)}
		}
	}

	// 连对：3对以上连续对子（不含2和王）
	if len(cards) >= 6 && len(cards)%2 == 0 {
		allPairs := true
		pairRanks := make([]int, 0)
		for rank, count := range rankCount {
			if count != 2 {
				allPairs = false
				break
			}
			pairRanks = append(pairRanks, int(rank))
		}
		if allPairs && len(pairRanks) == len(cards)/2 {
			sort.Ints(pairRanks)
			if isStraight(pairRanks) && Rank(pairRanks[len(pairRanks)-1]) < Rank2 {
				return HandResult{Type: CardTypePairStraight, MainRank: Rank(pairRanks[len(pairRanks)-1]), Length: len(pairRanks)}
			}
		}
	}

	// 飞机不带：2组以上连续三条
	tripleRanks := make([]int, 0)
	for rank, count := range rankCount {
		if count == 3 {
			tripleRanks = append(tripleRanks, int(rank))
		}
	}
	if len(tripleRanks) >= 2 && len(cards) == len(tripleRanks)*3 {
		sort.Ints(tripleRanks)
		if isStraight(tripleRanks) && Rank(tripleRanks[len(tripleRanks)-1]) < Rank2 {
			return HandResult{Type: CardTypePlane, MainRank: Rank(tripleRanks[len(tripleRanks)-1]), Length: len(tripleRanks)}
		}
	}

	// 飞机带单：连续三条 + 等量单牌
	if len(tripleRanks) >= 2 {
		sort.Ints(tripleRanks)
		if isStraight(tripleRanks) && Rank(tripleRanks[len(tripleRanks)-1]) < Rank2 {
			singleCount := len(cards) - len(tripleRanks)*3
			if singleCount == len(tripleRanks) {
				return HandResult{Type: CardTypePlaneOne, MainRank: Rank(tripleRanks[len(tripleRanks)-1]), Length: len(tripleRanks)}
			}
		}
	}

	// 飞机带双：连续三条 + 等量对子
	if len(tripleRanks) >= 2 {
		sort.Ints(tripleRanks)
		if isStraight(tripleRanks) && Rank(tripleRanks[len(tripleRanks)-1]) < Rank2 {
			pairCount := 0
			for _, count := range rankCount {
				if count == 2 {
					pairCount++
				}
			}
			if pairCount == len(tripleRanks) && len(cards) == len(tripleRanks)*5 {
				return HandResult{Type: CardTypePlaneTwo, MainRank: Rank(tripleRanks[len(tripleRanks)-1]), Length: len(tripleRanks)}
			}
		}
	}

	// 四带二：四条 + 2张单牌
	if len(cards) == 6 {
		for rank, count := range rankCount {
			if count == 4 {
				return HandResult{Type: CardTypeFourTwo, MainRank: rank}
			}
		}
	}

	// 四带两对：四条 + 2对
	if len(cards) == 8 {
		for rank, count := range rankCount {
			if count == 4 {
				pairCount := 0
				for _, c := range rankCount {
					if c == 2 {
						pairCount++
					}
				}
				if pairCount == 2 {
					return HandResult{Type: CardTypeFourFour, MainRank: rank}
				}
			}
		}
	}

	return HandResult{Type: CardTypeNone}
}

// isStraight 判断是否为连续序列
func isStraight(ranks []int) bool {
	if len(ranks) < 2 {
		return true
	}

	sort.Ints(ranks)
	for i := 1; i < len(ranks); i++ {
		if ranks[i]-ranks[i-1] != 1 {
			return false
		}
	}
	return true
}

// CanBeat 判断当前出牌是否能压过上一手
func CanBeat(current Cards, last HandResult) bool {
	currentResult := AnalyzeHand(current)

	if currentResult.Type == CardTypeNone {
		return false
	}

	// 火箭最大
	if currentResult.Type == CardTypeRocket {
		return true
	}

	// 炸弹可以压其他牌型
	if currentResult.Type == CardTypeBomb {
		if last.Type == CardTypeRocket {
			return false
		}
		if last.Type == CardTypeBomb {
			return currentResult.MainRank > last.MainRank
		}
		return true
	}

	// 同类型比较
	if currentResult.Type == last.Type {
		if currentResult.Type == CardTypeStraight || currentResult.Type == CardTypePairStraight {
			return currentResult.Length == last.Length && currentResult.MainRank > last.MainRank
		}
		return currentResult.MainRank > last.MainRank
	}

	return false
}

// FindAllHands 找出手牌中所有可以出的牌（用于AI提示）
func FindAllHands(hand Cards, lastResult HandResult) []Cards {
	results := make([]Cards, 0)

	if lastResult.Type == CardTypeNone {
		// 自由出牌，找出所有可能的牌型
		results = append(results, findAllPossibleHands(hand)...)
	} else {
		// 需要压牌，找出所有能压过的牌
		results = append(results, findBeatingHands(hand, lastResult)...)
	}

	return results
}

// findAllPossibleHands 找出所有可能的出牌组合
func findAllPossibleHands(hand Cards) []Cards {
	results := make([]Cards, 0)
	rankGroups := hand.GetRankGroups()

	// 单张
	for _, cards := range rankGroups {
		results = append(results, Cards{cards[0]})
	}

	// 对子
	for _, cards := range rankGroups {
		if len(cards) >= 2 {
			results = append(results, cards[:2])
		}
	}

	// 三条
	for _, cards := range rankGroups {
		if len(cards) >= 3 {
			results = append(results, cards[:3])
		}
	}

	// 三带一
	for rank, cards := range rankGroups {
		if len(cards) >= 3 {
			for otherRank, otherCards := range rankGroups {
				if otherRank != rank {
					triple := cards[:3]
					result := append(Cards{}, triple...)
					result = append(result, otherCards[0])
					results = append(results, result)
					break
				}
			}
		}
	}

	// 三带二
	for rank, cards := range rankGroups {
		if len(cards) >= 3 {
			for otherRank, otherCards := range rankGroups {
				if otherRank != rank && len(otherCards) >= 2 {
					triple := cards[:3]
					result := append(Cards{}, triple...)
					result = append(result, otherCards[:2]...)
					results = append(results, result)
					break
				}
			}
		}
	}

	// 炸弹
	for _, cards := range rankGroups {
		if len(cards) == 4 {
			results = append(results, cards)
		}
	}

	// 火箭
	hasBJ := false
	hasRJ := false
	for _, card := range hand {
		if card.Rank == RankBJ {
			hasBJ = true
		}
		if card.Rank == RankRJ {
			hasRJ = true
		}
	}
	if hasBJ && hasRJ {
		results = append(results, Cards{NewCard(SuitJoker, RankBJ), NewCard(SuitJoker, RankRJ)})
	}

	return results
}

// findBeatingHands 找出能压过指定牌型的所有出牌
func findBeatingHands(hand Cards, lastResult HandResult) []Cards {
	results := make([]Cards, 0)
	rankGroups := hand.GetRankGroups()

	switch lastResult.Type {
	case CardTypeSingle:
		// 找更大的单张
		for _, cards := range rankGroups {
			if cards[0].Rank > lastResult.MainRank {
				results = append(results, Cards{cards[0]})
			}
		}
		// 炸弹和火箭
		results = append(results, findBombsAndRockets(hand)...)

	case CardTypePair:
		// 找更大的对子
		for _, cards := range rankGroups {
			if len(cards) >= 2 && cards[0].Rank > lastResult.MainRank {
				results = append(results, cards[:2])
			}
		}
		results = append(results, findBombsAndRockets(hand)...)

	case CardTypeTriple:
		// 找更大的三条
		for _, cards := range rankGroups {
			if len(cards) >= 3 && cards[0].Rank > lastResult.MainRank {
				results = append(results, cards[:3])
			}
		}
		results = append(results, findBombsAndRockets(hand)...)

	case CardTypeTripleOne:
		// 找更大的三带一
		for rank, cards := range rankGroups {
			if len(cards) >= 3 && rank > lastResult.MainRank {
				for otherRank, otherCards := range rankGroups {
					if otherRank != rank {
						triple := cards[:3]
						result := append(Cards{}, triple...)
						result = append(result, otherCards[0])
						results = append(results, result)
						break
					}
				}
			}
		}
		results = append(results, findBombsAndRockets(hand)...)

	case CardTypeTripleTwo:
		// 找更大的三带二
		for rank, cards := range rankGroups {
			if len(cards) >= 3 && rank > lastResult.MainRank {
				for otherRank, otherCards := range rankGroups {
					if otherRank != rank && len(otherCards) >= 2 {
						triple := cards[:3]
						result := append(Cards{}, triple...)
						result = append(result, otherCards[:2]...)
						results = append(results, result)
						break
					}
				}
			}
		}
		results = append(results, findBombsAndRockets(hand)...)

	case CardTypeBomb:
		// 找更大的炸弹
		for _, cards := range rankGroups {
			if len(cards) == 4 && cards[0].Rank > lastResult.MainRank {
				results = append(results, cards)
			}
		}
		// 火箭
		results = append(results, findRocket(hand)...)
	}

	return results
}

// findBombsAndRockets 找出手牌中的炸弹和火箭
func findBombsAndRockets(hand Cards) []Cards {
	results := make([]Cards, 0)
	rankGroups := hand.GetRankGroups()

	for _, cards := range rankGroups {
		if len(cards) == 4 {
			results = append(results, cards)
		}
	}

	results = append(results, findRocket(hand)...)
	return results
}

// findRocket 找出手牌中的火箭
func findRocket(hand Cards) []Cards {
	hasBJ := false
	hasRJ := false
	for _, card := range hand {
		if card.Rank == RankBJ {
			hasBJ = true
		}
		if card.Rank == RankRJ {
			hasRJ = true
		}
	}

	if hasBJ && hasRJ {
		return []Cards{{NewCard(SuitJoker, RankBJ), NewCard(SuitJoker, RankRJ)}}
	}
	return nil
}
