package repository

import (
	"context"
	"errors"

	"github.com/ayehia0/org/pkg/database/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// the repository package contains the database operations for the organization model
type OrganizationRepository interface {
	Create(ctx context.Context, org *models.Organization) error                         // Create a new organization
	FindByID(ctx context.Context, id string) (*models.Organization, error)              // Find an organization by id
	Update(ctx context.Context, org *models.Organization) (*models.Organization, error) // Update an organization returns the updated organization
	Delete(ctx context.Context, id string) error                                        // Delete an organization
	AddMember(ctx context.Context, orgID string, member *models.Member) error           // Add a member to an organization
	RemoveMember(ctx context.Context, orgID string, memberID string) error              // Remove a member from an organization
}

// the organization repository struct
type organizationRepository struct {
	// the database connection
	col *mongo.Collection
}

// the function to create a new organization
func (r *organizationRepository) Create(ctx context.Context, org *models.Organization) error {
	_, err := r.col.InsertOne(ctx, org)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("organization already exists")
		}
	}
	return err
}

// the function to find an organization by id
func (r *organizationRepository) FindByID(ctx context.Context, id string) (*models.Organization, error) {
	var org models.Organization
	err := r.col.FindOne(ctx, models.Organization{ID: id}).Decode(&org)
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
	_, err := r.col.UpdateOne(ctx, models.Organization{ID: org.ID}, org)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("organization not found")
		}
	}
	return org, err
}

// the function to delete an organization
func (r *organizationRepository) Delete(ctx context.Context, id string) error {
	_, err := r.col.DeleteOne(ctx, models.Organization{ID: id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("organization not found")
		}
	}
	return err
}

// the function to add a member to an organization
func (r *organizationRepository) AddMember(ctx context.Context, orgID string, member *models.Member) error {
	_, err := r.col.UpdateOne(ctx, models.Organization{ID: orgID}, bson.M{"$push": bson.M{"members": member}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("organization not found")
		}
	}
	return err
}
