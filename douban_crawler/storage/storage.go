package storage

import "douban_crawler/models"

// Storage 数据存储接口，符合面向接口编程原则。
// 日后如果想把数据存入 MySQL 或 MongoDB，只需新建一个实现该接口的类，无需修改爬虫核心逻辑。
type Storage interface {
	Save(movies []models.Movie) error
}
