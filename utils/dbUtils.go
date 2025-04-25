package utils

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

func init() {
	dsn := "root:duoguan@mysql@123@(47.110.151.225:3307)/eth?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 开启单数表名
		},
	})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := Db.DB()
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
}
