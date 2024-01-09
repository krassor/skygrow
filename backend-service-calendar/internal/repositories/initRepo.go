package repositories

import (
	"context"
	"fmt"
	"os"

	"github.com/krassor/skygrow/backend-service-calendar/internal/models/domain"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository() *Repository {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, username, dbName, password)
	fmt.Println(dsn)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error().Msgf("Error gorm.Open(): %s", err)
	}
	log.Info().Msg("gorm have connected to database")

	err = conn.Debug().AutoMigrate(&domain.User{}, &domain.Calendar{}, &domain.CalendarEvent{}, &domain.GoogleAuthToken{}) //Миграция базы данных
	if err != nil {
		log.Error().Msgf("Error gorm.AutoMigrate(): %s", err)
	}
	log.Info().Msg("gorm have connected to database")

	return &Repository{
		DB: conn,
	}
}

func (r *Repository) Shutdown(ctx context.Context) error {
	op := "Repository.Shutdown"
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("Force exit %s: %w", op, ctx.Err())
		default:
			conn, _ := r.DB.DB()
			err := conn.Close()
			if err != nil {
				return fmt.Errorf("error exit %s: %w", op, err)
			}
			return nil
		}
	}
}
