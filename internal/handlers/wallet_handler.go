package handlers

import (
	"net/http"
	"raijai-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletHandler struct {
	db *gorm.DB
}

func NewWalletHandler(db *gorm.DB) *WalletHandler {
	return &WalletHandler{db: db}
}

// CreateWallet godoc
// @Summary Create a new wallet
// @Description Create a new wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param wallet body models.Wallet true "Wallet Data"
// @Success 201 {object} models.Wallet
// @Router /wallets [post]
// CreateWallet godoc
// @Summary Create a new wallet
// @Description Create a new wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param wallet body models.Wallet true "Wallet Data"
// @Success 201 {object} models.Wallet
// @Router /wallets [post]
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var wallet models.Wallet
	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet.ID = uuid.New().String()
	wallet.OwnerID = c.MustGet("user_id").(string)

	if result := h.db.Create(&wallet); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, wallet)
}

// GetWallets godoc
// @Summary Get all wallets by User ID
// @Description Get all wallets specific to a user
// @Tags wallets
// @Produce json
// @Success 200 {array} models.Wallet
// @Router /wallets [get]
func (h *WalletHandler) GetWallets(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	var wallets []models.Wallet

	if result := h.db.Where("owner_id = ?", userID).Find(&wallets); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, wallets)
}

// UpdateWallet godoc
// @Summary Update a wallet
// @Description Update a wallet by ID
// @Tags wallets
// @Accept json
// @Produce json
// @Param id path string true "Wallet ID"
// @Param wallet body models.Wallet true "Wallet Data"
// @Success 200 {object} models.Wallet
// @Router /wallets/{id} [put]
func (h *WalletHandler) UpdateWallet(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")
	var wallet models.Wallet

	if result := h.db.Where("id = ? AND owner_id = ?", id, userID).First(&wallet); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID and OwnerID are not changed
	wallet.ID = id
	wallet.OwnerID = userID

	h.db.Save(&wallet)
	c.JSON(http.StatusOK, wallet)
}

// DeleteWallet godoc
// @Summary Delete a wallet
// @Description Delete a wallet by ID
// @Tags wallets
// @Produce json
// @Param id path string true "Wallet ID"
// @Success 200 {object} map[string]string
// @Router /wallets/{id} [delete]
func (h *WalletHandler) DeleteWallet(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")
	if result := h.db.Where("id = ? AND owner_id = ?", id, userID).Delete(&models.Wallet{}, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet deleted successfully"})
}
