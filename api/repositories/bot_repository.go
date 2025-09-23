// repositories/bot_repository.go
package repositories

import (
	"gorm.io/gorm"

	"github.com/brando1998/docubot-api/models"
)

type BotRepository interface {
	GetBotByNumber(number string) (*models.Bot, error)
	GetOrCreateBot(number string, name string) (*models.Bot, error)
}

type botRepository struct {
	db *gorm.DB
}

func NewBotRepository(db *gorm.DB) BotRepository {
	return &botRepository{db}
}

func (r *botRepository) GetBotByNumber(number string) (*models.Bot, error) {
	var bot models.Bot
	err := r.db.Where("number = ?", number).First(&bot).Error
	return &bot, err
}

// get or create bot by number
func (r *botRepository) GetOrCreateBot(number string, name string) (*models.Bot, error) {
	var bot models.Bot
	err := r.db.Where("number = ?", number).First(&bot).Error
	if err == gorm.ErrRecordNotFound {
		bot = models.Bot{
			Number: number,
			Name:   name,
		}
		err = r.db.Create(&bot).Error
	}
	return &bot, err
}
