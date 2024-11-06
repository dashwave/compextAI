package handlers

import "errors"

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *CreateUserRequest) Validate() error {
	if r.Username == "" || r.Password == "" {
		return errors.New("username and password are required")
	}
	if len(r.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

type CreateUserResponse struct {
	APIToken string `json:"api_token"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Username == "" || r.Password == "" {
		return errors.New("username and password are required")
	}
	return nil
}

type LoginResponse struct {
	APIToken string `json:"api_token"`
}
