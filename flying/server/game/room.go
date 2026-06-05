package game

import (
	"fmt"
	"sync"
)

// GameState 游戏状态
type GameState int

const (
	StateWaiting  GameState = 0 // 等待玩家
	StateBidding  GameState = 1 // 叫地主
	StatePlaying  GameState = 2 // 出牌中
	StateGameOver GameState = 3 // 游戏结束
)

// Room 游戏房间
type Room struct {
	mu          sync.RWMutex
	ID          string
	Players     []Player
	State       GameState
	Deck        *Deck
	BottomCards Cards      // 底牌
	CurrentTurn int        // 当前轮到谁
	LastCards   Cards      // 上一手出的牌
	LastResult  HandResult // 上一手牌型
	LastPlayer  int        // 上一个出牌的玩家
	BidScore    int        // 当前叫分
	BidPlayer   int        // 当前叫分玩家
	PassCount   int        // 连续不出次数
	BombCount   int        // 炸弹数量

	// WebSocket 回调
	OnMessage func(playerID string, msg GameMessage)
}

// NewRoom 创建新房间
func NewRoom(id string) *Room {
	return &Room{
		ID:         id,
		Players:    make([]Player, 0, 3),
		State:      StateWaiting,
		LastPlayer: -1,
	}
}

// AddPlayer 添加玩家
func (r *Room) AddPlayer(p Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.Players) >= 3 {
		return fmt.Errorf("房间已满")
	}

	r.Players = append(r.Players, p)
	return nil
}

// RemovePlayer 移除玩家
func (r *Room) RemovePlayer(playerID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, p := range r.Players {
		if p.GetID() == playerID {
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			break
		}
	}
}

// GetPlayer 获取玩家
func (r *Room) GetPlayer(playerID string) Player {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.Players {
		if p.GetID() == playerID {
			return p
		}
	}
	return nil
}

// IsFull 房间是否已满
func (r *Room) IsFull() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.Players) >= 3
}

// StartGame 开始游戏
func (r *Room) StartGame() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 自动添加AI玩家
	for len(r.Players) < 3 {
		aiID := fmt.Sprintf("ai_%d", len(r.Players))
		ai := NewAIPlayer(aiID, fmt.Sprintf("电脑%d", len(r.Players)))
		r.Players = append(r.Players, ai)
	}

	// 创建牌组并发牌
	r.Deck = NewDeck()
	p1, p2, p3, bottom := r.Deck.Deal()
	r.BottomCards = bottom

	// 分配手牌
	r.Players[0].SetCards(p1)
	r.Players[1].SetCards(p2)
	r.Players[2].SetCards(p3)

	// 设置玩家状态
	for _, p := range r.Players {
		p.SetStatus(PlayerStatusFarmer)
	}

	// 开始叫地主
	r.State = StateBidding
	r.CurrentTurn = 0
	r.BidScore = 0
	r.BidPlayer = -1
	r.PassCount = 0
	r.BombCount = 0
	r.LastCards = nil
	r.LastResult = HandResult{}
	r.LastPlayer = -1

	// 通知房间状态更新
	r.broadcast(GameMessage{
		Type: MsgTypeRoomInfo,
		Data: r.getRoomInfo(),
	})

	// 通知发牌
	for _, p := range r.Players {
		if !p.IsAI() && r.OnMessage != nil {
			r.OnMessage(p.GetID(), GameMessage{
				Type: MsgTypeDealCards,
				Data: DealCardsData{
					Cards:      p.GetCards(),
					IsLandlord: false,
				},
			})
		}
	}

	// 通知轮到谁叫地主
	r.notifyTurn()

	return nil
}

// HandleBid 处理叫地主
func (r *Room) HandleBid(playerID string, score int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.State != StateBidding {
		return fmt.Errorf("当前不是叫地主阶段")
	}

	playerIdx := r.getPlayerIndex(playerID)
	if playerIdx != r.CurrentTurn {
		return fmt.Errorf("还没轮到你叫地主")
	}

	// 叫分逻辑
	if score > r.BidScore {
		r.BidScore = score
		r.BidPlayer = playerIdx
	}

	// 通知叫分结果
	r.broadcast(GameMessage{
		Type: MsgTypeBidResult,
		Data: BidResultData{
			PlayerID:   playerID,
			PlayerName: r.Players[playerIdx].GetName(),
			Bid:        score,
			IsLandlord: false,
		},
	})

	// 叫了3分直接成为地主
	if score == 3 {
		r.setLandlord(playerIdx)
		return nil
	}

	// 连续3人叫过或都叫过
	if score == 0 {
		r.PassCount++
	}

	if r.PassCount >= 3 {
		if r.BidPlayer >= 0 {
			r.setLandlord(r.BidPlayer)
		} else {
			// 都不叫，重新发牌
			r.restartGame()
		}
		return nil
	}

	// 下一个人叫
	r.CurrentTurn = (r.CurrentTurn + 1) % 3

	// 通知房间状态更新
	r.broadcast(GameMessage{
		Type: MsgTypeRoomInfo,
		Data: r.getRoomInfo(),
	})

	// 通知轮到谁
	r.notifyTurn()

	return nil
}

// setLandlord 设置地主
func (r *Room) setLandlord(playerIdx int) {
	r.State = StatePlaying

	// 设置地主
	r.Players[playerIdx].SetStatus(PlayerStatusLandlord)
	r.Players[playerIdx].SetCards(append(r.Players[playerIdx].GetCards(), r.BottomCards...))

	// 排序地主的牌
	r.Players[playerIdx].SetCards(SortCards(r.Players[playerIdx].GetCards()))

	// 通知所有人地主确定
	r.broadcast(GameMessage{
		Type: MsgTypeBidResult,
		Data: BidResultData{
			PlayerID:   r.Players[playerIdx].GetID(),
			PlayerName: r.Players[playerIdx].GetName(),
			Bid:        r.BidScore,
			IsLandlord: true,
		},
	})

	// 通知房间状态更新
	r.broadcast(GameMessage{
		Type: MsgTypeRoomInfo,
		Data: r.getRoomInfo(),
	})

	// 通知地主获得底牌
	if !r.Players[playerIdx].IsAI() && r.OnMessage != nil {
		r.OnMessage(r.Players[playerIdx].GetID(), GameMessage{
			Type: MsgTypeDealCards,
			Data: DealCardsData{
				Cards:      r.Players[playerIdx].GetCards(),
				IsLandlord: true,
			},
		})
	}

	// 地主先出牌
	r.CurrentTurn = playerIdx
	r.LastCards = nil
	r.LastResult = HandResult{}
	r.LastPlayer = -1
	r.PassCount = 0

	r.notifyTurn()
}

// HandlePlay 处理出牌
func (r *Room) HandlePlay(playerID string, cards Cards) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.State != StatePlaying {
		return fmt.Errorf("当前不是出牌阶段")
	}

	playerIdx := r.getPlayerIndex(playerID)
	if playerIdx != r.CurrentTurn {
		return fmt.Errorf("还没轮到你出牌")
	}

	// 不出
	if len(cards) == 0 {
		if r.LastPlayer == playerIdx || r.LastPlayer == -1 {
			return fmt.Errorf("必须出牌")
		}

		// 通知不出
		r.broadcast(GameMessage{
			Type: MsgTypePlayResult,
			Data: PlayResultData{
				PlayerID:   playerID,
				PlayerName: r.Players[playerIdx].GetName(),
				Cards:      nil,
				IsPass:     true,
			},
		})

		r.PassCount++
		r.CurrentTurn = (r.CurrentTurn + 1) % 3

		// 如果其他两人都不出，轮到上次出牌的人
		if r.PassCount >= 2 {
			r.CurrentTurn = r.LastPlayer
			r.LastCards = nil
			r.LastResult = HandResult{}
			r.PassCount = 0
		}

		// 通知房间状态更新
		r.broadcast(GameMessage{
			Type: MsgTypeRoomInfo,
			Data: r.getRoomInfo(),
		})

		r.notifyTurn()
		return nil
	}

	// 验证出牌
	result := AnalyzeHand(cards)
	if result.Type == CardTypeNone {
		return fmt.Errorf("无效的牌型")
	}

	// 需要压牌时验证
	if r.LastCards != nil && r.LastPlayer != playerIdx {
		if !CanBeat(cards, r.LastResult) {
			return fmt.Errorf("出的牌不够大")
		}
	}

	// 验证玩家手牌是否包含这些牌
	playerCards := r.Players[playerIdx].GetCards()
	for _, card := range cards {
		if !playerCards.Contains(card) {
			return fmt.Errorf("你没有这张牌: %s", card.String())
		}
	}

	// 移除手牌
	r.Players[playerIdx].SetCards(playerCards.Remove(cards))

	// 记录炸弹
	if result.Type == CardTypeBomb || result.Type == CardTypeRocket {
		r.BombCount++
	}

	// 更新状态
	r.LastCards = cards
	r.LastResult = result
	r.LastPlayer = playerIdx
	r.PassCount = 0

	// 通知出牌结果
	r.broadcast(GameMessage{
		Type: MsgTypePlayResult,
		Data: PlayResultData{
			PlayerID:   playerID,
			PlayerName: r.Players[playerIdx].GetName(),
			Cards:      cards,
			IsPass:     false,
			CardType:   result.Type,
		},
	})

	// 检查是否获胜
	if len(r.Players[playerIdx].GetCards()) == 0 {
		r.gameOver(playerIdx)
		return nil
	}

	// 下一个人出牌
	r.CurrentTurn = (r.CurrentTurn + 1) % 3

	// 通知房间状态更新
	r.broadcast(GameMessage{
		Type: MsgTypeRoomInfo,
		Data: r.getRoomInfo(),
	})

	r.notifyTurn()

	return nil
}

// gameOver 游戏结束
func (r *Room) gameOver(winnerIdx int) {
	r.State = StateGameOver

	winner := r.Players[winnerIdx]
	isLandlord := winner.GetStatus() == PlayerStatusLandlord

	// 计算分数
	baseScore := r.BidScore
	if r.BombCount > 0 {
		baseScore *= (1 << r.BombCount)
	}

	// 通知房间状态更新
	r.broadcast(GameMessage{
		Type: MsgTypeRoomInfo,
		Data: r.getRoomInfo(),
	})

	// 通知所有人
	r.broadcast(GameMessage{
		Type: MsgTypeGameOver,
		Data: GameOverData{
			Winner:     winner.GetID(),
			WinnerName: winner.GetName(),
			IsLandlord: isLandlord,
			Score:      baseScore,
		},
	})
}

// restartGame 重新开始游戏
func (r *Room) restartGame() {
	r.State = StateWaiting
	r.Deck = nil
	r.BottomCards = nil
	r.LastCards = nil
	r.LastResult = HandResult{}
	r.LastPlayer = -1
	r.BidScore = 0
	r.BidPlayer = -1
	r.PassCount = 0
	r.BombCount = 0

	// 重新开始
	r.StartGame()
}

// notifyTurn 通知轮到谁
func (r *Room) notifyTurn() {
	currentPlayer := r.Players[r.CurrentTurn]

	msg := GameMessage{
		Type: MsgTypeTurn,
		Data: TurnData{
			PlayerID:   currentPlayer.GetID(),
			PlayerName: currentPlayer.GetName(),
			LastCards:  r.LastCards,
			LastPlayer: func() string {
				if r.LastPlayer >= 0 {
					return r.Players[r.LastPlayer].GetID()
				}
				return ""
			}(),
			CanPass: r.LastPlayer != r.CurrentTurn && r.LastPlayer != -1,
		},
	}

	// 通知所有人
	r.broadcast(msg)

	// 如果是 AI，自动处理
	if currentPlayer.IsAI() {
		go r.handleAITurn()
	}
}

// handleAITurn 处理 AI 回合
func (r *Room) handleAITurn() {
	// 等待一段时间模拟思考
	// time.Sleep(1 * time.Second)

	ai := r.Players[r.CurrentTurn].(*AIPlayer)

	if r.State == StateBidding {
		// AI 叫地主
		bid := ai.Bid(r.BidScore)
		r.HandleBid(ai.GetID(), bid)
	} else if r.State == StatePlaying {
		// AI 出牌
		cards := ai.Play(r.LastCards, r.LastResult)
		r.HandlePlay(ai.GetID(), cards)
	}
}

// broadcast 广播消息
func (r *Room) broadcast(msg GameMessage) {
	if r.OnMessage == nil {
		return
	}

	for _, p := range r.Players {
		if !p.IsAI() {
			r.OnMessage(p.GetID(), msg)
		}
	}
}

// getPlayerIndex 获取玩家索引
func (r *Room) getPlayerIndex(playerID string) int {
	for i, p := range r.Players {
		if p.GetID() == playerID {
			return i
		}
	}
	return -1
}

// getRoomInfo 获取房间信息（内部调用，不获取锁）
func (r *Room) getRoomInfo() map[string]interface{} {
	players := make([]map[string]interface{}, 0)
	for _, p := range r.Players {
		players = append(players, map[string]interface{}{
			"id":         p.GetID(),
			"name":       p.GetName(),
			"card_count": len(p.GetCards()),
			"status":     p.GetStatus(),
			"is_ai":      p.IsAI(),
		})
	}

	return map[string]interface{}{
		"id":        r.ID,
		"state":     r.State,
		"players":   players,
		"bid_score": r.BidScore,
		"turn":      r.CurrentTurn,
	}
}

// GetRoomInfo 获取房间信息（外部调用，获取锁）
func (r *Room) GetRoomInfo() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.getRoomInfo()
}
