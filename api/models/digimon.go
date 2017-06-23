package model

import (
  "encoding/json"
)

// DigimonRepository デジモンリポジトリ
type DigimonRepository struct {}

// NewDigimonRepository リポジトリ取得
func NewDigimonRepository() DigimonRepository {
  return DigimonRepository{}
}

// GetAll 全デジモンを取得
func (m DigimonRepository) GetAll() *[]Digimon {
  db := connect()
  digimons := []Digimon{}
  db.Find(&digimons)
  for i := range digimons {
    json.Unmarshal([]byte(digimons[i].Data), &digimons[i].DataObj)
  }
  db.Close()
  return &digimons
}

// GetByID idでデジモンを取得
func (m DigimonRepository) GetByID(id int) *Digimon {
  db := connect()
  var ret *Digimon
  digimon := Digimon{}
  db.Find(&digimon, "id = ?", id)
  db.Close()
  if digimon.ID != 0 {
    json.Unmarshal([]byte(digimon.Data), &digimon.DataObj)
    ret = &digimon
  }
  return ret
}

// SetNewDigimon 新規デジモン作成
func (m DigimonRepository) SetNewDigimon(name string, data *[][]int) *Digimon {
  db := connect()
  digimon := Digimon{}
  digimon.Name = name
  bytes, _ := json.Marshal(data)
  digimon.Data = string(bytes)
  json.Unmarshal([]byte(digimon.Data), &digimon.DataObj)
  db.Create(&digimon)
  db.Close()
  return &digimon
}
