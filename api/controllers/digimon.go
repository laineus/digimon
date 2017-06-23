package controller

import (
  "../models"
  "github.com/gin-gonic/gin"
  "strconv"
)

// GetDigimon デジモン取得api
func GetDigimon(c *gin.Context) {
  digimonID, e := strconv.Atoi(c.Query("id"))
  repo := model.NewDigimonRepository()
  if e == nil {
    digimon := repo.GetByID(digimonID)
    c.JSON(200, digimon)
  } else {
    digimons := repo.GetAll()
    c.JSON(200, digimons)
  }
}

// Request リクエスト形式を定義する構造体
type Request struct {
  Name string `json:"name"`
  Data [][]int `json:"data"`
}

// PostDigimon デジモン投稿api
func PostDigimon(c *gin.Context) {
  var request Request
  c.BindJSON(&request)
  chk := checkSize(&request.Data, 16, 16)
  if chk {
    if request.Name == "" {
      request.Name = "no name"
    }
    repo := model.NewDigimonRepository()
    repo.SetNewDigimon(request.Name, &request.Data)
    digimons := repo.GetAll()
    c.JSON(200, digimons)
  } else {
    c.JSON(200, false)
  }
}

func checkSize(dots *[][]int, width, height int) bool {
  chk := true
  if len(*dots) != height {
    chk = false
  } else {
    for _, row := range *dots {
      if len(row) != width {
        chk = false
        break
      }
    }
  }
  return chk
}
