package interfaces

import (
	"audigo/models"
)

type Storage interface {
	Upload(bytes []byte, filename string) (string, string, error)
	SaveMetadata(audio *models.Audio) error
	List() ([]*models.Audio, error)
	GetByID(id string) (*models.Audio, error)
	Delete(id string) error
}
