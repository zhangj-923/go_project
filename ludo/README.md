# 飞行棋 (Ludo) 游戏项目

这是一个基于 Go + Vue 3 开发的现代版飞行棋网页游戏，支持单人（人机大战）和局域网多人的对战模式。

## 目录结构
- `/server`: Go 后端服务，基于 Gorilla WebSocket 提供实时通信能力。
- `/web`: Vue 3 + Vite 前端工程，包含游戏渲染、界面交互及状态管理。
- `/docs`: 项目文档。
  - [`design.md`](docs/design.md): 项目设计方案、规则逻辑、接口协议等。
  - [`acceptance.md`](docs/acceptance.md): 游戏验收总结文档。
- `start.sh`: 游戏一键启动脚本。

## 快速开始

运行游戏根目录下的启动脚本：

```bash
chmod +x start.sh
./start.sh
```

- 本地访问地址: `http://localhost:5173`
- 脚本同时还会打印可供局域网内其他设备访问的网络 IP 地址。

## 技术栈
- **后端**: Go 1.21+, net/http, gorilla/websocket
- **前端**: Vue 3 (Composition API), Vite, Pinia, TypeScript, 原生 SVG + CSS 动画

## 特性概览
- **三种模式**: 单机游戏(1v3 AI)、本地同屏、局域网联机。
- **动态特效**: 毛玻璃 UI、3D 骰子物理翻滚、棋子跳跃/飞跃动画。
- **智能化 AI**: 后端内置机器人的回合状态机，能够在玩家缺席时或者单人游戏中自主接管操作。
