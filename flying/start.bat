@echo off
chcp 65001 >nul

echo ========================================
echo       斗地主游戏启动脚本
echo ========================================
echo.

echo 检查依赖...

where go >nul 2>nul
if %errorlevel% neq 0 (
    echo 错误: 未安装 Go
    pause
    exit /b 1
)

where node >nul 2>nul
if %errorlevel% neq 0 (
    echo 错误: 未安装 Node.js
    pause
    exit /b 1
)

where npm >nul 2>nul
if %errorlevel% neq 0 (
    echo 错误: 未安装 npm
    pause
    exit /b 1
)

echo 依赖检查通过
echo.

echo 安装前端依赖...
cd web
call npm install
cd ..
echo 前端依赖安装完成
echo.

echo 启动后端服务器...
cd server
start "后端服务器" cmd /k "go run main.go"
cd ..
echo 后端服务器已启动
echo.

echo 启动前端开发服务器...
cd web
start "前端开发服务器" cmd /k "npm run dev"
cd ..
echo 前端开发服务器已启动
echo.

echo ========================================
echo       游戏已启动
echo ========================================
echo.
echo 后端地址: http://localhost:8080
echo 前端地址: http://localhost:3000
echo.
echo 按任意键停止所有服务...
pause >nul

echo 正在停止服务...
taskkill /FI "WindowTitle eq 后端服务器*" /F >nul 2>nul
taskkill /FI "WindowTitle eq 前端开发服务器*" /F >nul 2>nul
echo 所有服务已停止
pause
