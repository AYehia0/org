package mongodb

import "github.com/ayehia0/org/pkg/database/mongodb/repository"

// here we gather all the repositories to allow access different collections
type Store struct {
	UserRepository         repository.UserRepository
	SessionRepository      repository.SessionRepository
	OrganizationRepository repository.OrganizationRepository
}

func NewStore(conn *MongoDBConn) *Store {

	// create the collections for the different models
	userCol := conn.Database.Collection("users")
	sessionCol := conn.Database.Collection("sessions")
	orgCol := conn.Database.Collection("organizations")

	// create the user repository
	user := repository.NewUserRepository(userCol)
	session := repository.NewSessionRepository(sessionCol)
	org := repository.NewOrganizationRepository(orgCol)

	return &Store{
		UserRepository:         user,
		SessionRepository:      session,
		OrganizationRepository: org,
	}
}
