package handler

import (
	"annisa-salon/auth"
	"annisa-salon/helper"
	"annisa-salon/input"
	"annisa-salon/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type treatmentsHandler struct {
	treatmentService service.ServiceTreatments
	authService      auth.UserAuthService
}

func NewTreatmentsHandler(treatmentService service.ServiceTreatments, authService auth.UserAuthService) *treatmentsHandler {
	return &treatmentsHandler{treatmentService, authService}
}


// @Summary Create New treatment
// @Description Create New treatment 
// @Accept json
// @Produce json
// @Tags Treatment
// @Security BearerAuth
// @Param treatment_name formData string true "TreatmentName"
// @Param price formData string true "Price"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Router /api/treatment [post]
func (h *treatmentsHandler) CreateTreatments (c *gin.Context) {
	var inputTreatments input.InputTreatments

	err := c.ShouldBindJSON(&inputTreatments)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	createTreatment, err := h.treatmentService.CreateTreatment(inputTreatments)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, createTreatment)
	c.JSON(http.StatusOK, response)
}

// @Summary Update treatment by slug
// @Description Update treatment by slug 
// @Accept json
// @Produce json
// @Tags Treatment
// @Security BearerAuth
// @Param slug path string true "slug treatment"
// @Param treatment_name formData string true "TreatmentName"
// @Param price formData string true "Price"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Router /api/treatment/{slug} [put]
func (h *treatmentsHandler) UpdatedTreatment(c *gin.Context) {
	slug := c.Param("slug")

	var inputTreatments input.InputTreatments

	err := c.ShouldBindJSON(&inputTreatments)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	UpdatedTreatment, err := h.treatmentService.UpdateTreatment(slug, inputTreatments)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIresponse(http.StatusOK, UpdatedTreatment)
	c.JSON(http.StatusOK, response)
}

// GetAllTreatments 
// @Summary Get All treatment
// @Description Get All treatment
// @Accept json
// @Produce json
// @Tags treatment
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Router /api/treatment [get]
func (h *treatmentsHandler) GetAllTreatments (c *gin.Context){
	// slug := c.Param("slug")
	// finalSlug = c.Param("finalSlug")
	
	Blog, err := h.treatmentService.FindAllTreatment()

	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
        c.JSON(http.StatusBadRequest, response)
        return
	}
	
	response := helper.APIresponse(http.StatusOK, Blog)
	c.JSON(http.StatusOK, response)
}

// GetOneTreatment 
// @Summary Get a single treatment by slug
// @Description Retrieve a single treatment using its slug
// @Tags treatment
// @Accept json
// @Produce json
// @Param slug path string true "Slug of the treatment"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Router /api/treatments/{slug} [get]
func (h *treatmentsHandler) GetOneTreatment (c *gin.Context) {
	slug := c.Param("slug")

	Blog, err := h.treatmentService.FindTreatmentBySlug(slug)

	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
        c.JSON(http.StatusBadRequest, response)
        return
	}
	
	response := helper.APIresponse(http.StatusOK, Blog)
	c.JSON(http.StatusOK, response)
}

// DeleteTreatment 
// @Summary Delete a treatment by slug
// @Description Delete a treatment by its slug
// @Tags treatment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Slug of the treatment to be deleted"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Router /api/treatment/{slug} [delete]
func (h *treatmentsHandler) DeleteTreatment (c *gin.Context) {
	slug := c.Param("slug")

	err := h.treatmentService.DeleteTreatment(slug)

	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
        c.JSON(http.StatusBadRequest, response)
        return
	}
	
	response := helper.APIresponse(http.StatusOK, "blog has succesfuly deleted")
	c.JSON(http.StatusOK, response)
}