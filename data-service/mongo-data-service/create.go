package mongo_data_service

import (
	"context"
	"time"

	"github.com/devingen/data-api/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Sets the data value types based on the collection configuration
// TODO validate the input date based on field configurations
func (service MongoDataService) processData(ctx context.Context, base, collection string, isCreating bool, data map[string]interface{}) error {

	tableConfig, err := service.GetTableConfig(ctx, base, collection)
	if err != nil {
		return err
	}

	if tableConfig != nil {
		if isCreating {
			// Generate an ID upfront to make the ID a string value in DB instead of ObjectId
			data["_id"] = primitive.NewObjectID().Hex()
		}

		if tableConfig.Fields != nil {
			for _, field := range tableConfig.Fields {
				fieldId := field.GetID()
				_, containsValue := data[fieldId]
				if !containsValue {
					// Skip the field if the input doesn't contain value for it
					continue
				}

				switch field.GetType() {
				case "date":
					// Convert string to ISO Date
					stringVal, ok := data[fieldId].(string)
					if ok {
						dt, err := time.Parse(time.RFC3339, stringVal)
						if err != nil {
							return err
						}
						data[fieldId] = dt
					}
				case "reference":
					// Tidy up the reference field data to prevent saving unnecessary fields and save only the required fields
					refValue, ok := data[fieldId].(map[string]interface{})
					if ok {
						id := refValue["_id"].(string)
						data[fieldId] = map[string]string{
							"_id":  id,
							"_ref": field.GetString("otherCollection"),
						}
					}
					arrayVal, ok := data[fieldId].([]interface{})
					if ok {
						for i, v := range arrayVal {
							ref := v.(map[string]interface{})
							id := ref["_id"].(string)
							data[fieldId].([]interface{})[i] = map[string]string{
								"_id":  id,
								"_ref": field.GetString("otherCollection"),
							}
						}
					}
				}
			}
		}
	}
	return nil
}

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

	err := service.processData(ctx, base, collection, true, data)
	if err != nil {
		return "", 0, err
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
