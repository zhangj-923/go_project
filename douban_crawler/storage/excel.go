package storage

import (
	"fmt"
	"strconv"

	"douban_crawler/models"

	"github.com/xuri/excelize/v2"
)

// ExcelStorage Excel 存储实现
type ExcelStorage struct {
	FileName string
}

// Save 实现 Storage 接口的 Save 方法
func (e *ExcelStorage) Save(movies []models.Movie) error {
	// 创建一个新的 Excel 文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("关闭 Excel 文件失败:", err)
		}
	}()

	sheetName := "Top250"
	// 重命名默认的工作表
	f.SetSheetName("Sheet1", sheetName)

	// 设置表头
	headers := []string{"排名", "电影名称", "评分", "经典短评"}
	for col, header := range headers {
		// CoordinatesToCellName 将列号和行号转换为单元格坐标，例如 (1, 1) -> A1
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据
	for row, movie := range movies {
		// 行号从 2 开始，因为第 1 行是表头
		f.SetCellValue(sheetName, "A"+strconv.Itoa(row+2), movie.Rank)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row+2), movie.Title)
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row+2), movie.Rating)
		f.SetCellValue(sheetName, "D"+strconv.Itoa(row+2), movie.Quote)
	}

	// 保存文件到指定路径
	if err := f.SaveAs(e.FileName); err != nil {
		return err
	}
	return nil
}
