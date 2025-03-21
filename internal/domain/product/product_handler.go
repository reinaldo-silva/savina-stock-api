package product

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/reinaldo-silva/savina-stock/internal/domain/image_service"
	"github.com/reinaldo-silva/savina-stock/internal/domain/product_image"
	"github.com/reinaldo-silva/savina-stock/package/response/error"
	"github.com/reinaldo-silva/savina-stock/package/response/response"
	"github.com/reinaldo-silva/savina-stock/utils"
)

type ProductHandler struct {
	useCase      *ProductUseCase
	imageService *image_service.ImageService
}

func NewProductHandler(uc *ProductUseCase, cs *image_service.ImageService) *ProductHandler {
	return &ProductHandler{uc, cs}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	nameFilter := r.URL.Query().Get("name")
	categoryIDsStr := r.URL.Query()["category_ids"]

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	var categoryIDs []uint
	for _, idStr := range categoryIDsStr {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err == nil {
			categoryIDs = append(categoryIDs, uint(id))
		}
	}

	products, total, err := h.useCase.GetAll(r.Context(), page, pageSize, nameFilter, categoryIDs, r.Host)
	if err != nil {
		appError := error.NewAppError(err.Error(), http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(products, "Products fetched successfully", &total)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *ProductHandler) GetProductsToAdmin(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	nameFilter := r.URL.Query().Get("name")
	categoryIDsStr := r.URL.Query()["category_ids"]

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	var categoryIDs []uint
	for _, idStr := range categoryIDsStr {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err == nil {
			categoryIDs = append(categoryIDs, uint(id))
		}
	}

	products, total, err := h.useCase.GetAllToAdmin(r.Context(), page, pageSize, nameFilter, categoryIDs, r.Host)
	if err != nil {
		appError := error.NewAppError(err.Error(), http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(products, "Products fetched successfully", &total)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct Product

	err := json.NewDecoder(r.Body).Decode(&newProduct)
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

	appResponse := response.NewAppResponse(createdProduct, "Product created successfully", nil, http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appResponse.StatusCode)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *ProductHandler) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	product, err := h.useCase.GetBySlug(slug)
	if err != nil {
		appError := error.NewAppError("Product not found", http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	host := r.Host

	for i := range product.Images {
		product.Images[i].ImageURL = utils.GenerateImageURL(host, product.Images[i].PublicID)
	}

	appResponse := response.NewAppResponse(product, "Product fetched successfully", nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appResponse.StatusCode)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *ProductHandler) GetProductBySlugToAdmin(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	product, err := h.useCase.GetBySlugToAdmin(slug)
	if err != nil {
		appError := error.NewAppError("Product not found", http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	host := r.Host

	for i := range product.Images {
		product.Images[i].ImageURL = utils.GenerateImageURL(host, product.Images[i].PublicID)
	}

	appResponse := response.NewAppResponse(product, "Product fetched successfully", nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appResponse.StatusCode)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	slug := chi.URLParam(r, "slug")

	err := h.useCase.Delete(slug)
	if err != nil {

		appError := error.NewAppError("Product not found", http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(nil, "Product deleted successfully", nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appResponse.StatusCode)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	slug := chi.URLParam(r, "slug")

	var updatedProduct Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		appError := error.NewAppError("Invalid input data", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	if updatedProduct.Name == "" || updatedProduct.Price <= 0 {
		appError := error.NewAppError("Invalid product data", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	updatedProductRes, err := h.useCase.Update(slug, updatedProduct)
	if err != nil {
		appError := error.NewAppError(err.Error(), http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(updatedProductRes, "Product updated successfully", nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appResponse.StatusCode)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *ProductHandler) UploadImages(w http.ResponseWriter, r *http.Request) {

	slug := chi.URLParam(r, "slug")

	r.ParseMultipartForm(10 << 20)

	files := r.MultipartForm.File["images"]

	if len(files) > 5 {
		appError := error.NewAppError("You can only upload up to 5 images", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	var uploadedImages []product_image.UploadedImage
	for _, fileHeader := range files {

		file, err := fileHeader.Open()
		if err != nil {
			appError := error.NewAppError("Failed to open the image", http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(appError.StatusCode)
			json.NewEncoder(w).Encode(appError)
			return
		}
		defer file.Close()

		tempFile, err := os.CreateTemp("", "upload-*.png")
		if err != nil {
			appError := error.NewAppError("Failed to create temp file", http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(appError.StatusCode)
			json.NewEncoder(w).Encode(appError)
			return
		}
		defer tempFile.Close()

		_, err = io.Copy(tempFile, file)
		if err != nil {
			appError := error.NewAppError("Failed to save the image", http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(appError.StatusCode)
			json.NewEncoder(w).Encode(appError)
			return
		}

		host := r.Host

		publicID, err := h.imageService.Upload(tempFile.Name())
		uploadedURL := utils.GenerateImageURL(host, publicID)

		if err != nil {
			appError := error.NewAppError("Failed to upload the image", http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(appError.StatusCode)
			json.NewEncoder(w).Encode(appError)
			return
		}

		uploadedImages = append(uploadedImages, product_image.UploadedImage{URL: uploadedURL, PublicID: publicID})

		os.Remove(tempFile.Name())
	}

	err := h.useCase.AddImagesToProduct(slug, uploadedImages)
	if err != nil {
		appError := error.NewAppError(err.Error(), http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(uploadedImages, "Images uploaded successfully", nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appResponse.StatusCode)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *ProductHandler) GetProductImages(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	product, err := h.useCase.GetBySlug(slug)
	if err != nil {
		appError := error.NewAppError("Product not found", http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	images, err := h.useCase.GetProductImages(product.ID)
	if err != nil {
		appError := error.NewAppError("Failed to fetch product images", http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	host := r.Host

	for i := range images {
		images[i].ImageURL = utils.GenerateImageURL(host, images[i].PublicID)
	}

	appResponse := response.NewAppResponse(images, "Images fetched successfully", nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appResponse.StatusCode)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *ProductHandler) LinkCategories(w http.ResponseWriter, r *http.Request) {
	categories := r.URL.Query().Get("ids")
	slug := chi.URLParam(r, "slug")

	if categories == "" {
		err := h.useCase.UpdateProductCategories(slug, []int{})
		if err != nil {
			appError := error.NewAppError(err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(appError)
			return
		}

		appResponse := response.NewAppResponse(nil, "All product categories removed successfully", nil, http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appResponse.StatusCode)
		json.NewEncoder(w).Encode(appResponse)
		return
	}

	stringArray := strings.Split(categories, ",")
	var intArray []int
	uniqueIDs := make(map[int]bool)

	for _, s := range stringArray {
		num, err := strconv.Atoi(s)
		if err != nil {
			appError := error.NewAppError("Invalid category ID format")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(appError)
			return
		}

		if !uniqueIDs[num] {
			intArray = append(intArray, num)
			uniqueIDs[num] = true
		}
	}

	err := h.useCase.UpdateProductCategories(slug, intArray)
	if err != nil {
		appError := error.NewAppError(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(nil, "Product categories linked successfully", nil, http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appResponse.StatusCode)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *ProductHandler) SwitchAvailable(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	err := h.useCase.SwitchAvailable(slug)
	if err != nil {
		appError := error.NewAppError(err.Error(), http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(nil, "Visibilidade alterada com sucesso", nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *ProductHandler) ProductStockEntry(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var body struct {
		Quantity int `json:"quantity"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		appError := error.NewAppError("Invalid request payload", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	err = h.useCase.ProductStockEntry(slug, body.Quantity)
	if err != nil {
		appError := error.NewAppError(err.Error(), http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(nil, "Quantidade do produto atualizada com sucesso", nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *ProductHandler) ProductStockOut(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var body struct {
		Quantity int `json:"quantity"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		appError := error.NewAppError("Invalid request payload", http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	err = h.useCase.ProductStockOut(slug, body.Quantity)
	if err != nil {
		appError := error.NewAppError(err.Error(), http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(nil, "Saída de estoque do produto registrada com sucesso", nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appResponse)
}
