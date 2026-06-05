package storage

import (
	"encoding/csv"
	"os"

	"douban_crawler/models"
)

// CSVStorage CSV 存储实现
type CSVStorage struct {
	FileName string
}

// Save 实现 Storage 接口的 Save 方法
func (c *CSVStorage) Save(movies []models.Movie) error {
	file, err := os.Create(c.FileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入 UTF-8 BOM 以防止 Excel 打开乱码 (可选)
	file.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	writer.Write([]string{"排名", "电影名称", "评分", "经典短评"})

	// 写入数据行
	for _, movie := range movies {
		writer.Write([]string{movie.Rank, movie.Title, movie.Rating, movie.Quote})
	}
	return nil
}
