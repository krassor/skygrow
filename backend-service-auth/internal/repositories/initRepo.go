package repositories

import (
	"fmt"
	"github.com/krassor/skygrow/backend-service-auth/internal/config"
	"github.com/krassor/skygrow/backend-service-auth/internal/models/domain"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(cfg *config.Config) *Repository {
	username := cfg.DBConfig.User
	password := cfg.DBConfig.Password
	dbName := cfg.DBConfig.Name
	dbHost := cfg.DBConfig.Host
	dbPort := cfg.DBConfig.Port

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, username, dbName, password)
	fmt.Println(dsn)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error().Msgf("Error gorm.Open(): %s", err)
	}
	log.Info().Msg("gorm have connected to database")

	err = conn.Debug().AutoMigrate(&domain.User{}) //Миграция базы данных
	if err != nil {
		log.Error().Msgf("Error gorm.AutoMigrate(): %s", err)
	}
	log.Info().Msg("gorm have connected to database")

	return &Repository{
		DB: conn,
	}
}
