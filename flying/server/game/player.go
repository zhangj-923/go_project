package game

// PlayerStatus 玩家状态
type PlayerStatus int

const (
	PlayerStatusWaiting  PlayerStatus = 0 // 等待中
	PlayerStatusPlaying  PlayerStatus = 1 // 游戏中
	PlayerStatusReady    PlayerStatus = 2 // 准备
	PlayerStatusLandlord PlayerStatus = 3 // 地主
	PlayerStatusFarmer   PlayerStatus = 4 // 农民
)

// Player 玩家接口
type Player interface {
	GetID() string
	GetName() string
	GetCards() Cards
	SetCards(cards Cards)
	GetStatus() PlayerStatus
	SetStatus(status PlayerStatus)
	IsAI() bool

	// 游戏相关
	Bid(currentBid int) int                            // 叫地主，返回叫的分数（0表示不叫）
	Play(lastCards Cards, lastResult HandResult) Cards // 出牌
}

// HumanPlayer 人类玩家
type HumanPlayer struct {
	ID     string
	Name   string
	Cards  Cards
	Status PlayerStatus
}

func NewHumanPlayer(id, name string) *HumanPlayer {
	return &HumanPlayer{
		ID:     id,
		Name:   name,
		Status: PlayerStatusWaiting,
	}
}

func (p *HumanPlayer) GetID() string                 { return p.ID }
func (p *HumanPlayer) GetName() string               { return p.Name }
func (p *HumanPlayer) GetCards() Cards               { return p.Cards }
func (p *HumanPlayer) SetCards(cards Cards)          { p.Cards = cards }
func (p *HumanPlayer) GetStatus() PlayerStatus       { return p.Status }
func (p *HumanPlayer) SetStatus(status PlayerStatus) { p.Status = status }
func (p *HumanPlayer) IsAI() bool                    { return false }

func (p *HumanPlayer) Bid(currentBid int) int {
	// 人类玩家通过WebSocket交互，这里返回-1表示等待输入
	return -1
}

func (p *HumanPlayer) Play(lastCards Cards, lastResult HandResult) Cards {
	// 人类玩家通过WebSocket交互，这里返回nil表示等待输入
	return nil
}

// GameMessage 游戏消息
type GameMessage struct {
	Type     string      `json:"type"`
	Data     interface{} `json:"data"`
	PlayerID string      `json:"player_id,omitempty"`
}

// 消息类型常量
const (
	MsgTypeStartGame   = "start_game"   // 开始游戏
	MsgTypeDealCards   = "deal_cards"   // 发牌
	MsgTypeBid         = "bid"          // 叫地主
	MsgTypeBidResult   = "bid_result"   // 叫地主结果
	MsgTypePlay        = "play"         // 出牌
	MsgTypePlayResult  = "play_result"  // 出牌结果
	MsgTypePass        = "pass"         // 不出
	MsgTypeGameOver    = "game_over"    // 游戏结束
	MsgTypeTurn        = "turn"         // 轮到谁
	MsgTypeError       = "error"        // 错误
	MsgTypeRoomInfo    = "room_info"    // 房间信息
	MsgTypePlayerJoin  = "player_join"  // 玩家加入
	MsgTypePlayerLeave = "player_leave" // 玩家离开
)

// DealCardsData 发牌数据
type DealCardsData struct {
	Cards      Cards `json:"cards"`
	IsLandlord bool  `json:"is_landlord"`
}

// BidData 叫地主数据
type BidData struct {
	CurrentBid int `json:"current_bid"`
	MinBid     int `json:"min_bid"`
}

// BidResultData 叫地主结果
type BidResultData struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	Bid        int    `json:"bid"`
	IsLandlord bool   `json:"is_landlord"`
}

// PlayData 出牌数据
type PlayData struct {
	LastCards  Cards  `json:"last_cards"`
	LastPlayer string `json:"last_player"`
	CanPass    bool   `json:"can_pass"`
}

// PlayResultData 出牌结果
type PlayResultData struct {
	PlayerID   string   `json:"player_id"`
	PlayerName string   `json:"player_name"`
	Cards      Cards    `json:"cards"`
	IsPass     bool     `json:"is_pass"`
	CardType   CardType `json:"card_type"`
}

// GameOverData 游戏结束数据
type GameOverData struct {
	Winner     string `json:"winner"`
	WinnerName string `json:"winner_name"`
	IsLandlord bool   `json:"is_landlord"`
	Score      int    `json:"score"`
}

// TurnData 轮次数据
type TurnData struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	LastCards  Cards  `json:"last_cards"`
	LastPlayer string `json:"last_player"`
	CanPass    bool   `json:"can_pass"`
}
