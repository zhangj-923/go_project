package main

import (
	"douban_crawler/scraper"
	"douban_crawler/storage"
	"fmt"
)

func main() {
	fmt.Println("=== 开始豆瓣电影 Top 250 抓取任务 ===")

	// 1. 依赖注入：初始化数据存储模块 (这里改为使用 ExcelStorage)
	store := &storage.ExcelStorage{
		FileName: "douban_top250.xlsx", // 后缀改为 .xlsx
	}

	// 2. 初始化爬虫核心模块，并将存储模块传入
	doubanScraper := scraper.NewDoubanScraper(store)

	// 3. 启动爬虫任务
	doubanScraper.Start()

	fmt.Println("=== 抓取任务圆满结束 ===")
}
