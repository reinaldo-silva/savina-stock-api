package api_category

import (
	"encoding/json"
	"net/http"

	"github.com/reinaldo-silva/savina-stock/internal/domain/category"
	usecase_category "github.com/reinaldo-silva/savina-stock/internal/usecase/category"
	"github.com/reinaldo-silva/savina-stock/package/response/error"
	"github.com/reinaldo-silva/savina-stock/package/response/response"
)

type CategoryHandler struct {
	useCase *usecase_category.CategoryUseCase
}

func NewCategoryHandler(uc *usecase_category.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{uc}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category category.Category

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		appError := error.NewAppError("Invalid request payload", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	createdCategory, err := h.useCase.CreateCategory(&category)
	if err != nil {
		appError := error.NewAppError(err.Error(), http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(createdCategory, "Category created successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(appResponse)
}
