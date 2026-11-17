package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/laineus/digimon/api/models"
	"net/http"
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
	Name string  `json:"name"`
	Data [][]int `json:"data"`
}

// PostDigimon デジモン投稿api
func PostDigimon(c *gin.Context) {
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	chk := checkSize(&request.Data, 16, 16)
	if chk {
		if request.Name == "" {
			request.Name = "no name"
		}
		repo := model.NewDigimonRepository()
		if err := repo.SetNewDigimon(request.Name, &request.Data); err != nil {
			if errors.Is(err, model.ErrPostRateLimited) {
				c.Header("Retry-After", strconv.Itoa(model.PostIntervalSeconds))
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "投稿は全体で60秒に1件までです"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save post"})
			}
			return
		}
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
