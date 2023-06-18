package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
var DB *gorm.DB

func LoadDB() {
	var err error
	db_string:=os.Getenv("DB_STRING")
	dsn := db_string
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
   if(err!=nil){
	log.Fatal("Failed to connect to DB")
   }
}
