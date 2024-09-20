package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/surajNirala/user_services/app/commons"
	"github.com/surajNirala/user_services/app/config"
	"github.com/surajNirala/user_services/app/models"
	"github.com/surajNirala/user_services/app/validation"
)

var DB = config.DB

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type userResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type sumInterface interface {
	sum() int
	multiply() int
}

type sumStruct struct {
	first  int
	second int
}

func (v sumStruct) multiply() int {
	return v.first * v.second
}
func (v sumStruct) sum() int {
	return v.first + v.second
}

func UserList(c *gin.Context) {

	data := sumStruct{
		first:  4,
		second: 2,
	}
	var s sumInterface = data
	fmt.Println(s.sum())

	var m sumInterface = data
	fmt.Println("Multiple : ", m.multiply())

	var users []models.User
	DB.Select("id", "name", "email").Order("created_at DESC").Find(&users)
	// newUser := make(map[string]interface{})
	var userResponses []userResponse
	for _, item := range users {
		userResponse := userResponse{
			ID:    item.ID,
			Name:  item.Name,
			Email: item.Email,
		}
		userResponses = append(userResponses, userResponse)
	}
	// res := Response{
	// 	Status:  200,
	// 	Message: "Get All User List",
	// 	Data:    userResponses,
	// }
	// c.JSON(200, res)
	commons.ResponseSuccess(c, 200, "Get all user list.", userResponses)
	// return
}

func UserStore(c *gin.Context) {
	var request models.User
	// err := c.BindJSON(&user)

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Request binding error: %v", err)
		customErrors := validation.TranslateValidationErrors(err)
		res := gin.H{
			"status":  400,
			"message": "Invalid Request",
			"errors":  customErrors,
		}
		c.JSON(400, res)
		return
	}

	hashedPassword, err := commons.HashPassword(request.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  500,
			"message": "Failed to hash password.",
		})
		return
	}
	userdata := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
	}

	if err := DB.Create(&userdata).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			// Return a custom error message for duplicate email
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		res := gin.H{
			"status":  500,
			"message": "User not created successfully",
			"error":   err.Error(),
			"data":    nil,
		}
		c.JSON(500, res)
		return
	}
	userResponse := userResponse{
		ID:    userdata.ID,
		Name:  userdata.Name,
		Email: userdata.Email,
	}
	res := gin.H{
		"status":  201,
		"message": "User Created Successfully",
		"data":    userResponse,
	}
	c.JSON(201, res)
}

func UserDetail(c *gin.Context) {
	var user models.User
	user_id := c.Param("user_id")
	result := DB.Select("id", "name", "email").Where("id = ?", user_id).Find(&user)
	if result.RowsAffected == 0 {
		res := Response{
			Status:  409,
			Message: "User not found.",
			Data:    nil,
		}
		c.JSON(409, res)
		return
	}
	userResponse := userResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	res := Response{
		Status:  200,
		Message: "Fetch User Detail",
		Data:    userResponse,
	}
	c.JSON(200, res)
}

func UserUpdate(c *gin.Context) {
	var user models.User
	user_id := c.Param("user_id")
	result := DB.Select("id", "name", "email").Where("id = ?", user_id).Find(&user)
	if result.RowsAffected == 0 {
		res := Response{
			Status:  409,
			Message: "User not found.",
			Data:    nil,
		}
		c.JSON(409, res)
		return
	}

	var request models.User
	// err := c.BindJSON(&user)

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Request binding error: %v", err)
		customErrors := validation.TranslateValidationErrors(err)
		res := gin.H{
			"status":  400,
			"message": "Invalid Request",
			"errors":  customErrors,
		}
		c.JSON(400, res)
		return
	}
	hashedPassword, err := commons.HashPassword(request.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  500,
			"message": "Failed to hash password.",
		})
		return
	}
	userdata := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
	}

	if err := DB.Where("id = ?", user_id).Updates(&userdata).Error; err != nil {
		res := gin.H{
			"status":  500,
			"message": "User is not updated.",
			"error":   err.Error(),
			"data":    nil,
		}
		c.JSON(500, res)
		return
	}
	userID64, _ := strconv.ParseUint(user_id, 10, 64)
	userResponse := userResponse{
		ID:    uint(userID64),
		Name:  userdata.Name,
		Email: userdata.Email,
	}
	res := Response{
		Status:  200,
		Message: "User Detail Updated Successfully.",
		Data:    userResponse,
	}
	c.JSON(200, res)
}

func UserDelete(c *gin.Context) {
	var user models.User
	user_id := c.Param("user_id")
	result := DB.Select("id", "name", "email").Where("id = ?", user_id).Find(&user)
	if result.RowsAffected == 0 {
		res := Response{
			Status:  409,
			Message: "User not found.",
			Data:    nil,
		}
		c.JSON(409, res)
		return
	}

	if err := DB.Where("id = ?", user_id).Delete(&user).Error; err != nil {
		res := gin.H{
			"status":  500,
			"message": "User is not deleted.",
			"error":   err.Error(),
			"data":    nil,
		}
		c.JSON(500, res)
		return
	}

	res := Response{
		Status:  200,
		Message: "User Deleted Successfully.",
		Data:    nil,
	}
	c.JSON(200, res)
}
