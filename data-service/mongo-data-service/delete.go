package mongo_data_service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
)

func (service MongoDataService) Delete(ctx context.Context, base, collection, id string) (string, int, error) {

	coll := service.Database.Client.Database(base).Collection(collection)
	result, err := coll.DeleteOne(
		context.Background(),
		bson.M{"_id": id},
	)

	if result.DeletedCount == 0 {
		// TODO decide this with env var, header or some other parameter
		// try object id if nothing matches
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return "", 0, err
		}

		result, err = coll.DeleteOne(
			context.Background(),
			bson.M{"_id": oid},
		)
	}

	return strconv.Itoa(int(result.DeletedCount)), 0, err
}
