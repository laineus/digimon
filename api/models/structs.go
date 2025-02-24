package model

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  "os"
  "time"
)

// .envファイルの情報でmysqlに接続
func connect() *gorm.DB {
  dbUser := os.Getenv("DB_USER")
  dbPass := os.Getenv("DB_PASS")
  dbName := os.Getenv("DB_NAME")
  dbHost := os.Getenv("DB_HOST")
  dbPort := os.Getenv("DB_PORT")
  db, err := gorm.Open("mysql", dbUser+":"+dbPass+"@tcp(" + dbHost + ":" + dbPort + ")/"+dbName+"?charset=utf8&parseTime=true&loc=Asia%2FTokyo&allowNativePasswords=true")
  if err != nil {
    panic("failed to connect database")
  }
  return db
}

// TimeStamp 共通の日付カラム
type TimeStamp struct {
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt *time.Time
}

// Digimon デジモン
type Digimon struct {
  ID int
  Name string
  Data string
  DataObj *[][]int `gorm:"-"`
  TimeStamp
}
