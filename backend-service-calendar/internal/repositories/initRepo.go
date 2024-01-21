package repositories

import (
	"context"
	"fmt"
	"github.com/krassor/skygrow/backend-service-calendar/internal/config"
	"github.com/krassor/skygrow/backend-service-calendar/internal/models/domain"
	"github.com/krassor/skygrow/backend-service-calendar/internal/utils/logger/sl"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

type Repository struct {
	DB  *gorm.DB
	log *slog.Logger
}

func NewCalendarRepository(log *slog.Logger, cfg *config.Config) *Repository {
	op := "repositories.NewCalendarRepository()"
	log.With(
		slog.String("op", op))

	username := cfg.DBConfig.User
	password := cfg.DBConfig.Password
	dbName := cfg.DBConfig.Name
	dbHost := cfg.DBConfig.Host
	dbPort := cfg.DBConfig.Port

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, username, dbName, password)
	fmt.Println(dsn)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error("error connecting to database", sl.Err(err))
		panic("error connecting to database")
	}
	log.Debug("gorm have connected to database")

	err = conn.Debug().AutoMigrate(&domain.User{}, &domain.Calendar{}, &domain.CalendarEvent{}, &domain.GoogleAuthToken{}) //Миграция базы данных
	if err != nil {
		log.Error("error database auto migrate", sl.Err(err))
		panic("error database auto migrate")
	}
	log.Debug("success auto migrate")

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
