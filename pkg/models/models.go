package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BPReading is used to represent patient profile data
type BPReading struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Email       string             `bson:"email,omitempty"`
	Systolic    string             `bson:"systolic,omitempty"`
	Diastolic   string             `bson:"diastolic,omitempty"`
	Category    string             `bson:"category,omitempty"`
	ReadingTime time.Time          `bson:"readingtime,omitempty"`
}
