package main

import (
	"github.com/gin-gonic/gin"
	pagination "github.com/wujiyu98/gin-pagination"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)


type Product struct {
	ID                uint64         `gorm:"primaryKey;column:id;type:bigint(20) unsigned;not null"`
	ProductCategoryID int            `gorm:"index:product_category_id;column:product_category_id;type:int(11);not null"`
	ManufacturerID    int            `gorm:"index:manufacturer_id;column:manufacturer_id;type:int(11);not null"`
	Title             string         `gorm:"unique;column:title;type:varchar(255);not null"`
	Pathname          string         `gorm:"index:pathname;column:pathname;type:varchar(255);default:''"`
	RangePrice        int8           `gorm:"column:range_price;type:tinyint(4);default:1"`
	ImageSrc          string         `gorm:"column:image_src;type:varchar(255);default:''"`
	PdfSrc            string         `gorm:"column:pdf_src;type:varchar(255);default:''"`
	Summary           string         `gorm:"column:summary;type:varchar(255);default:''"`
	Stock             int            `gorm:"column:stock;type:int(11);not null;default:0"`
	Price             float64        `gorm:"column:price;type:decimal(10,4);not null"`
	Content           string         `gorm:"column:content;type:longtext"`
	Hot               int8           `gorm:"index:hot_index;column:hot;type:tinyint(4);not null;default:0"`
	New               int8           `gorm:"index:new_index;column:new;type:tinyint(4);not null;default:0"`
	Special           int8           `gorm:"index:special_index;column:special;type:tinyint(4);not null;default:0"`
	SortOrder         int64          `gorm:"column:sort_order;type:bigint(20);not null;default:0"`
	CreatedAt         time.Time      `gorm:"column:created_at;type:timestamp"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;type:timestamp"`
}

var (
	db *gorm.DB
	err error
)

func init()  {
	dsn := "root:password@tcp(127.0.0.1:3306)/p_database?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("example/gin-gorm/template/*")

	r.GET("/", func(ctx *gin.Context) {
		var products []Product
		tx := db.Model(&Product{}).Where("manufacturer_id",1)
		p := pagination.GinOrmInit(ctx,tx, func(tx *gorm.DB) *gorm.DB {
			return tx
		},&products,10,5)
		ctx.HTML(200, "index.html", gin.H{"bsPage": p.BsPage(),"products": products})

	})

	r.Run()
	
}
