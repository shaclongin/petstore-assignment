package app

import "context"

type PetService interface {
	AddPet(ctx context.Context, p Pet) (Pet, error)
	GetPetById(ctx context.Context, id PetId) (Pet, error)
}
