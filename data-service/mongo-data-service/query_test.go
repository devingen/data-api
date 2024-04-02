package mongo_data_service

import (
	coremodel "github.com/devingen/api-core/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPickQueryType(t *testing.T) {

	assert.Equal(t, PickQueryType(&coremodel.QueryConfig{
		Limit: 1000,
		Fields: []coremodel.Field{
			coremodel.New(coremodel.FieldTypeText, "title"),
			coremodel.NewReverseReference(
				"lastGroup",
				"groups",
				"space",
				true,
			).SetSort(&coremodel.SortConfig{
				ID:    "createdAt",
				Order: -1,
			}).SetLimit(1).ToField(),
		},
	}), QueryTypeAdvanced)
}
