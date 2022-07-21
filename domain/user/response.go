package user

import "github.com/mongmx/echo-template/model"

type profileResponse struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func marshalProfileResponse(s *Store, userID string, user *model.User) *profileResponse {
	f, _ := s.IsFollower(user.ID, userID)
	return &profileResponse{
		Username:  user.Name,
		Bio:       user.Bio,
		Image:     user.Avatar,
		Following: f,
	}
}

type singleProfileResponse struct {
	Profile *profileResponse `json:"profile"`
}

func newProfileResponse(s *Store, userID string, user *model.User) *singleProfileResponse {
	return &singleProfileResponse{
		Profile: marshalProfileResponse(s, userID, user),
	}
}
