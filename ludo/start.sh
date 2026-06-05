#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}=======================================${NC}"
echo -e "${BLUE}       启动飞行棋 (Ludo) 游戏服务      ${NC}"
echo -e "${BLUE}=======================================${NC}"

# 获取项目根目录
ROOT_DIR=$(dirname "$0")
cd "$ROOT_DIR"
ROOT_DIR=$(pwd)

# 确保端口未被占用 (可选)
# lsof -ti:8080 | xargs kill -9 2>/dev/null
# lsof -ti:5173 | xargs kill -9 2>/dev/null

echo -e "\n${YELLOW}▶ 正在启动 Go 后端服务 (端口 18189)...${NC}"
cd "$ROOT_DIR/server"
go run main.go > /dev/null 2>&1 &
SERVER_PID=$!
echo "✓ 后端服务启动成功 [PID: $SERVER_PID]"

echo -e "\n${YELLOW}▶ 正在启动 Vue 前端服务 (端口 15173)...${NC}"
cd "$ROOT_DIR/web"
npm run dev -- --host --port 15173 > /dev/null 2>&1 &
WEB_PID=$!
echo "✓ 前端服务启动成功 [PID: $WEB_PID]"

echo -e "\n${GREEN}=======================================${NC}"
echo -e "${GREEN}🎉 游戏已成功运行！${NC}"
echo -e "${GREEN}=======================================${NC}"
echo -e "\n请在浏览器中访问以下地址开始游戏："
echo -e "👉 本地访问: ${YELLOW}http://localhost:15173${NC}"

# 获取本机局域网 IP
LAN_IP=$(ipconfig getifaddr en0 2>/dev/null || ipconfig getifaddr en1 2>/dev/null)
if [ ! -z "$LAN_IP" ]; then
    echo -e "👉 局域网访问 (供其他人加入): ${YELLOW}http://${LAN_IP}:15173${NC}"
else
    echo -e "👉 局域网访问 (供其他人加入): 请查看 npm run dev 的网络输出地址"
fi
echo -e "\n后端 WebSocket 监听地址: ws://localhost:18189/ws"
echo -e "\n${YELLOW}提示: 按下 Ctrl+C 即可同时停止前后端服务${NC}\n"

# 捕获 Ctrl+C 信号，优雅关闭后台进程
cleanup() {
    echo -e "\n${YELLOW}正在停止所有服务...${NC}"
    kill $SERVER_PID 2>/dev/null
    kill $WEB_PID 2>/dev/null
    echo -e "✓ 服务已完全退出"
    exit 0
}
trap cleanup SIGINT SIGTERM

# 保持脚本运行，等待结束信号
wait
