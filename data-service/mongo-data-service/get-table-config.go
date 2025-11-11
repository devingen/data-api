package mongo_data_service

import (
	"context"

	"github.com/devingen/api-core/database"
	"github.com/devingen/data-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (service MongoDataService) GetTableConfig(ctx context.Context, base, collection string) (*model.TableConfig, error) {

	var result *model.TableConfig
	err := service.Database.Find(ctx, base, "harman-table-configs", bson.M{"databaseTable": collection}, database.FindOptions{}, func(cur *mongo.Cursor) error {
		err := cur.Decode(&result)
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}
