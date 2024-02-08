package models

// Session is a struct that represents the session model, the session model is used to store the user session (access_token, refresh_token, and the user id)
type Session struct {
	ID                  string `json:"id" bson:"_id,omitempty"`
	AccessToken         string `json:"access_token" bson:"access_token"`
	AccessTokenExpires  int64  `json:"access_token_expires" bson:"access_token_expires"`
	RefreshToken        string `json:"refresh_token" bson:"refresh_token"`
	RefreshTokenExpires int64  `json:"refresh_token_expires" bson:"refresh_token_expires"`
	UserID              string `json:"user_id" bson:"user_id"`
}
