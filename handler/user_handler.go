package handler

import (
	"annisa-salon/auth"
	"annisa-salon/helper"
	"annisa-salon/input"
	"annisa-salon/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

	type userHandler struct {
		userService service.ServiceUser
		authService auth.UserAuthService
	}

	func NewUserHandler(userService service.ServiceUser, authService auth.UserAuthService) *userHandler {
		return &userHandler{userService, authService}
	}

	// @Summary Register new user
	// @Description Register a new user with the provided information
	// @Accept json
	// @Produce json
	// @Tags Users
	// @Param body body input.UserInput true "User registration details"
	// @Success 200 {object} map[string]interface{}
	// @Failure 400 {object} map[string]interface{}
	// @Failure 409 {object} map[string]interface{}
	// @Failure 422 {object} map[string]interface{}
	// @Failure 500 {object} map[string]interface{}
	// @Router /api/user/register [post]
	func (h *userHandler) RegisterUser(c *gin.Context) {
		var input input.UserInput

		err := c.ShouldBindJSON(&input)
		if err != nil {
			errors := helper.FormatValidationError(err)
			errorMessage := gin.H{"errors": errors}
			response := helper.APIresponse(http.StatusUnprocessableEntity, errorMessage)
			c.JSON(http.StatusUnprocessableEntity, response)
			return
		}

		isEmailAvailable, err := h.userService.IsEmailAvailability(input.Email)
		if err != nil {
		    response := helper.APIresponse(http.StatusInternalServerError, nil)
		    c.JSON(http.StatusInternalServerError, response)
		    return
		}

		// Jika email tidak tersedia, kirim respons kesalahanx
		if !isEmailAvailable {
		    response := helper.APIresponse(http.StatusConflict, nil)
		    c.JSON(http.StatusConflict, response)
		    return
		}

		newUser, err := h.userService.RegisterUser(input)
		if err != nil {
			response := helper.APIresponse(http.StatusUnprocessableEntity, "erorrrr")
			c.JSON(http.StatusUnprocessableEntity, response)
			return
		}

		response := helper.APIresponse(http.StatusOK, newUser)
		c.JSON(http.StatusOK, response)
	}

	// @Summary User login
	// @Description Log in an existing user using email and password
	// @Accept json
	// @Produce json
	// @Tags Users
	// @Param body body input.UserInput true "User login details"
	// @Success 200 {object} map[string]interface{}
	// @Failure 400 {object} map[string]interface{}
	// @Failure 422 {object} map[string]interface{}
	// @Failure 500 {object} map[string]interface{}
	// @Router /api/user/login [post]
func (h *userHandler) Login(c *gin.Context) {
	var input input.UserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIresponse(http.StatusUnprocessableEntity, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.LoginUser(input)
	if err != nil {
		response := helper.APIresponse(http.StatusUnprocessableEntity, nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	token, err := h.authService.GenerateToken(loggedinUser.ID, loggedinUser.Role)
	if err != nil {
		response := helper.APIresponse(http.StatusUnprocessableEntity, nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// formatter := user.FormatterUser(loggedinUser, token)
	response := helper.APIresponse(http.StatusOK, token)
	c.JSON(http.StatusOK, response)
}