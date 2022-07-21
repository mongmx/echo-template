package user

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/gopherlibs/gravatar/gravatar"
	"github.com/labstack/echo/v4"
	"github.com/mongmx/echo-template/middleware"
	"github.com/mongmx/echo-template/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	accessTokenExp  = time.Hour
	refreshTokenExp = time.Hour * 24 * 7
)

type Handler struct {
	store *Store
	rd *redis.Client
}

func NewHandler(store *Store, db *gorm.DB, rd *redis.Client) *Handler {
	return &Handler{store: store, rd: rd}
}

func (h *Handler) CurrentUser(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusOK, "Could not find current user")
	}
	user, err := h.store.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusOK, "Could not find current user")
	}
	return c.JSON(http.StatusOK, map[string]*model.User{"currentUser": user})
}

func (h *Handler) GetHistory(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusOK, "Could not find current user")
	}
	histories, err := h.store.GetHistory(userID)
	if err != nil {
		return c.JSON(http.StatusOK, "Could not find histories")
	}
	return c.JSON(http.StatusOK, histories)
}

func (h *Handler) RefreshToken(c echo.Context) error {
	var reqBody struct {
		RefreshToken string `json:"refreshToken"`
	}
	err := c.Bind(&reqBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if reqBody.RefreshToken == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("input cannot empty"))
	}

	oldToken, err := jwt.Parse(reqBody.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method: HS256")
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	}

	oldClaims, ok := oldToken.Claims.(jwt.MapClaims)
	if !ok || !oldToken.Valid {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh claim")
	}

	uid, err := h.getRefreshTokenFromRedis(c, oldClaims["sub"].(string))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Could not find current user")
	}

	user, err := h.store.GetUserByID(uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Could not find current user")
	}
	accessTokenUUID, err := uuid.NewV4()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	claims := &middleware.JwtCustomClaims{
		ID:    user.ID,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenExp).Unix(),
			Issuer:    "thruster-engine",
			IssuedAt:  time.Now().Unix(),
			Subject:   accessTokenUUID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshTokenUUID, err := uuid.NewV4()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenExp).Unix(),
			Issuer:    "thruster-engine",
			IssuedAt:  time.Now().Unix(),
			Subject:   refreshTokenUUID.String(),
		},
	)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	auth := model.Auth{
		AuthState: model.AuthState{
			Name: user.Name,
			UID:  user.ID,
		},
		AccessToken:           tokenString,
		ExpiresIn:             int64(accessTokenExp.Minutes()),
		AccessTokenUUID:       accessTokenUUID.String(),
		RefreshToken:          refreshTokenString,
		RefreshTokenExpiresIn: int64(accessTokenExp.Minutes()),
		RefreshTokenUUID:      refreshTokenUUID.String(),
	}

	err = h.storeJWTAuthToRedis(c, user.ID, &auth)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, auth)
}

func (h *Handler) getRefreshTokenFromRedis(c echo.Context, refreshTokenUUID string) (string, error) {
	uid, err := h.rd.Get(c.Request().Context(), refreshTokenUUID).Result()
	if err != nil {
		return "", err
	}
	return uid, nil
}

func (h *Handler) SignIn(c echo.Context) error {
	var reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := c.Bind(&reqBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if reqBody.Email == "" || reqBody.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("input cannot empty"))
	}

	user, err := h.store.GetUserByEmail(reqBody.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Could not find current user")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Could not find current user")
	}
	err = h.store.CreateHistory(&model.History{UserID: user.ID ,Email: user.Email})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	accessTokenUUID, err := uuid.NewV4()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	claims := &middleware.JwtCustomClaims{
		ID:    user.ID,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenExp).Unix(),
			Issuer:    "thruster-engine",
			IssuedAt:  time.Now().Unix(),
			Subject:   accessTokenUUID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshTokenUUID, err := uuid.NewV4()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenExp).Unix(),
			Issuer:    "thruster-engine",
			IssuedAt:  time.Now().Unix(),
			Subject:   refreshTokenUUID.String(),
		},
	)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	auth := model.Auth{
		AuthState: model.AuthState{
			Name: user.Name,
			UID:  user.ID,
		},
		AccessToken:           tokenString,
		ExpiresIn:             int64(accessTokenExp.Minutes()),
		AccessTokenUUID:       accessTokenUUID.String(),
		RefreshToken:          refreshTokenString,
		RefreshTokenExpiresIn: int64(refreshTokenExp.Minutes()),
		RefreshTokenUUID:      refreshTokenUUID.String(),
	}

	err = h.storeJWTAuthToRedis(c, user.ID, &auth)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, auth)
}

func (h *Handler) storeJWTAuthToRedis(c echo.Context, uid string, a *model.Auth) error {
	atExp := time.Duration(a.ExpiresIn) * time.Minute
	rtExp := time.Duration(a.RefreshTokenExpiresIn) * time.Minute
	err := h.rd.Set(c.Request().Context(), a.AccessTokenUUID, uid, atExp).Err()
	if err != nil {
		return err
	}
	err = h.rd.Set(c.Request().Context(), a.RefreshTokenUUID, uid, rtExp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) SignOut(c echo.Context) error {
	claims, ok := c.Get("user").(*jwt.Token).Claims.(*middleware.JwtCustomClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid token")
	}
	err := h.rd.Del(c.Request().Context(), claims.Subject).Err()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	c.Set("user_id", nil)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Sign out successfully",
	})
}

func (h *Handler) SignUp(c echo.Context) error {
	var reqBody struct {
		Name     string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := c.Bind(&reqBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if reqBody.Email == "" || reqBody.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "input cannot empty")
	}

	user := model.User{
		Name:  reqBody.Name,
		Email: reqBody.Email,
	}
	var count int64
	h.store.db.Model(&model.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, "อีเมล์นี้ถูกใช้แล้ว กรุณาใช้อีเมล์อื่น")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	user.Password = string(hashedPassword)
	img, err := gravatar.NewImage(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Could not generate avatar")
	}
	err = img.SetSize(200)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Could not generate avatar")
	}
	imgURL, err := img.URL()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Could not generate avatar")
	}
	user.Avatar = imgURL.String()
	err = h.store.CreateUser(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// User created publish
	// err = h.ec.Publish("user:created", user)
	// if err != nil {
	// 	return c.String(http.StatusInternalServerError, err.Error())
	// }
	//Email created publish
	// msg := map[string]string{
	// 	"Email":   user.Email,
	// 	"Subject": "Thank you for registering an account!",
	// 	"Text":    "Hello " + user.Name + ". Thank you for registering an account with aquabidthai.com!",
	// }
	// err = h.ec.Publish("email:created", msg)
	// if err != nil {
	// 	return c.String(http.StatusInternalServerError, err.Error())
	// }

	accessTokenUUID, err := uuid.NewV4()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	claims := &middleware.JwtCustomClaims{
		ID:    user.ID,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenExp).Unix(),
			Issuer:    "thruster-engine",
			IssuedAt:  time.Now().Unix(),
			Subject:   accessTokenUUID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshTokenUUID, err := uuid.NewV4()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenExp).Unix(),
			Issuer:    "thruster-engine",
			IssuedAt:  time.Now().Unix(),
			Subject:   refreshTokenUUID.String(),
		},
	)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	auth := model.Auth{
		AuthState: model.AuthState{
			Name: user.Name,
			UID:  user.ID,
		},
		AccessToken:           tokenString,
		ExpiresIn:             int64(accessTokenExp.Minutes()),
		AccessTokenUUID:       accessTokenUUID.String(),
		RefreshToken:          refreshTokenString,
		RefreshTokenExpiresIn: int64(accessTokenExp.Minutes()),
		RefreshTokenUUID:      refreshTokenUUID.String(),
	}

	err = h.storeJWTAuthToRedis(c, user.ID, &auth)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, auth)
}

func (h *Handler) UpdatePassword(c echo.Context) error {
	var reqBody struct {
		OldPassword     string `json:"oldPassword"`
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	err := c.Bind(&reqBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if reqBody.OldPassword == "" || reqBody.NewPassword == "" || reqBody.ConfirmPassword == "" {
		return c.JSON(http.StatusBadRequest, errors.New("input cannot empty"))
	}

	if reqBody.NewPassword != reqBody.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, errors.New("passwords do not match"))
	}
	currentUser, ok := c.Get("currentUser").(*middleware.JwtCustomClaims)
	if !ok {
		return c.JSON(http.StatusOK, "Could not find current user")
	}
	user, err := h.store.GetUserByID(currentUser.ID)
	if err != nil {
		return c.JSON(http.StatusOK, "Could not find current user")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.OldPassword))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Old password is incorrect")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user.Password = string(hashedPassword)
	err = h.store.UpdateUser(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Could not save user")
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) GetProfile(c echo.Context) error {
	username := c.Param("username")
	u, err := h.store.GetProfile(username)
	userID, ok := c.Get("user_id").(string)
	if !ok {
		userID = ""
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Could not find user")
	}
	return c.JSON(http.StatusOK, newProfileResponse(h.store, userID, u))
}

func (h *Handler) FollowProfile(c echo.Context) error {
	username := c.Param("username")
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not find current user")
	}
	u, err := h.store.GetProfile(username)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Could not find user")
	}
	err = h.store.AddFollower(u, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not follow user")
	}
	return c.JSON(http.StatusOK, newProfileResponse(h.store, userID, u))
}

func (h *Handler) UnfollowProfile(c echo.Context) error {
	username := c.Param("username")
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not find current user")
	}
	u, err := h.store.GetProfile(username)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Could not find user")
	}
	err = h.store.RemoveFollower(u, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not unfollow user")
	}
	return c.JSON(http.StatusOK, newProfileResponse(h.store, userID, u))
}
