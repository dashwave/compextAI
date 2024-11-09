package handlers

import "errors"

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (r *CreateUserRequest) Validate() error {
	if r.Username == "" {
		return errors.New("username is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
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
