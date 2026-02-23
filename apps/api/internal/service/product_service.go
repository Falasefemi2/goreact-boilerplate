package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/falasefemi2/goreact-boilerplate/internal/db"
	"github.com/google/uuid"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrForbidden       = errors.New("forbidden")
)

type ProductService struct {
	queries db.Querier
}

func NewProductService(queries db.Querier) *ProductService {
	return &ProductService{queries: queries}
}

type CreateProductInput struct {
	Name        string
	Description string
	Price       string
	Stock       int32
}

type UpdateProductInput struct {
	Name        string
	Description string
	Price       string
	Stock       int32
}

func (s *ProductService) Create(ctx context.Context, userID string, input CreateProductInput) (db.Product, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return db.Product{}, ErrForbidden
	}

	return s.queries.CreateProduct(ctx, db.CreateProductParams{
		UserID: uid,
		Name:   input.Name,
		Description: sql.NullString{
			String: input.Description,
			Valid:  input.Description != "",
		},
		Price: input.Price,
		Stock: input.Stock,
	})
}

func (s *ProductService) GetByID(ctx context.Context, userID, productID string) (db.Product, error) {
	uid, _ := uuid.Parse(userID)
	pid, err := uuid.Parse(productID)
	if err != nil {
		return db.Product{}, ErrProductNotFound
	}

	product, err := s.queries.GetProductByID(ctx, db.GetProductByIDParams{
		ID:     pid,
		UserID: uid,
	})
	if err != nil {
		return db.Product{}, ErrProductNotFound
	}

	return product, nil
}

func (s *ProductService) List(ctx context.Context, userID string) ([]db.Product, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrForbidden
	}

	return s.queries.ListProductsByUser(ctx, uid)
}

func (s *ProductService) Update(ctx context.Context, userID, productID string, input UpdateProductInput) (db.Product, error) {
	uid, _ := uuid.Parse(userID)
	pid, err := uuid.Parse(productID)
	if err != nil {
		return db.Product{}, ErrProductNotFound
	}

	product, err := s.queries.UpdateProduct(ctx, db.UpdateProductParams{
		ID:     pid,
		UserID: uid,
		Name:   input.Name,
		Description: sql.NullString{
			String: input.Description,
			Valid:  input.Description != "",
		},
		Price: input.Price,
		Stock: input.Stock,
	})
	if err != nil {
		return db.Product{}, ErrProductNotFound
	}

	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, userID, productID string) error {
	uid, _ := uuid.Parse(userID)
	pid, err := uuid.Parse(productID)
	if err != nil {
		return ErrProductNotFound
	}

	return s.queries.DeleteProduct(ctx, db.DeleteProductParams{
		ID:     pid,
		UserID: uid,
	})
}
