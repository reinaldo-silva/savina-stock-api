package api_image

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	usecase_image "github.com/reinaldo-silva/savina-stock/internal/usecase/image"
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
