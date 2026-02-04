package adapters

import (
	"context"
	"fmt"
	"sync"

	"github.com/tanninio/home-assignment/internal/app"
	"github.com/tanninio/home-assignment/internal/common"
)

type MemRepository struct {
	pets  map[app.PetId]app.Pet
	mutex *sync.RWMutex
}

func NewMemRepository() *MemRepository {
	return &MemRepository{
		pets:  map[app.PetId]app.Pet{},
		mutex: &sync.RWMutex{},
	}
}

func (r *MemRepository) AddPet(ctx context.Context, p app.Pet) (app.Pet, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, ok := r.pets[p.Id]; ok {
		return app.Pet{}, fmt.Errorf("can't add pet with id %d: %w", p.Id, common.ErrAlreadyExists)
	}
	r.pets[p.Id] = p
	return p, nil
}

func (r *MemRepository) GetPetById(ctx context.Context, id app.PetId) (app.Pet, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	p, ok := r.pets[id]
	if !ok {
		return app.Pet{}, fmt.Errorf("can't get pet with id %d: %w", id, common.ErrNotFound)
	}
	return p, nil
}
