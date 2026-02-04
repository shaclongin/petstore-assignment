package adapters_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tanninio/home-assignment/internal/adapters"
	"github.com/tanninio/home-assignment/internal/app"
	"github.com/tanninio/home-assignment/internal/common"
)

func TestMemRepository(t *testing.T) {
	t.Run("sanity-add-and-get", func(t *testing.T) {
		t.Parallel()
		r := adapters.NewMemRepository()
		p := app.Pet{Id: 1, Name: "first-pet"}
		newp, err := r.AddPet(context.TODO(), p)
		require.NoError(t, err)
		require.Equal(t, p, newp, "want %v, got %v", p, newp)
		gotp, err := r.GetPetById(context.TODO(), p.Id)
		require.NoError(t, err)
		require.Equal(t, p, gotp, "want %v, got %v", p, gotp)
	})
	t.Run("id-conflict", func(t *testing.T) {
		t.Parallel()
		r := adapters.NewMemRepository()
		p := app.Pet{Id: 1, Name: "first-pet"}
		newp, err := r.AddPet(context.TODO(), p)
		require.NoError(t, err)
		require.Equal(t, p, newp, "want %v, got %v", p, newp)
		anotherp := app.Pet{Id: 1, Name: "another-first-pet"}
		_, err = r.AddPet(context.TODO(), anotherp)
		require.ErrorIs(t, err, common.ErrAlreadyExists)
	})
	t.Run("add-and-get-another", func(t *testing.T) {
		t.Parallel()
		r := adapters.NewMemRepository()
		p := app.Pet{Id: 1, Name: "first-pet"}
		newp, err := r.AddPet(context.TODO(), p)
		require.NoError(t, err)
		require.Equal(t, p, newp, "want %v, got %v", p, newp)
		_, err = r.GetPetById(context.TODO(), p.Id+1)
		require.ErrorIs(t, err, common.ErrNotFound)
	})
	t.Run("get-add-get", func(t *testing.T) {
		t.Parallel()
		r := adapters.NewMemRepository()
		p := app.Pet{Id: 1, Name: "first-pet"}
		_, err := r.GetPetById(context.TODO(), p.Id)
		require.ErrorIs(t, err, common.ErrNotFound)
		newp, err := r.AddPet(context.TODO(), p)
		require.NoError(t, err)
		require.Equal(t, p, newp, "want %v, got %v", p, newp)
		gotp, err := r.GetPetById(context.TODO(), p.Id)
		require.NoError(t, err)
		require.Equal(t, p, gotp, "want %v, got %v", p, gotp)
	})
}
