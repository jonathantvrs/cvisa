package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jonathantvrs/cvisa/internal/database"
	"github.com/jonathantvrs/cvisa/internal/handler"
	"github.com/jonathantvrs/cvisa/internal/repository"
)

func main() {
	db := database.InitDB()
	repo := &repository.AccountLimitRepository{DB: db}
	h := &handler.AccountLimitHandler{AccountLimitRepo: repo}

	r := gin.Default()

	accountGroup := r.Group("/accounts/:accountId")
	{
		accountGroup.POST("/limit", h.CreateAccountLimit)
		accountGroup.GET("/limit", h.GetAccountLimit)
		accountGroup.PATCH("/limit", h.UpdateAccountLimit)
	}

	r.Run(":8000")
}
