package handler

import (
	"net/http"
	"strconv"

	"github.com/jonathantvrs/cvisa/internal/model"
	"github.com/jonathantvrs/cvisa/internal/repository"

	"github.com/gin-gonic/gin"
)

type AccountLimitHandler struct {
	AccountLimitRepo *repository.AccountLimitRepository
}

func (h *AccountLimitHandler) CreateAccountLimit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("accountId"))
	var input struct {
		Limit int64 `json:"available_credit_limit"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}
	newLimit := model.AccountLimit{
		AccountID:            uint(id),
		AvailableCreditLimit: input.Limit,
	}
	if err := h.AccountLimitRepo.CreateAccountLimit(&newLimit); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Limite de Conta já existe"})
		return
	}
	c.JSON(http.StatusCreated, newLimit)
}

func (h *AccountLimitHandler) GetAccountLimit(c *gin.Context) {
	acc_id, _ := strconv.Atoi(c.Param("accountId"))
	acc_limit, err := h.AccountLimitRepo.GetLimitByAccountId(uint(acc_id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conta não encontrada"})
		return
	}
	c.JSON(http.StatusOK, acc_limit)
}

func (h *AccountLimitHandler) UpdateAccountLimit(c *gin.Context) {
	acc_id, _ := strconv.Atoi(c.Param("accountId"))
	var input struct {
		Amount int64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valor inválido"})
		return
	}

	if err := h.AccountLimitRepo.UpdateAccountLimit(uint(acc_id), input.Amount); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Limite atualizado com sucesso!"})
}
