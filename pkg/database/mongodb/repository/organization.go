package repository

import (
	"context"
	"errors"

	"github.com/ayehia0/org/pkg/database/mongodb/models"
	"github.com/ayehia0/org/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// the repository package contains the database operations for the organization model
type OrganizationRepository interface {
	Create(ctx context.Context, org *models.Organization) (string, error)               // Create a new organization and return the id as a string
	FindByID(ctx context.Context, id string) (*models.Organization, error)              // Find an organization by id
	Update(ctx context.Context, org *models.Organization) (*models.Organization, error) // Update an organization returns the updated organization
	Delete(ctx context.Context, id string) error                                        // Delete an organization
	AddMember(ctx context.Context, orgID string, member *models.Member) error           // Add a member to an organization
	FindAll(ctx context.Context) ([]models.Organization, error)                         // Find all organizations
	InviteUserToOrganization(ctx context.Context, orgID string, member models.Member) error
	IsUserInOrganization(ctx context.Context, orgID string, email string) (bool, error)
}

// the organization repository struct
type organizationRepository struct {
	// the database connection
	col *mongo.Collection
}

// create a new organization repository
func NewOrganizationRepository(col *mongo.Collection) OrganizationRepository {
	return &organizationRepository{col: col}
}

// the function to create a new organization
func (r *organizationRepository) Create(ctx context.Context, org *models.Organization) (string, error) {
	// create a new organization and return the created organization
	res, err := r.col.InsertOne(ctx, org)
	if err != nil {
		return "", err
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("error converting id to string")
	}
	return id.Hex(), nil
}

// the function to find an organization by id
func (r *organizationRepository) FindByID(ctx context.Context, id string) (*models.Organization, error) {
	var org models.Organization

	objectID, err := utils.StringToObjectID(id)
	if err != nil {
		return nil, err
	}
	err = r.col.FindOne(ctx, bson.M{"_id": objectID}).Decode(&org)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("organization not found")
		}
	}
	return &org, err
}

// the function to update an organization
// update only the given fields
func (r *organizationRepository) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {

	objectID, err := utils.StringToObjectID(org.ID)
	if err != nil {
		return nil, err
	}
	updateFields := bson.M{
		"$set": bson.M{
			"name": org.Name,
			"desc": org.Desc,
		},
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": objectID}, updateFields)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("organization not found")
		}
	}
	return org, err
}

// the function to delete an organization
func (r *organizationRepository) Delete(ctx context.Context, id string) error {

	objectID, err := utils.StringToObjectID(id)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("organization not found")
		}
	}
	return err
}

// the function to add a member to an organization
func (r *organizationRepository) AddMember(ctx context.Context, orgID string, member *models.Member) error {

	objectID, err := utils.StringToObjectID(orgID)
	if err != nil {
		return err
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$push": bson.M{"members": member}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("organization not found")
		}
	}
	return err
}

func (r *organizationRepository) FindAll(ctx context.Context) ([]models.Organization, error) {
	var orgs []models.Organization
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &orgs); err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *organizationRepository) InviteUserToOrganization(ctx context.Context, orgID string, member models.Member) error {
	objectID, err := utils.StringToObjectID(orgID)
	if err != nil {
		return err
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$push": bson.M{"members": member}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("organization not found")
		}
	}
	return err
}

func (r *organizationRepository) IsUserInOrganization(ctx context.Context, orgID string, email string) (bool, error) {
	orgObjectId, err := utils.StringToObjectID(orgID)
	if err != nil {
		return false, err
	}

	count, err := r.col.CountDocuments(ctx, bson.M{"_id": orgObjectId, "members.email": email})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
