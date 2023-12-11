package model

import (
	"github.com/why-xn/alap-backend/pkg/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ChatMessage struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	ChatWindow    string             `json:"chatWindow"`
	Sender        dto.MinUserDTO     `json:"sender"`
	Msg           string             `json:"msg"`
	Status        string             `json:"status" bson:"status"`
	CreatedAt     time.Time          `json:"createDate" bson:"createdAt"`
	LastUpdatedAt time.Time          `json:"lastUpdatedAt" bson:"lastUpdatedAt"`
}
