package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

// convert string id to object id
func StringToObjectID(id string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return objectID, nil
}
