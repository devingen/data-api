package mongo_data_service

import (
	"context"
	"github.com/devingen/data-api/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
)

func (service MongoDataService) Update(ctx context.Context, base, collection, id string, config *dto.UpdateConfig) (string, int, error) {

	coll := service.Database.Client.Database(base).Collection(collection)
	result, err := coll.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.D{
			{"$set", config.Data},
		},
	)

	if result.ModifiedCount == 0 {
		// TODO decide this with env var, header or some other parameter
		// try object id if nothing matches
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return "", 0, err
		}

		result, err = coll.UpdateOne(
			context.Background(),
			bson.M{"_id": oid},
			bson.D{
				{"$set", config.Data},
			},
		)
	}

	return strconv.Itoa(int(result.ModifiedCount)), 0, err
}
