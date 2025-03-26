package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectBD() error {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", user, password, dbname, host, port, sslmode)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 1 * time.Second,
			LogLevel:      logger.Error,
			Colorful:      true,
		},
	)

	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	// err = DB.AutoMigrate(
	// 	&models.User{},
	// 	&models.Project{},
	// )
	// if err != nil {
	// 	return fmt.Errorf("erro ao realizar a migração do banco de dados: %w", err)
	// }

	log.Println("Banco de dados conectado e migrado com sucesso!")
	return nil
}
