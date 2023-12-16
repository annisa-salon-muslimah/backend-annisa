package handler

import (
	"annisa-salon/auth"
	"annisa-salon/cdn"
	"annisa-salon/helper"
	"annisa-salon/input"
	"annisa-salon/service"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type blogHandler struct {
	blogService service.ServiceBlog
	authService auth.UserAuthService
}

func NewBlogHandler(blogService service.ServiceBlog, authService auth.UserAuthService) *blogHandler {
	return &blogHandler{blogService, authService}
}

// @Summary Create New blog
// @Description Create New blog 
// @Accept json
// @Produce json
// @Tags blog
// @Security BearerAuth
// @Param file formData file true "File gambar"
// @Param title formData string true "Title"
// @Param description formData string true "Description"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Router /api/blog [post]
func (h *blogHandler) CreateBlog(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, "Failed to get file from request")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	src, err := file.Open()
	if err != nil {
		response := helper.APIresponse(http.StatusInternalServerError, "Failed to open file")
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	defer src.Close()
	
	buf:=bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		fmt.Printf("error read file %v",err)
		return 
	}	

	img,err:=cdn.Base64toEncode(buf.Bytes())
	if err!=nil{
		fmt.Printf("error reading image %v",err)
	}

	fmt.Printf("image base 64 format : %v",img)

	imageKitURL, err := cdn.ImageKit(context.Background(), img)
	if err != nil {
		// Tangani jika terjadi kesalahan saat upload gambar
		// Misalnya, Anda dapat mengembalikan respon error ke klien jika diperlukan
		response := helper.APIresponse(http.StatusInternalServerError, "Failed to upload image")
		c.JSON(http.StatusInternalServerError, response)
		return
	}


	var inputBlog input.InputBlog

	err = c.ShouldBind(&inputBlog)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIresponse(http.StatusUnprocessableEntity, errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	if err != nil {
		//inisiasi data yang tujuan dalam return hasil ke postman
		data := gin.H{"is_uploaded": false}
		response := helper.APIresponse(http.StatusUnprocessableEntity, data)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	_, err = h.blogService.CreateBlog(inputBlog, imageKitURL)
	if err != nil {
		// data := gin.H{"is_uploaded": false}
		response := helper.APIresponse(http.StatusUnprocessableEntity, err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	data := gin.H{"is_uploaded": true}
	response := helper.APIresponse(http.StatusOK, data)
	c.JSON(http.StatusOK, response)

}

// @Summary Update blog by slug
// @Description Update blog by slug 
// @Accept json
// @Produce json
// @Tags blog
// @Security BearerAuth
// @Param slug path string true "slug blog"
// @Param file formData file true "File gambar"
// @Param title formData string true "Title"
// @Param description formData string true "Description"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Router /api/blog/{slug} [put]
func (h *blogHandler) UpdateBlog (c *gin.Context){

	file, err := c.FormFile("file")
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, "Failed to get file from request")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	src, err := file.Open()
	if err != nil {
		response := helper.APIresponse(http.StatusInternalServerError, "Failed to open file")
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	defer src.Close()
	
	buf:=bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		fmt.Printf("error read file %v",err)
		return 
	}	

	img,err:=cdn.Base64toEncode(buf.Bytes())
	if err!=nil{
		fmt.Printf("error reading image %v",err)
	}

	fmt.Printf("image base 64 format : %v",img)

	imageKitURL, err := cdn.ImageKit(context.Background(), img)
	if err != nil {
		// Tangani jika terjadi kesalahan saat upload gambar
		// Misalnya, Anda dapat mengembalikan respon error ke klien jika diperlukan
		response := helper.APIresponse(http.StatusInternalServerError, "Failed to upload image")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var updateInput input.InputBlog

	err = c.ShouldBind(&updateInput)
	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, "Invalid input")
        c.JSON(http.StatusBadRequest, response)
        return
	}

	slug := c.Param("slug")
	// finalSlug = c.Param("finalSlug")
	
	_, err = h.blogService.UpdateBlog(slug, updateInput, imageKitURL)

	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
        c.JSON(http.StatusBadRequest, response)
        return
	}
	
	data := gin.H{"is_updated": true}
	response := helper.APIresponse(http.StatusOK, data)
	c.JSON(http.StatusOK, response)
}

// GetAllBlog 
// @Summary Get all blogs
// @Description Retrieve all blogs with pagination
// @Tags Blog
// @Accept json
// @Produce json
// @Param page query integer false "Page number for pagination (default is 1)"
// @Param limit query integer false "Number of items per page (default is 10)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Router /api/blog [get]
func (h *blogHandler) GetAllBlog(c *gin.Context) {
    page, err := strconv.Atoi(c.Query("page"))
    if err != nil || page <= 0 {
        page = 1
    }

    limit, err := strconv.Atoi(c.Query("limit"))
    if err != nil || limit <= 0 {
        limit = 10
    }

    blogs, err := h.blogService.FindAllBlog(page, limit)
    if err != nil {
        response := helper.APIresponse(http.StatusBadRequest, err.Error())
        c.JSON(http.StatusBadRequest, response)
        return
    }

    response := helper.APIresponse(http.StatusOK, blogs)
    c.JSON(http.StatusOK, response)
}

// GetOneBlog 
// @Summary Get a single blog by slug
// @Description Retrieve a single blog using its slug
// @Tags Blog
// @Accept json
// @Produce json
// @Param slug path string true "Slug of the blog"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Router /api/blog/{slug} [get]
func (h *blogHandler) GetOneBlog (c *gin.Context) {
	slug := c.Param("slug")

	Blog, err := h.blogService.FindBlogBySlug(slug)

	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
        c.JSON(http.StatusBadRequest, response)
        return
	}
	
	response := helper.APIresponse(http.StatusOK, Blog)
	c.JSON(http.StatusOK, response)
}

// DeleteBlog 
// @Summary Delete a blog by slug
// @Description Delete a blog by its slug
// @Tags Blog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Slug of the blog to be deleted"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Router /api/blog/{slug} [delete]
func (h *blogHandler) DeleteBlog (c *gin.Context) {
	slug := c.Param("slug")

	err := h.blogService.DeleteBlog(slug)

	if err != nil {
		response := helper.APIresponse(http.StatusBadRequest, err.Error())
        c.JSON(http.StatusBadRequest, response)
        return
	}
	
	response := helper.APIresponse(http.StatusOK, "blog has succesfuly deleted")
	c.JSON(http.StatusOK, response)
}

