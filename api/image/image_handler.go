package api_image

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	usecase_image "github.com/reinaldo-silva/savina-stock/internal/usecase/image"
	"github.com/reinaldo-silva/savina-stock/package/response/error"
	"github.com/reinaldo-silva/savina-stock/package/response/response"
)

type ImageHandler struct {
	useCase *usecase_image.ImageUseCase
}

func NewImageHandler(uc *usecase_image.ImageUseCase) *ImageHandler {
	return &ImageHandler{uc}
}

func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	imageBuffer, contentType, err := h.useCase.GetImage(uuid)
	if err != nil {
		http.Error(w, "Imagem não encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "inline")

	_, err = io.Copy(w, imageBuffer)
	if err != nil {
		http.Error(w, "Erro ao enviar imagem", http.StatusInternalServerError)
		return
	}
}

func (h *ImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	err := h.useCase.DeleteImage(uuid)
	if err != nil {
		appError := error.NewAppError(err.Error(), http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(nil, "Imagem deletada com sucesso", nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appResponse)
}

func (h *ImageHandler) SetImageAsCover(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	slug := chi.URLParam(r, "slug")

	err := h.useCase.SetImageAsCover(uuid, slug)
	if err != nil {
		appError := error.NewAppError("Erro ao definir a imagem como capa", http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appError.StatusCode)
		json.NewEncoder(w).Encode(appError)
		return
	}

	appResponse := response.NewAppResponse(nil, "Imagem definida como capa com sucesso", nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appResponse)
}
