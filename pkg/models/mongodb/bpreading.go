package mongodb

import (
	"context"
	"errors"

	"github.com/anupx73/go-bpcalc-backend-k8s/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BPReadingModel represent a mgo database session with a bpReading data model.
type BPReadingModel struct {
	C *mongo.Collection
}

// All method will be used to get all records from the bpReading table.
func (m *BPReadingModel) All() ([]models.BPReading, error) {
	// Define variables
	ctx := context.TODO()
	mm := []models.BPReading{}

	// Find all bpReading
	bpReadingsCursor, err := m.C.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = bpReadingsCursor.All(ctx, &mm)
	if err != nil {
		return nil, err
	}

	return mm, err
}

// FindByID will be used to find a new movie registry by id
func (m *BPReadingModel) FindByID(id string) (*models.BPReading, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Find bpReading by id
	var bpReading = models.BPReading{}
	err = m.C.FindOne(context.TODO(), bson.M{"_id": p}).Decode(&bpReading)
	if err != nil {
		// Checks if the bpReading was not found
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ErrNoDocuments")
		}
		return nil, err
	}

	return &bpReading, nil
}

// Insert will be used to insert a new bpReading registry
func (m *BPReadingModel) Insert(bpReading models.BPReading) (*mongo.InsertOneResult, error) {
	return m.C.InsertOne(context.TODO(), bpReading)
}

// Delete will be used to delete a bpReading registry
func (m *BPReadingModel) Delete(id string) (*mongo.DeleteResult, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return m.C.DeleteOne(context.TODO(), bson.M{"_id": p})
}
