package scraper

import (
	"fmt"
	"strings"
	"time"

	"douban_crawler/models"
	"douban_crawler/storage"

	"github.com/gocolly/colly/v2"
)

// DoubanScraper 豆瓣爬虫核心类
type DoubanScraper struct {
	Collector *colly.Collector
	Storage   storage.Storage
	Movies    []models.Movie
}

// NewDoubanScraper 初始化豆瓣爬虫
func NewDoubanScraper(store storage.Storage) *DoubanScraper {
	// 初始化 Colly 框架
	c := colly.NewCollector(
		colly.AllowedDomains("movie.douban.com"),
		colly.Async(false), // 使用同步抓取，避免并发过高被封
	)

	// 配置限制规则 (礼貌爬取)
	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*douban.com*",
		RandomDelay: 2 * time.Second, // 每次请求随机延迟 0~2 秒
	})

	// 伪装 HTTP 请求头，极其关键
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
		// 添加一个虚假的 bid cookie，有些时候可以绕过豆瓣的安全拦截
		r.Headers.Set("Cookie", `bid="VjS-xXW0_aE";`)
		fmt.Println("正在请求:", r.URL.String())
	})

	// 增加错误捕获，看看是不是被豆瓣拦截了或者返回了 403
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("请求失败: URL=%s, 错误=%v, 状态码=%d\n", r.Request.URL, err, r.StatusCode)
	})

	scraper := &DoubanScraper{
		Collector: c,
		Storage:   store,
		Movies:    make([]models.Movie, 0),
	}

	scraper.setupCallbacks()
	return scraper
}

// setupCallbacks 设置解析回调函数
func (s *DoubanScraper) setupCallbacks() {
	// 解析每个电影卡片 (HTML DOM 解析)
	s.Collector.OnHTML(".item", func(e *colly.HTMLElement) {
		movie := models.Movie{
			Rank:   e.ChildText(".pic em"),
			Title:  e.ChildText(".info .hd .title"),
			Rating: e.ChildText(".info .bd .rating_num"),
			Quote:  e.ChildText(".info .bd .quote .inq"),
		}

		// 简单的数据清洗：去掉标题多余的别名部分
		titles := strings.Split(movie.Title, "\u00a0")
		if len(titles) > 0 {
			movie.Title = titles[0]
		}

		s.Movies = append(s.Movies, movie)
	})

	// 自动翻页逻辑：寻找“后页”的链接并自动访问
	s.Collector.OnHTML(".paginator .next a", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")
		if nextPage != "" {
			// 构建下一页的绝对 URL 并访问
			_ = e.Request.Visit(e.Request.AbsoluteURL(nextPage))
		}
	})
}

// Start 启动爬虫
func (s *DoubanScraper) Start() {
	// 启动种子 URL
	_ = s.Collector.Visit("https://movie.douban.com/top250?start=0")

	// 等待所有的请求处理完毕
	s.Collector.Wait()

	// 抓取结束后的回调 (同步模式下，Visit 执行完毕代表抓取结束)
	fmt.Printf("抓取完成，共获取 %d 条电影数据，正在保存...\n", len(s.Movies))
	err := s.Storage.Save(s.Movies)
	if err != nil {
		fmt.Println("保存数据失败:", err)
	} else {
		fmt.Println("数据成功保存到文件！")
	}
}
