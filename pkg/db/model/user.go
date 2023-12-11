package model

import (
	"github.com/why-xn/alap-backend/pkg/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Username      string             `json:"username" bson:"username"`
	FullName      string             `json:"fullName" bson:"fullName"`
	Email         string             `json:"email" bson:"email"`
	Status        enum.Status        `json:"status" bson:"status"`
	CreatedAt     time.Time          `json:"createDate" bson:"createdAt"`
	LastUpdatedAt time.Time          `json:"lastUpdatedAt" bson:"lastUpdatedAt"`
}
