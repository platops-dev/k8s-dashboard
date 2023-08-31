package db

import (
	"database/sql"
	"fmt"
	"test4/config"
	"test4/model"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wonderivan/logger"
)

var (
	isInit bool
	GORM   *gorm.DB
	err error
)

//db的初始化函数, 与数据库建立连接
func Init()  {
	//判断是否已经初始化了
	if isInit {
		return
	}

	//判断数据库是否存在, 不存在则创建
	db, errs := sql.Open(config.DbType, fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&parseTime=True&loc=Local", 
		config.DbUser, config.DbPwd, config.DbHost, config.DbPort))
	if errs != nil {
		panic("数据库连接失败: " + errs.Error())
	}
	defer db.Close()
	//创建数据库(如果不存在)
	if _, errs := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARSET UTF8", config.DbName)); errs != nil {
		panic("创建数据库失败: " + errs.Error())
	}

	//组装连接配置
	//ParseTime 是查询结果是否自动解析为时间
	// loc是mysql的时区配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
	config.DbUser,
	config.DbPwd,
	config.DbHost,
	config.DbPort,
	config.DbName)

	//与数据库建立连接, 生成一个*gorm.BD类型的对象
	GORM, err = gorm.Open(config.DbType, dsn)
	if err != nil {
		panic("数据库连接失败. " + err.Error())
	}

	//打印sql语句
	GORM.LogMode(config.LogMode)

	//迁移数据表
	GORM.Set("gorm:table_options", "CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ENGINE=InnoDB").AutoMigrate(&model.Workflow{})
	logger.Info("自动迁移数据库表成功")

	//开启连接池
	//连接池最大允许的空闲连接数, 如果sql任务需要执行的连接数大于20, 超过的连接数会被连接池关闭
	GORM.DB().SetMaxIdleConns(config.MaxIdleConns)
	//设置连接可复用的最大连接时间
	GORM.DB().SetMaxOpenConns(config.MaxOpenConns)
	GORM.DB().SetConnMaxLifetime(time.Duration(config.MaxLifeTime))

	isInit = true
	logger.Info("连接数据库成功")
}