package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"app/main.go/internal/models/repositories"

	"github.com/google/uuid"
)

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (repositories.User, error) {
	var userEntity repositories.User
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE email = $1 LIMIT 1`

	err := r.DB.GetContext(ctx, &userEntity, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return repositories.User{}, fmt.Errorf("user not found with email: %s", email)
		}
		return repositories.User{}, fmt.Errorf("error in FindUserByEmail(): %w", err)
	}

	return userEntity, nil
}

func (r *Repository) FindOrCreateUser(ctx context.Context, user repositories.User) (repositories.User, error) {
	op := "calendarRepo.CreateNewUser()"

	// проверка существования пользователя
	var existingUser repositories.User
	checkQuery := `SELECT id, name, email, created_at, updated_at FROM users WHERE email = $1 LIMIT 1`
	err := r.DB.GetContext(ctx, &existingUser, checkQuery, user.Email)

	if err != nil && err != sql.ErrNoRows {
		return repositories.User{}, fmt.Errorf("%s: %w", op, err)
	}

	// если пользователь уже существует, возвращаем его
	if err != sql.ErrNoRows {
		return existingUser, nil
	}

	// создание нового пользователя
	user.ID = uuid.New()

	insertQuery := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3)`

	_, err = r.DB.ExecContext(ctx, insertQuery, user.ID, user.Name, user.Email)
	if err != nil {
		return repositories.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user repositories.User) (repositories.User, error) {
	updateQuery := `UPDATE users SET name = $1, email = $2 WHERE id = $3`

	_, err := r.DB.ExecContext(ctx, updateQuery, user.Name, user.Email, user.ID)
	if err != nil {
		return repositories.User{}, fmt.Errorf("error in UpdateUser(): %w", err)
	}

	return user, nil
}
