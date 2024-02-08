package models

// Organization is a struct that represents the organization model
// the organization have a name, description and members
type Member struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	Email       string `json:"email" bson:"email"`
	AccessLevel string `json:"accessLevel" bson:"accessLevel"`
}

type Organization struct {
	ID      string   `json:"id" bson:"_id,omitempty"`
	Name    string   `json:"name" bson:"name"`
	Desc    string   `json:"desc" bson:"desc"`
	Members []Member `json:"members" bson:"members"`
}
