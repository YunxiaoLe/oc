package dbclt

import (
	"context"
	"fmt"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/conf"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"time"
)

var Engine *gorm.DB

func CreateTableIfNotExist(Engine *gorm.DB, tableModels []interface{}) {
	for _, value := range tableModels {
		err := Engine.AutoMigrate(value)
		if err != nil {
			fmt.Println("Create table ", reflect.TypeOf(value), " error!")
		}
	}
}
func dsn(settings conf.DbSettings) string {
	// https://stackoverflow.com/questions/45040319/unsupported-scan-storing-driver-value-type-uint8-into-type-time-time
	// Add ?parseTime=true
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4,utf8", settings.Username, settings.Password, settings.Hostname, settings.Dbname)
}

// InitDb 依据配置文件信息连接数据库
func InitDb(settings conf.DbSettings) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(10*time.Second))
	go func(ctx context.Context) {
		connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4,utf8",
			settings.Username, settings.Password, settings.Hostname, settings.Dbname)
		var err1 error

		Engine, err1 = gorm.Open(mysql.Open(connStr), &gorm.Config{})
		if err1 != nil {
			panic(any("Database connect error," + err1.Error()))
		}
		sqlDB, err := Engine.DB()
		if err != nil {
			panic(any("Database error"))
		}
		var temp []interface{}
		var user model.AuthUser
		//var college models.College
		//var branch models.Branch
		//var admin models.Admin
		temp = append(temp, &user)
		//temp = append(temp, &user, &college, &branch, &admin)
		CreateTableIfNotExist(Engine, temp)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(10000)
		sqlDB.SetConnMaxLifetime(time.Second * 3)
		cancel()
	}(ctx)

	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			fmt.Println("context timeout exceeded")
			panic(any("Timeout when initialize database connection"))
		case context.Canceled:
			fmt.Println("context cancelled by force. whole process is complete")
		}
	}
}
