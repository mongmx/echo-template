package auth

type profileResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func newProfileResponse(u *User) *profileResponse {
	return &profileResponse{
		Username: u.Name,
		Email:    u.Email,
	}
}
type authResponse struct {
	AccessToken           string `json:"token"`
	AccessTokenExpiresIn  int64  `json:"expiresIn"`
	RefreshToken          string `json:"refreshToken"`
	RefreshTokenExpiresIn int64  `json:"refreshTokenExpiresIn"`
}

func newAuthResponse(a *Auth) *authResponse {
	return &authResponse{
		AccessToken:           a.AccessToken,
		AccessTokenExpiresIn:  a.AccessTokenExpiresIn,
		RefreshToken:          a.RefreshToken,
		RefreshTokenExpiresIn: a.RefreshTokenExpiresIn,
	}
}
