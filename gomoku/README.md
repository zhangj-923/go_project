# Gomoku (五子棋) - Modern AI Board Game

## 项目介绍
基于 Go 语言和 Ebiten 2D 游戏引擎开发的一款现代化五子棋桌面小游戏。
具备流畅的交互体验、美观的国风木质 UI 设计，并集成了从新手到专家的多难度 AI。

## 特性
- **人机对战**: 黑白交替，与智能 AI 展开对弈。
- **三种难度模式**: 
  - Easy: 简单随机落子，适合新手上手体验。
  - Normal: 启发式评分引擎，具备基础的攻防意识。
  - Hard: Minimax + Alpha-Beta 剪枝搜索，提供强大的压制力。
- **现代化 UI**: 木质棋盘、高质量阴影、平滑的界面体验。
- **60 FPS**: 优化渲染管线，保证极低延迟和丝滑动画。

## 项目结构
```text
gomoku/
├── cmd/
│   └── game/        # 游戏入口
├── internal/
│   ├── ai/          # AI 逻辑 (Minimax, Heuristic)
│   ├── board/       # 棋盘数据与胜负判定
│   ├── engine/      # 游戏状态机循环
│   └── ui/          # Ebiten UI 渲染
├── go.mod           # 模块配置
└── README.md
```

## 启动方式
确保您已安装 Go 1.22+，在项目根目录下执行：
```bash
go run ./cmd/game
```
