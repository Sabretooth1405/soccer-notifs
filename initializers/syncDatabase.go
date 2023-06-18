package initializers

import "soccer-notifs/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
