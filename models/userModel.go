package models

import "gorm.io/gorm"


type User struct{
	gorm.Model
	Email string `gorm:"unique"`
	Email_is_verified bool 
	Password string
	Timezone string `gorm:"default:IST"`
}
