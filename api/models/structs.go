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
  db, err := gorm.Open("mysql", dbUser+":"+dbPass+"@/"+dbName+"?charset=utf8&parseTime=true&loc=Asia%2FTokyo")
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
