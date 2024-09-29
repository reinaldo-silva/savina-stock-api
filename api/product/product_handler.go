package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/reinaldo-silva/savina-stock/internal/domain/product"
	usecase "github.com/reinaldo-silva/savina-stock/internal/usecase/product"
	"github.com/reinaldo-silva/savina-stock/package/response/error"
	"github.com/reinaldo-silva/savina-stock/package/response/response"
)

type ProductHandler struct {
	useCase *usecase.ProductUseCase
}

func NewProductHandler(uc *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{uc}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.useCase.GetAll()
	if err != nil {
		appError := error.NewAppError(err.Error(), http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(products, "Products fetched successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appResponse.StatusCode)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct product.Product

	err := json.NewDecoder(r.Body).Decode(&newProduct)

	fmt.Println(newProduct)

	if err != nil {
		appError := error.NewAppError("Invalid input data")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	if newProduct.Name == "" || newProduct.Price <= 0 {
		appError := error.NewAppError("Invalid product data")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	createdProduct, err := h.useCase.Create(newProduct)
	if err != nil {
		appError := error.NewAppError(err.Error(), http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(createdProduct, "Product created successfully", http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appResponse.StatusCode)
	json.NewEncoder(w).Encode(appResponse)
}
