package game

import (
	"math/rand"
	"time"
)

// Deck 牌组
type Deck struct {
	cards Cards
}

// NewDeck 创建一副新牌（54张）
func NewDeck() *Deck {
	d := &Deck{}
	d.init()
	return d
}

// init 初始化一副牌
func (d *Deck) init() {
	d.cards = make(Cards, 0, 54)

	// 4种花色各13张
	suits := []Suit{SuitSpade, SuitHeart, SuitClub, SuitDiamond}
	for _, suit := range suits {
		for rank := Rank3; rank <= Rank2; rank++ {
			d.cards = append(d.cards, NewCard(suit, rank))
		}
	}

	// 大小王
	d.cards = append(d.cards, NewCard(SuitJoker, RankBJ))
	d.cards = append(d.cards, NewCard(SuitJoker, RankRJ))
}

// Shuffle 洗牌
func (d *Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// Deal 发牌：返回3个玩家的牌和底牌
func (d *Deck) Deal() (player1, player2, player3, landlord Cards) {
	d.Shuffle()

	// 每人17张
	player1 = d.cards[0:17]
	player2 = d.cards[17:34]
	player3 = d.cards[34:51]

	// 底牌3张
	landlord = d.cards[51:54]

	// 排序
	player1 = SortCards(player1)
	player2 = SortCards(player2)
	player3 = SortCards(player3)

	return
}

// SortCards 对牌组进行排序
func SortCards(cards Cards) Cards {
	sorted := make(Cards, len(cards))
	copy(sorted, cards)

	// 使用简单的冒泡排序
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Rank > sorted[i].Rank {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted
}
