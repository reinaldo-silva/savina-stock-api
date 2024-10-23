package api_category

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {

	categories, err := h.useCase.GetAllCategories()
	if err != nil {
		appError := error.NewAppError("Failed to fetch categories", http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	if len(categories) == 0 {
		appError := error.NewAppError("No categories found", http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(categories, "Categories fetched successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	categoryID, err := strconv.Atoi(id)
	if err != nil {
		appError := error.NewAppError("Invalid category ID format", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	err = h.useCase.DeleteCategory(uint(categoryID))
	if err != nil {
		appError := error.NewAppError(err.Error(), http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(nil, "Category deleted successfully", http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appResponse.StatusCode)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	categoryID, err := strconv.Atoi(id)
	if err != nil {
		appError := error.NewAppError("Invalid category ID format", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	category, err := h.useCase.GetCategoryByID(uint(categoryID))
	if err != nil {
		appError := error.NewAppError("Category not found", http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(category, "Category fetched successfully", http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appResponse.StatusCode)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	categoryID, err := strconv.Atoi(id)
	if err != nil {
		appError := error.NewAppError("Invalid category ID format", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	var updatedCategory category.Category

	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
		appError := error.NewAppError("Invalid request payload", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	updatedCategory.ID = uint(categoryID)
	err = h.useCase.UpdateCategory(&updatedCategory)
	if err != nil {
		appError := error.NewAppError(err.Error(), http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(updatedCategory, "Category updated successfully", http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appResponse.StatusCode)
	json.NewEncoder(w).Encode(appResponse)
}
