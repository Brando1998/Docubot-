package repositories

import (
	"errors"
	"log"

	"gorm.io/gorm"

	"github.com/brando1998/docubot-api/models"
)

type ClientRepository interface {
	CreateClient(user *models.Client) error
	GetClientByID(id uint) (*models.Client, error)
	GetClientByPhone(phone string) (*models.Client, error)
	GetOrCreateClient(phone, name, email string) (*models.Client, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateClient(user *models.Client) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetClientByID(id uint) (*models.Client, error) {
	var user models.Client
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) GetClientByPhone(phone string) (*models.Client, error) {
	var user models.Client
	err := r.db.Where("phone = ?", phone).First(&user).Error
	return &user, err
}

func (r *userRepository) GetOrCreateClient(phone, name, email string) (*models.Client, error) {
	var user models.Client
	err := r.db.Where("phone = ?", phone).First(&user).Error

	if err == nil {
		return &user, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Cliente no encontrado, creando nuevo con phone: %s", phone)

		var emailPtr *string
		if email != "" {
			emailPtr = &email
		}

		user = models.Client{
			Phone: phone,
			Name:  name,
			Email: emailPtr,
		}

		if createErr := r.db.Create(&user).Error; createErr != nil {
			log.Printf("Error creando cliente: %v", createErr)
			return nil, createErr
		}
		return &user, nil
	}

	log.Printf("Error buscando cliente: %v", err)
	return nil, err
}
