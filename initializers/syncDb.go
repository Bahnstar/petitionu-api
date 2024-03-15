package initializers

import (
	"github.com/Bahnstar/petitionu-api/models"
)

// SyncDb is a function that syncs the database
// by creating the tables for the models
func SyncDb() {
	err := DB.SetupJoinTable(&models.User{}, "Petitions", &models.UserPetition{})
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(&models.Organization{}, &models.Petition{}, &models.Comment{}, &models.Preference{}, &models.User{})
	if err != nil {
		panic(err)
	}
}
