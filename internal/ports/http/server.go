package ports

import (
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/tanninio/home-assignment/internal/app"
	"github.com/tanninio/home-assignment/internal/common"
)

type HttpServer struct {
	svc app.PetService
}

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -package ports -generate strict-server,gorilla,types -o openapi.gen.go ../../../api/petstore.yml
func NewHttpServer(svc app.PetService) HttpServer {
	return HttpServer{svc: svc}
}

func unimplementedMethod(method string) error {
	return fmt.Errorf("%s: %w", method, common.ErrUnimplemented)
}

func (h HttpServer) AddPet(ctx context.Context, request AddPetRequestObject) (AddPetResponseObject, error) {
	var resp AddPet200JSONResponse
	var addent app.Pet
	if err := copier.Copy(&addent, request.Body); err != nil {
		return nil, common.ErrIncorrectInput
	}
	newent, err := h.svc.AddPet(ctx, addent)
	if err != nil {
		return nil, fmt.Errorf("can't add pet: %w", err)
	}
	if err := copier.Copy(&resp, &newent); err != nil {
		return nil, common.ErrUnknown
	}
	return resp, err
}

func (h HttpServer) UpdatePet(ctx context.Context, request UpdatePetRequestObject) (UpdatePetResponseObject, error) {
	return nil, unimplementedMethod("UpdatePet")
}

func (h HttpServer) FindPetsByStatus(ctx context.Context, request FindPetsByStatusRequestObject) (FindPetsByStatusResponseObject, error) {
	return nil, unimplementedMethod("FindPetsByStatus")
}

func (h HttpServer) FindPetsByTags(ctx context.Context, request FindPetsByTagsRequestObject) (FindPetsByTagsResponseObject, error) {
	return nil, unimplementedMethod("FindPetsByTags")
}

func (h HttpServer) DeletePet(ctx context.Context, request DeletePetRequestObject) (DeletePetResponseObject, error) {
	return nil, unimplementedMethod("DeletePet")
}

func (h HttpServer) GetPetById(ctx context.Context, request GetPetByIdRequestObject) (GetPetByIdResponseObject, error) {
	var found GetPetById200JSONResponse
	ent, err := h.svc.GetPetById(ctx, request.PetId)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return GetPetById404Response{}, err
		}
		return nil, fmt.Errorf("can't get pet: %w", err)
	}
	if err := copier.Copy(&found, &ent); err != nil {
		return nil, common.ErrUnknown
	}
	return found, err
}

func (h HttpServer) UpdatePetWithForm(ctx context.Context, request UpdatePetWithFormRequestObject) (UpdatePetWithFormResponseObject, error) {
	return nil, unimplementedMethod("UpdatePetWithForm")
}

func (h HttpServer) UploadFile(ctx context.Context, request UploadFileRequestObject) (UploadFileResponseObject, error) {
	return nil, unimplementedMethod("UploadFile")
}
