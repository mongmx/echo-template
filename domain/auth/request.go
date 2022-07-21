package auth

import (
	"github.com/labstack/echo/v4"
)

type registerRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

func (r *registerRequest) bind(c echo.Context, u *User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Email = r.Email
	u.Password = r.Password
	return nil
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r *loginRequest) bind(c echo.Context, u *User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Email = r.Email
	u.Password = r.Password
	return nil
}
