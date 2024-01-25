package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"oliva-back/internal/models"
	"oliva-back/pkg/config"
)

var (
	db  *gorm.DB
	err error
)

func NewDatabaseConnection(cfg *config.Config) *gorm.DB {
	db, err = gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			cfg.DB.Host,
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.DBName,
			cfg.DB.Port,
			cfg.DB.SSLMode,
			cfg.DB.Timezone)),
		&gorm.Config{DisableForeignKeyConstraintWhenMigrating: false})

	if err != nil {
		panic(err)
	}

	log.Println("Database connection successful")

	db.AutoMigrate(&models.Users{})
	return db
}
