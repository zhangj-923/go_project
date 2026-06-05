package config

import "time"

type Config struct {
	BoardWidth   int
	BoardHeight  int
	CellWidth    int
	CellHeight   int
	InitialSpeed time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		BoardWidth:   10,
		BoardHeight:  20,
		CellWidth:    2,
		CellHeight:   1, // 终端字符通常是瘦高的，宽2高1可以近似为正方形
		InitialSpeed: 800 * time.Millisecond,
	}
}
