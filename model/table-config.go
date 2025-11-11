package model

import coremodel "github.com/devingen/api-core/model"

type TableConfig struct {
	DatabaseTable *string           `json:"databaseTable,omitempty"  bson:"databaseTable,omitempty"`
	Fields        []coremodel.Field `json:"fields,omitempty"  bson:"fields,omitempty"`
}
