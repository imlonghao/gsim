package db

import (
	"fmt"
	"github.com/imlonghao/gsim/types"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

var DB *gorm.DB

func init() {
	target := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBHOST"), os.Getenv("DBNAME"))
	var err error
	DB, err = gorm.Open("mysql", target)
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&types.Task{})
	DB.AutoMigrate(&types.Token{})
	DB.AutoMigrate(&types.Whitelist{})
	DB.AutoMigrate(&types.Result{})
}
