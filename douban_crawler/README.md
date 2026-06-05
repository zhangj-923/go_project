# 豆瓣电影爬虫练手项目

这是一个使用 Go 语言和 Colly 框架编写的入门级爬虫项目。
它抓取了豆瓣电影 Top 250 的数据，并最终将结果导出到 Excel 文件。

## 如何运行

1. 下载项目依赖：
```sh
go mod tidy
```

2. 运行爬虫：
```sh
go run main.go
```

3. 运行完成后，你会在当前目录下找到名为 `douban_top250.xlsx` 的 Excel 文件。
