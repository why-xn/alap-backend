package model

import (
	"github.com/why-xn/alap-backend/pkg/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ChatWindow struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Users         []string           `bson:"users" json:"users"`
	Status        enum.Status        `json:"status" bson:"status"`
	CreatedAt     time.Time          `json:"createDate" bson:"createdAt"`
	LastUpdatedAt time.Time          `json:"lastUpdatedAt" bson:"lastUpdatedAt"`
}
