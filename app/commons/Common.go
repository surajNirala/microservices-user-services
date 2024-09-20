package commons

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/surajNirala/user_services/app/config"
	"github.com/surajNirala/user_services/app/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SuccessResponseOrder struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponseOrder struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"` // Make error field optional
}

func ResponseSuccess(c *gin.Context, status int, message string, data interface{}) {
	res := SuccessResponseOrder{
		Status:  status,
		Message: message,
	}
	if data != nil {
		res.Data = data
	}
	c.JSON(status, res)
}

// ResponseError formats error API responses
func ResponseError(c *gin.Context, status int, message string, optionalErr ...error) {
	res := ErrorResponseOrder{
		Status:  status,
		Message: message,
	}
	if len(optionalErr) > 0 && optionalErr[0] != nil {
		res.Error = optionalErr[0] // Set error only if it's not nil
	}
	c.JSON(status, res)
	c.Abort()
}

var DB = config.DB

func CheckDuplicateEmail(email string) bool {
	var user models.User
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		log.Printf("Error checking for duplicate email: %v", err)
	}
	return true
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
