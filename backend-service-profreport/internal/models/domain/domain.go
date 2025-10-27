package domain

// "time"

// "github.com/google/uuid"

// type BaseModel struct {
// 	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }
// type User struct {
// 	BaseModel
// 	FirstName  string
// 	MiddleName string
// 	SecondName string
// 	Email      string `gorm:"column:email"`
// 	Phone      string `gorm:"column:phone"`
// 	TelegramId string `gorm:"index;column:telegram_id"`
// }

type User struct {
	// ID          string    `json:"id"`
	// CreatedAt   time.Time `json:"created_at"`
	// UpdatedAt   time.Time `json:"updated_at"`
	// FirstName   string    `json:"first_name"`
	// MiddleName  string    `json:"middle_name"`
	// SecondName  string    `json:"second_name"`
	Name  string
	Email string
	// Phone       string    `json:"phone"`
	// TelegramId  string    `json:"telegram_id"`
}
