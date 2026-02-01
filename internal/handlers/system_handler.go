package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SystemHandler struct {
	db *gorm.DB
}

func NewSystemHandler(db *gorm.DB) *SystemHandler {
	return &SystemHandler{db: db}
}

// UploadImage
// @Summary Upload an image
// @Description Upload an image file (returns a mock URL)
// @Tags system
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image File"
// @Success 200 {object} map[string]string "Image URL"
// @Router /upload/image [post]
func (h *SystemHandler) UploadImage(c *gin.Context) {
	// Mock Image Upload
	c.JSON(http.StatusOK, gin.H{
		"url": "https://via.placeholder.com/300",
		"message": "Image uploaded successfully (Mock)",
	})
}
