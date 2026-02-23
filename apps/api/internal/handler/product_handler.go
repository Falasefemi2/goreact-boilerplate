package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/falasefemi2/goreact-boilerplate/internal/middleware"
	"github.com/falasefemi2/goreact-boilerplate/internal/response"
	"github.com/falasefemi2/goreact-boilerplate/internal/service"
	appvalidator "github.com/falasefemi2/goreact-boilerplate/internal/validator"
	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

type createProductRequest struct {
	Name        string `json:"name"        validate:"required,min=1,max=255"`
	Description string `json:"description"`
	Price       string `json:"price"       validate:"required"`
	Stock       int32  `json:"stock"       validate:"min=0"`
}

type updateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Stock       int32  `json:"stock"`
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req createProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	product, err := h.productService.Create(r.Context(), userID, service.CreateProductInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not create product")
		return
	}

	response.JSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	products, err := h.productService.List(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not fetch products")
		return
	}

	response.JSON(w, http.StatusOK, products)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	productID := chi.URLParam(r, "id")

	product, err := h.productService.GetByID(r.Context(), userID, productID)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			response.Error(w, http.StatusNotFound, "product not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "could not fetch product")
		return
	}

	response.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	productID := chi.URLParam(r, "id")

	var req updateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	product, err := h.productService.Update(r.Context(), userID, productID, service.UpdateProductInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			response.Error(w, http.StatusNotFound, "product not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "could not update product")
		return
	}

	response.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	productID := chi.URLParam(r, "id")

	if err := h.productService.Delete(r.Context(), userID, productID); err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			response.Error(w, http.StatusNotFound, "product not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "could not delete product")
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
