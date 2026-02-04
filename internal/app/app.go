package app

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tanninio/home-assignment/internal/common"
)

type Application struct {
	repo PetRepository
}

const minPetNameLength = 4

func NewApplication(repo PetRepository) Application {
	return Application{repo: repo}
}

func (a Application) AddPet(ctx context.Context, p Pet) (Pet, error) {
	logrus.Infof("Adding pet with name: %s", p.Name)
	if len(p.Name) < minPetNameLength {
		return Pet{}, fmt.Errorf("can't add pet with short name: %w", common.ErrIncorrectInput)
	}
	return a.repo.AddPet(ctx, p)
}

func (a Application) GetPetById(ctx context.Context, id PetId) (Pet, error) {
	logrus.Infof("Getting pet with id: %d", id)
	return a.repo.GetPetById(ctx, id)
}
