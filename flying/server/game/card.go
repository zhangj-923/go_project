package game

import "fmt"

// 花色定义
type Suit int

const (
	SuitSpade   Suit = 0 // 黑桃
	SuitHeart   Suit = 1 // 红桃
	SuitClub    Suit = 2 // 梅花
	SuitDiamond Suit = 3 // 方块
	SuitJoker   Suit = 4 // 王
)

// 牌面值定义
type Rank int

const (
	Rank3  Rank = 3
	Rank4  Rank = 4
	Rank5  Rank = 5
	Rank6  Rank = 6
	Rank7  Rank = 7
	Rank8  Rank = 8
	Rank9  Rank = 9
	Rank10 Rank = 10
	RankJ  Rank = 11
	RankQ  Rank = 12
	RankK  Rank = 13
	RankA  Rank = 14
	Rank2  Rank = 15
	RankBJ Rank = 16 // 小王
	RankRJ Rank = 17 // 大王
)

// 牌型定义
type CardType int

const (
	CardTypeNone         CardType = 0  // 无效
	CardTypeSingle       CardType = 1  // 单张
	CardTypePair         CardType = 2  // 对子
	CardTypeTriple       CardType = 3  // 三条
	CardTypeTripleOne    CardType = 4  // 三带一
	CardTypeTripleTwo    CardType = 5  // 三带二
	CardTypeStraight     CardType = 6  // 顺子
	CardTypePairStraight CardType = 7  // 连对
	CardTypePlane        CardType = 8  // 飞机不带
	CardTypePlaneOne     CardType = 9  // 飞机带单
	CardTypePlaneTwo     CardType = 10 // 飞机带双
	CardTypeFourTwo      CardType = 11 // 四带二
	CardTypeFourFour     CardType = 12 // 四带两对
	CardTypeBomb         CardType = 13 // 炸弹
	CardTypeRocket       CardType = 14 // 火箭
)

// Card 扑克牌结构
type Card struct {
	Suit Suit `json:"Suit"`
	Rank Rank `json:"Rank"`
}

// NewCard 创建一张牌
func NewCard(suit Suit, rank Rank) Card {
	return Card{Suit: suit, Rank: rank}
}

// String 返回牌的字符串表示
func (c Card) String() string {
	suits := map[Suit]string{
		SuitSpade:   "♠",
		SuitHeart:   "♥",
		SuitClub:    "♣",
		SuitDiamond: "♦",
		SuitJoker:   "🃏",
	}

	ranks := map[Rank]string{
		Rank3:  "3",
		Rank4:  "4",
		Rank5:  "5",
		Rank6:  "6",
		Rank7:  "7",
		Rank8:  "8",
		Rank9:  "9",
		Rank10: "10",
		RankJ:  "J",
		RankQ:  "Q",
		RankK:  "K",
		RankA:  "A",
		Rank2:  "2",
		RankBJ: "小王",
		RankRJ: "大王",
	}

	if c.Suit == SuitJoker {
		return ranks[c.Rank]
	}
	return fmt.Sprintf("%s%s", suits[c.Suit], ranks[c.Rank])
}

// Equal 判断两张牌是否相同
func (c Card) Equal(other Card) bool {
	return c.Suit == other.Suit && c.Rank == other.Rank
}

// Cards 牌组类型
type Cards []Card

// Len 实现 sort.Interface
func (c Cards) Len() int {
	return len(c)
}

// Less 实现 sort.Interface
func (c Cards) Less(i, j int) bool {
	if c[i].Rank != c[j].Rank {
		return c[i].Rank < c[j].Rank
	}
	return c[i].Suit < c[j].Suit
}

// Swap 实现 sort.Interface
func (c Cards) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Contains 判断牌组是否包含指定牌
func (c Cards) Contains(card Card) bool {
	for _, cc := range c {
		if cc.Equal(card) {
			return true
		}
	}
	return false
}

// Remove 从牌组中移除指定牌
func (c Cards) Remove(cards Cards) Cards {
	result := make(Cards, 0, len(c))
	removed := make(map[int]bool)

	for _, card := range c {
		found := false
		for j, toRemove := range cards {
			if !removed[j] && card.Equal(toRemove) {
				removed[j] = true
				found = true
				break
			}
		}
		if !found {
			result = append(result, card)
		}
	}
	return result
}

// CountRank 统计某个点数的牌数
func (c Cards) CountRank(rank Rank) int {
	count := 0
	for _, card := range c {
		if card.Rank == rank {
			count++
		}
	}
	return count
}

// GetRankGroups 按点数分组
func (c Cards) GetRankGroups() map[Rank]Cards {
	groups := make(map[Rank]Cards)
	for _, card := range c {
		groups[card.Rank] = append(groups[card.Rank], card)
	}
	return groups
}
