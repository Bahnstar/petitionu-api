package initializers

import (
	"github.com/Bahnstar/petitionu-api/models"
)

// SyncDb is a function that syncs the database
// by creating the tables for the models
func SyncDb() {
	DB.AutoMigrate(&models.Organization{}, &models.Petition{}, &models.Comment{}, &models.Preference{}, &models.User{})
}
