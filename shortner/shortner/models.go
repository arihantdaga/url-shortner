package shortner

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShortUrl struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	Shortkey  string             `bson:"shortkey" json:"shortkey"`
	Url       string             `bson:"original_url" json:"original_url"`
	CreatedAt time.Time          `bson:"CreatedAt" json:"CreatedAt"`
}
