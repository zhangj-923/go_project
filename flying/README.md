# 斗地主单机游戏

一个使用 Go + Vue 3 + TypeScript 开发的单机斗地主游戏，支持两个电脑 AI 玩家。

## 功能特性

- 完整的斗地主规则实现
- 两个 AI 电脑玩家
- WebSocket 实时通信
- 精美的游戏界面
- 叫地主、出牌、提示等功能

## 项目结构

```
flying/
├── server/                    # Go 后端
│   ├── main.go               # 入口
│   ├── game/                 # 游戏核心逻辑
│   │   ├── card.go           # 牌定义
│   │   ├── deck.go           # 发牌器
│   │   ├── player.go         # 玩家接口
│   │   ├── ai.go             # AI 策略
│   │   ├── rules.go          # 出牌规则
│   │   └── room.go           # 游戏房间
│   └── handler/              # WebSocket 处理
│       └── websocket.go
├── web/                      # Vue 前端
│   ├── src/
│   │   ├── components/       # 游戏组件
│   │   ├── composables/      # WebSocket 连接
│   │   ├── stores/           # 状态管理
│   │   └── types/            # 类型定义
│   └── package.json
└── README.md
```

## 快速开始

### 1. 启动后端

```bash
cd server
go run main.go
```

后端将在 http://localhost:8080 启动

### 2. 启动前端

```bash
cd web
npm install
npm run dev
```

前端将在 http://localhost:3000 启动

### 3. 开始游戏

1. 打开浏览器访问 http://localhost:3000
2. 点击"创建房间"
3. 点击"开始游戏"
4. 系统会自动加入两个 AI 玩家

## 游戏规则

### 基本规则

- 3人游戏，使用54张牌（含大小王）
- 每人17张牌，3张底牌
- 叫分最高者成为地主，获得底牌
- 地主 vs 两个农民

### 牌型

- 单张：任意一张牌
- 对子：两张相同点数的牌
- 三条：三张相同点数的牌
- 三带一：三条 + 一张单牌
- 三带二：三条 + 一对
- 顺子：5张或以上连续单牌（不含2和王）
- 连对：3对或以上连续对子
- 飞机：2组或以上连续三条
- 炸弹：四张相同点数的牌
- 火箭：大小王

### 出牌规则

- 首轮由地主先出牌
- 后续玩家必须出比上家更大的牌，或选择不出
- 炸弹可以压任何牌型（火箭除外）
- 火箭最大，可以压任何牌

## 技术栈

- **后端**: Go + Gorilla WebSocket
- **前端**: Vue 3 + TypeScript + Pinia
- **通信**: WebSocket

## 开发说明

### 添加新功能

1. 后端游戏逻辑在 `server/game/` 目录
2. 前端组件在 `web/src/components/` 目录
3. WebSocket 消息类型在 `server/game/player.go` 定义

### AI 策略优化

AI 策略在 `server/game/ai.go` 中实现，可以通过以下方式优化：

1. 改进 `evaluateHand` 函数的估值算法
2. 优化 `playBeat` 函数的压牌策略
3. 添加记牌器功能
4. 实现更智能的配合策略

## License

MIT
