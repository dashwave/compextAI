package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/burnerlee/compextAI/controllers"
	"github.com/burnerlee/compextAI/utils/responses"
)

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := controllers.CreateUser(s.DB, &controllers.CreateUserRequest{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
	})

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, CreateUserResponse{APIToken: user.APIToken})
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		responses.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := controllers.Login(s.DB, &controllers.LoginRequest{
		Username: request.Username,
		Password: request.Password,
	})

	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, LoginResponse{APIToken: user.APIToken})
}
