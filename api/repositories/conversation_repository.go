package repositories

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/brando1998/docubot-api/models"
)

type ConversationRepository interface {
	SaveMessage(ctx context.Context, userID uint, botID uint, message models.Message) error
	GetConversationByUserID(ctx context.Context, userID uint) (*models.Conversation, error)
}

type conversationRepository struct {
	collection *mongo.Collection
}

// Constructor
func NewConversationRepository(client *mongo.Client) ConversationRepository {
	collection := client.Database(os.Getenv("MONGO_DB")).Collection("conversations")
	return &conversationRepository{collection}
}

// Implementación de SaveMessage
func (r *conversationRepository) SaveMessage(ctx context.Context, userID uint, botID uint, message models.Message) error {
	filter := bson.M{"user_id": userID, "bot_id": botID}
	update := bson.M{
		"$push": bson.M{"messages": message},
		"$setOnInsert": bson.M{
			"user_id":    userID,
			"bot_id":     botID,
			"created_at": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// Implementación de GetConversationByUserID
func (r *conversationRepository) GetConversationByUserID(ctx context.Context, userID uint) (*models.Conversation, error) {
	filter := bson.M{"user_id": userID}
	var conversation models.Conversation
	err := r.collection.FindOne(ctx, filter).Decode(&conversation)
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}
