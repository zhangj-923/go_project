package models

// Movie 电影数据模型，负责定义我们需要抓取的结构化字段
type Movie struct {
	Rank   string // 排名
	Title  string // 电影名称
	Rating string // 豆瓣评分
	Quote  string // 经典短评
}
