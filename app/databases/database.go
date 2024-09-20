package databases

import (
	"log"

	"github.com/surajNirala/user_services/app/config"
	"github.com/surajNirala/user_services/app/models"
)

func DatabaseUp() {
	DB := config.DB
	err := DB.AutoMigrate(
		&models.User{},
	)
	// err = DB.AutoMigrate(&models.Student{})
	// err = DB.AutoMigrate(&models.User{}, &models.Post{})
	if err != nil {
		log.Fatalf("Error migrating the database: %v", err)
	}
}
