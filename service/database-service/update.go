package database_service

import (
	"context"
	"github.com/devingen/data-api/dto"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
)

func (service DatabaseService) Update(base, collection, id string, config *dto.UpdateConfig) (string, int, error) {

	coll := service.Database.Client.Database(base).Collection(collection)
	result, err := coll.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.D{
			{"$set", config.Data},
		},
	)

	return strconv.Itoa(int(result.ModifiedCount)), 0, err
}
