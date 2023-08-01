package repositories

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/krassor/skygrow/backend-serverHttp/internal/models/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	DB *gorm.DB
}

func NewRepository() *repository {
	username := os.Getenv("ORDERS_DB_USER")
	password := os.Getenv("ORDERS_DB_PASSWORD")
	dbName := os.Getenv("ORDERS_DB_NAME")
	dbHost := os.Getenv("ORDERS_DB_HOST")
	dbPort := os.Getenv("ORDERS_DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, username, dbName, password)
	fmt.Println(dsn)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error().Msgf("Error gorm.Open(): %s", err)
	}
	log.Info().Msg("gorm have connected to database")

	err = conn.Debug().AutoMigrate(&entities.BookOrder{}, &entities.Subscriber{}, &entities.User{}, &entities.Mentor{}) //Миграция базы данных
	if err != nil {
		log.Error().Msgf("Error gorm.AutoMigrate(): %s", err)
	}
	log.Info().Msg("gorm have connected to database")

	return &repository{
		DB: conn,
	}
}
