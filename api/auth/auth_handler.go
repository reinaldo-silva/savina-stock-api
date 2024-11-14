package api_auth

import (
	"encoding/json"
	"net/http"

	"github.com/reinaldo-silva/savina-stock/internal/domain/user"
	usecase "github.com/reinaldo-silva/savina-stock/internal/usecase/user"
	error_response "github.com/reinaldo-silva/savina-stock/package/response/error"
	"github.com/reinaldo-silva/savina-stock/package/response/response"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	useCase *usecase.UserUseCase
}

func NewAuthHandler(uc *usecase.UserUseCase) *AuthHandler {
	return &AuthHandler{useCase: uc}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var newUser user.User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		appError := error_response.NewAppError("Invalid input data", http.StatusBadRequest)
		h.sendErrorResponse(w, appError)
		return
	}

	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" {
		appError := error_response.NewAppError("Name, email, and password are required", http.StatusBadRequest)
		h.sendErrorResponse(w, appError)
		return
	}

	userFound, _ := h.useCase.GetByEmail(newUser.Email)

	if userFound != nil {
		appError := error_response.NewAppError("Email ja está em uso!", http.StatusBadRequest)
		h.sendErrorResponse(w, appError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		appError := error_response.NewAppError("Error generating password hash", http.StatusInternalServerError)
		h.sendErrorResponse(w, appError)
		return
	}
	newUser.Password = string(hashedPassword)

	createdUser, err := h.useCase.Create(newUser)
	if err != nil {
		appError := error_response.NewAppError(err.Error(), http.StatusInternalServerError)
		h.sendErrorResponse(w, appError)
		return
	}

	appResponse := response.NewAppResponse(createdUser, "User created successfully", nil)
	h.sendSuccessResponse(w, appResponse)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		appError := error_response.NewAppError("Invalid input data", http.StatusBadRequest)
		h.sendErrorResponse(w, appError)
		return
	}

	if loginData.Email == "" || loginData.Password == "" {
		appError := error_response.NewAppError("Email and password are required", http.StatusBadRequest)
		h.sendErrorResponse(w, appError)
		return
	}

	user, token, err := h.useCase.SignInUseCase(loginData.Email, loginData.Password)
	if err != nil {
		appError := error_response.NewAppError(err.Error(), http.StatusUnauthorized)
		h.sendErrorResponse(w, appError)
		return
	}

	appResponse := response.NewAppResponse(map[string]interface{}{
		"user":  user,
		"token": token,
	}, "User signed in successfully", nil)
	h.sendSuccessResponse(w, appResponse)
}

func (h *AuthHandler) sendErrorResponse(w http.ResponseWriter, appError error_response.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appError.StatusCode)
	json.NewEncoder(w).Encode(appError)
}

func (h *AuthHandler) sendSuccessResponse(w http.ResponseWriter, appResponse response.AppResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appResponse)
}
