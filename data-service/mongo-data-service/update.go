package mongo_data_service

import (
	"context"
	"strconv"
	"time"

	"github.com/devingen/data-api/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (service MongoDataService) Update(ctx context.Context, base, collection, id string, config *dto.UpdateConfig) (string, int, error) {

	data := map[string]interface{}{
		"_updated": time.Now(),
	}
	if item, ok := config.Data.(map[string]interface{}); ok {
		for k, v := range item {
			data[k] = v
		}
	}

	err := service.processData(ctx, base, collection, false, data)
	if err != nil {
		return "", 0, err
	}

	update := bson.D{
		{"$inc", bson.M{"_revision": 1}},
		{"$set", data},
	}

	coll := service.Database.Client.Database(base).Collection(collection)
	result, err := coll.UpdateOne(
		ctx,
		bson.M{"_id": id},
		update,
	)

	if result.ModifiedCount == 0 {
		// TODO decide this with env var, header or some other parameter
		// try object id if nothing matches
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return "", 0, err
		}

		result, err = coll.UpdateOne(
			ctx,
			bson.M{"_id": oid},
			update,
		)
	}

	return strconv.Itoa(int(result.ModifiedCount)), 0, err
}
