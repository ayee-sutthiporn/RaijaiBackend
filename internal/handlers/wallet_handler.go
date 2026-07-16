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
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var wallet models.Wallet
	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(string)
	wallet.ID = uuid.New().String()
	wallet.OwnerID = userID

	// Convert empty string pointer for BookID to nil
	if wallet.BookID != nil && *wallet.BookID == "" {
		wallet.BookID = nil
	}

	if wallet.BookID != nil && !canEditInBook(h.db, *wallet.BookID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to add wallets to this book"})
		return
	}

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
// @Param book_id query string false "Filter by Book ID"
// @Success 200 {array} models.Wallet
// @Router /wallets [get]
func (h *WalletHandler) GetWallets(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	bookID := c.Query("book_id")
	var wallets []models.Wallet

	var query *gorm.DB
	if bookID != "" {
		if !isBookMember(h.db, bookID, userID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this book"})
			return
		}
		// Book members share visibility of all wallets in the book, not just their own.
		query = h.db.Where("book_id = ?", bookID)
	} else {
		query = h.db.Where("owner_id = ?", userID)
	}

	if result := query.Find(&wallets); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, wallets)
}

// GetWallet godoc
// @Summary Get a wallet by ID
// @Description Get a wallet by ID
// @Tags wallets
// @Produce json
// @Param id path string true "Wallet ID"
// @Success 200 {object} models.Wallet
// @Router /wallets/{id} [get]
func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")
	var wallet models.Wallet

	if result := h.db.First(&wallet, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	isOwner := wallet.OwnerID == userID
	isSharedMember := wallet.BookID != nil && isBookMember(h.db, *wallet.BookID, userID)
	if !isOwner && !isSharedMember {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, wallet)
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

	if result := h.db.First(&wallet, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	isOwner := wallet.OwnerID == userID
	canEditViaBook := wallet.BookID != nil && canEditInBook(h.db, *wallet.BookID, userID)
	if !isOwner && !canEditViaBook {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to edit this wallet"})
		return
	}

	originalOwnerID := wallet.OwnerID
	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID and original owner are not changed
	wallet.ID = id
	wallet.OwnerID = originalOwnerID

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

	var wallet models.Wallet
	if result := h.db.First(&wallet, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	isOwner := wallet.OwnerID == userID
	canEditViaBook := wallet.BookID != nil && canEditInBook(h.db, *wallet.BookID, userID)
	if !isOwner && !canEditViaBook {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this wallet"})
		return
	}

	if result := h.db.Delete(&models.Wallet{}, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet deleted successfully"})
}
