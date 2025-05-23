package controllers

import (
	"api-service/services"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type APIController struct {
	APIService services.APIService
}

func NewAPIController(APIService services.APIService) *APIController {
	return &APIController{APIService: APIService}
}

func (ctrl *APIController) Collect(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":  "error",
			"message": "Method not allowed",
		})
		return
	}
	
	extracted_id := c.Query("id")
	if extracted_id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "ID parameter is require.*",
		})
		return
	}

	id, err := strconv.Atoi(extracted_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid ID format, Id must be integer.",
		})
		return
	}

	data, err := ctrl.APIService.GetProductById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch subject list",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}