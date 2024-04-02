package mongo_data_service

import (
	"context"
	"github.com/devingen/data-api/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (service MongoDataService) Create(ctx context.Context, base, collection string, config *dto.CreateConfig) (string, int, error) {

	data := map[string]interface{}{
		"_created":  time.Now(),
		"_updated":  time.Now(),
		"_revision": 1,
	}
	if item, ok := config.Data.(map[string]interface{}); ok {
		for k, v := range item {
			data[k] = v
		}
	}

	coll := service.Database.Client.Database(base).Collection(collection)
	result, err := coll.InsertOne(
		context.Background(),
		data,
	)

	id, ok := result.InsertedID.(string)
	if !ok {
		id = (result.InsertedID.(primitive.ObjectID)).Hex()
	}
	return id, 0, err
}
