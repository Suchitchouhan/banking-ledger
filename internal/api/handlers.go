package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"banking-ledger/internal/models"
	"banking-ledger/internal/service"
)

type Handler struct {
	accountService     *service.AccountService
	transactionService *service.TransactionService
}

func NewHandler(
	accountService *service.AccountService,
	transactionService *service.TransactionService,
) *Handler {
	return &Handler{
		accountService:     accountService,
		transactionService: transactionService,
	}
}

// Creation of a new account
func (h *Handler) CreateAccountHandler(c *gin.Context) {
	var req models.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	account, err := h.accountService.CreateAccount(req.Name, req.InitialAmount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    account,
	})
}

// Retrieves an account by ID
func (h *Handler) GetAccountHandler(c *gin.Context) {
	accountID := c.Param("id")

	account, err := h.accountService.GetAccount(accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Account not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    account,
	})
}

// Lists all accounts
func (h *Handler) ListAccountsHandler(c *gin.Context) {
	accounts, err := h.accountService.ListAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve accounts",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    accounts,
	})
}

// Deposits to account
func (h *Handler) DepositHandler(c *gin.Context) {
	accountID := c.Param("id")

	var req models.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	transaction, err := h.transactionService.CreateDeposit(accountID, req.Amount, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    transaction,
	})
}

// Withdrawals from account
func (h *Handler) WithdrawHandler(c *gin.Context) {
	accountID := c.Param("id")

	var req models.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	transaction, err := h.transactionService.CreateWithdrawal(accountID, req.Amount, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    transaction,
	})
}
