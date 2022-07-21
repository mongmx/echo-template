package auth

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Register(c echo.Context) error {
	merchantID := c.Request().Header.Get("X-Merchant-Id")
	log.Printf("merchantID: %s", merchantID)

	u := User{}
	r := registerRequest{}
	if err := r.bind(c, &u); err != nil {
		return err
	}

	err := h.repo.CreateUser(&u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, ErrCreateUser)
	}

	token, err := u.GenAuthToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, ErrJWTInvalid)
	}

	err = h.repo.StoreJWT(u.ID, token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = h.repo.CreateLog(&History{UserID: u.ID, Name: u.Name})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, newAuthResponse(token))

	// return c.JSON(http.StatusOK, "")
}

func (h *Handler) Login(c echo.Context) error {
	u := User{}
	req := loginRequest{}
	if err := req.bind(c, &u); err != nil {
		return err
	}

	queryUser, err := h.repo.GetUserByEmail(u.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = u.ComparePassword(queryUser.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password")
	}

	token, err := u.GenAuthToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = h.repo.StoreJWT(u.ID, token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = h.repo.CreateLog(&History{UserID: u.ID, Name: u.Name})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, newAuthResponse(token))
}

func (h *Handler) GetProfile(c echo.Context) error {
	user, err := h.repo.GetUserByID(c.Get("user_id").(string))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, newProfileResponse(user))
}

// func (h *Handler) SelectMerchant(c echo.Context) error {
// 	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
// 	if err != nil {
// 		return err
// 	}
// 	user, err := h.repo.GetMerchantByID(id)
// 	if err != nil {
// 		return err
// 	}
// 	return c.JSON(http.StatusOK, newMerchantResponse(user))
// }
