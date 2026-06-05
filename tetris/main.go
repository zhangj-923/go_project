package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"tetris/internal/config"
	"tetris/internal/game"
	"tetris/internal/ui"

	"github.com/gdamore/tcell/v2"
)

func main() {
	// 自定义 Usage 提示，当用户输入 -h 或参数错误时，打印友好的说明
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "=== 俄罗斯方块 (Tetris) ===\n")
		fmt.Fprintf(os.Stderr, "用法: ./tetris [选项]\n\n")
		fmt.Fprintf(os.Stderr, "例如: \n")
		fmt.Fprintf(os.Stderr, "  ./tetris                (使用默认大小启动)\n")
		fmt.Fprintf(os.Stderr, "  ./tetris -w 20 -h 30    (开启 20列30行 的超大棋盘)\n")
		fmt.Fprintf(os.Stderr, "  ./tetris -cw 4 -ch 2    (将方块渲染得更加巨大)\n\n")
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
	}

	// 定义命令行参数
	boardW := flag.Int("w", 10, "棋盘的宽度 (格数)")
	boardH := flag.Int("h", 20, "棋盘的高度 (格数)")
	cellW := flag.Int("cw", 2, "单个方块渲染的字符宽度")
	cellH := flag.Int("ch", 1, "单个方块渲染的字符高度")
	helpFlag := flag.Bool("help", false, "显示帮助信息")

	flag.Parse()

	// 拦截 --help 参数，主动打印提示并退出
	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	// 终端尺寸合法性校验，避免因参数过小导致渲染崩溃
	if *boardW < 4 || *boardH < 4 || *cellW < 1 || *cellH < 1 {
		fmt.Fprintf(os.Stderr, "错误: 尺寸参数过小。棋盘宽高至少为 4，方块宽高至少为 1\n")
		os.Exit(1)
	}

	// 初始化配置
	cfg := config.DefaultConfig()
	cfg.BoardWidth = *boardW
	cfg.BoardHeight = *boardH
	cfg.CellWidth = *cellW
	cfg.CellHeight = *cellH

	// 初始化 Tcell 屏幕
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating screen: %v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing screen: %v\n", err)
		os.Exit(1)
	}
	screen.SetStyle(tcell.StyleDefault.Background(tcell.NewRGBColor(15, 15, 25)))
	screen.EnableMouse()
	screen.Clear()
	defer screen.Fini()

	// 组装系统
	engine := game.NewEngine(cfg)
	renderer := ui.NewRenderer(screen, cfg)

	// 事件通道
	eventCh := make(chan tcell.Event)
	go func() {
		for {
			eventCh <- screen.PollEvent()
		}
	}()

	ticker := time.NewTicker(engine.DropInterval)
	defer ticker.Stop()

	// 游戏主循环
	for {
		engine.Mu.Lock()
		renderer.Draw(engine)
		engine.Mu.Unlock()

		select {
		case ev := <-eventCh:
			switch ev := ev.(type) {
			case *tcell.EventResize:
				engine.Mu.Lock()
				renderer.UpdateLayout()
				screen.Sync()
				engine.Mu.Unlock()

			case *tcell.EventKey:
				engine.Mu.Lock()
				switch {
				case ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' || ev.Rune() == 'Q':
					engine.Mu.Unlock()
					return

				case ev.Rune() == 'r' || ev.Rune() == 'R':
					if engine.GameOver {
						engine.Reset()
					}

				case ev.Rune() == 'p' || ev.Rune() == 'P':
					engine.TogglePause()

				case !engine.GameOver && !engine.Paused:
					switch {
					case ev.Key() == tcell.KeyLeft:
						engine.MoveLeft()
					case ev.Key() == tcell.KeyRight:
						engine.MoveRight()
					case ev.Key() == tcell.KeyDown:
						if engine.MoveDown() {
							engine.Score++
						}
					case ev.Key() == tcell.KeyUp:
						engine.Rotate()
					case ev.Rune() == ' ':
						engine.HardDrop()
					}
				}
				engine.Mu.Unlock()
			}

		case <-ticker.C:
			engine.Mu.Lock()
			engine.Tick()
			// 刷新下落速度
			ticker.Reset(engine.DropInterval)
			engine.Mu.Unlock()
		}
	}
}
