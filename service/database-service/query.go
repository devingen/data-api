package database_service

import (
	coremodel "github.com/devingen/api-core/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AggregateResult struct {
	Results     []*coremodel.DataModel `bson:"results"`
	MetaWrapper []*coremodel.Meta      `bson:"meta"`
}

type QueryPipeline []bson.M

func NewPipeline() QueryPipeline {
	return make(QueryPipeline, 0)
}

func (q *QueryPipeline) AddLimit(limit int) {
	//*q = append(*q, bson.M{"$limit": limit})

	// add total number meta
	*q = append(*q, bson.M{
		"$facet": bson.M{
			"meta": []bson.M{{
				"$group": bson.M{
					"_id":   nil,
					"total": bson.M{"$sum": 1},
				},
			}},
			"results": []bson.M{{
				"$limit": limit,
			}},
		},
	})
}

func (q *QueryPipeline) AddSkip(skip int) {
	*q = append(*q, bson.M{"$skip": skip})
}

func (q *QueryPipeline) AddSort(id string, order int) {
	*q = append(*q, bson.M{"$sort": bson.M{id: order}})
}

func (q *QueryPipeline) AddProject(fields []coremodel.Field) {

	project := bson.M{
		"_id": 1,
	}

	for _, field := range fields {
		project[field.GetID()] = 1
	}

	*q = append(*q, bson.M{
		"$project": project,
	})
}

func (q *QueryPipeline) AddMatch(config *coremodel.QueryConfig) {
	*q = append(*q, bson.M{"$match": config.Filter.ToMatchQuery(config)})
}

func (q *QueryPipeline) AddReference(fields []coremodel.Field, ref coremodel.ReferenceField) {
	*q = append(*q, bson.M{
		"$lookup": bson.M{
			"from":         ref.GetOtherCollection(),
			"localField":   ref.GetID() + "._id",
			"foreignField": "_id",
			"as":           ref.GetID(),
		},
	})

	// filter
	////////////////////////////////////////////////////////////////
	if ref.GetFilter() != nil {
		project := bson.M{
			"_id": 1,
		}

		// keep parent fields
		for _, field := range fields {
			if ref.GetID() != field.GetID() {
				project[field.GetID()] = 1
			}
		}

		project[ref.GetID()] = bson.M{
			"$filter": bson.M{
				"input": "$" + ref.GetID(),
				"as":    ref.GetID(),
				"cond":  ref.GetFilter().ToFilterQuery(ref.GetID()),
			},
		}

		*q = append(*q, bson.M{
			"$project": project,
		})

		// remove the relations that doesn't contain the other collection
		*q = append(*q, bson.M{
			"$match": bson.M{ref.GetID(): bson.M{"$not": bson.M{"$size": 0}}},
		})
	}
	////////////////////////////////////////////////////////////////

	// filter fields
	////////////////////////////////////////////////////////////////
	project := bson.M{
		"_id": 1,
	}

	// keep parent fields
	for _, field := range fields {
		if ref.GetID() != field.GetID() {
			project[field.GetID()] = 1
		}
	}

	// keep reference fields
	project[ref.GetID()+"._id"] = 1
	if ref.GetFields() != nil {
		for _, f := range ref.GetFields() {
			project[ref.GetID()+"."+f.GetID()] = 1
		}
	}

	*q = append(*q, bson.M{
		"$project": project,
	})
	////////////////////////////////////////////////////////////////

	if ref.IsSingle() {
		*q = append(*q, bson.M{
			"$unwind": bson.M{
				"path":                       "$" + ref.GetID(),
				"preserveNullAndEmptyArrays": true,
			},
		})
	}
}

func (q *QueryPipeline) AddReverseReference(fields []coremodel.Field, ref coremodel.ReverseReferenceField) {

	// Sort can only be used on the singular fields.
	// If the reference is saved in an array, this sort and limit doesn't work.
	if ref.IsSingle() && ref.GetLimit() != 0 && ref.GetSort() != nil {
		*q = append(*q, bson.M{
			"$lookup": bson.M{
				"from": ref.GetOtherCollection(),
				"as":   ref.GetID(),
				"let":  bson.M{"id": "$_id"},
				"pipeline": bson.A{
					bson.M{"$match": bson.M{"$expr": bson.M{
						"$eq": bson.A{
							"$$id", "$" + ref.GetNameInOtherCollection() + "._id",
						}}},
					},
					bson.M{"$sort": bson.M{ref.GetSort().ID: ref.GetSort().Order}},
					bson.M{"$limit": ref.GetLimit()},
				},
			},
		})
	} else {
		*q = append(*q, bson.M{
			"$lookup": bson.M{
				"from":         ref.GetOtherCollection(),
				"localField":   "_id",
				"foreignField": ref.GetNameInOtherCollection() + "._id",
				"as":           ref.GetID(),
			},
		})
	}

	// filter
	if ref.GetFilter() != nil {
		project := bson.M{
			"_id": 1,
		}

		for _, field := range fields {
			if ref.GetID() != field.GetID() {
				project[field.GetID()] = 1
			}
		}
		project[ref.GetID()] = bson.M{
			"$filter": bson.M{
				"input": "$" + ref.GetID(),
				"as":    ref.GetID(),
				"cond":  ref.GetFilter().ToFilterQuery(ref.GetID()),
			},
		}

		*q = append(*q, bson.M{
			"$project": project,
		})
	}

	// filter fields START
	project := bson.M{
		"_id": 1,
	}

	// keep parent fields
	for _, field := range fields {
		if ref.GetID() != field.GetID() {
			project[field.GetID()] = 1
		}
	}

	// keep reference fields
	project[ref.GetID()+"._id"] = 1
	if ref.GetFields() != nil {
		for _, f := range ref.GetFields() {
			project[ref.GetID()+"."+f.GetID()] = 1
		}
	}

	*q = append(*q, bson.M{
		"$project": project,
	})
	// filter fields END

	if ref.IsSingle() {
		*q = append(*q, bson.M{
			"$unwind": bson.M{
				"path":                       "$" + ref.GetID(),
				"preserveNullAndEmptyArrays": true,
			},
		})
	}
}

func (q *QueryPipeline) AddRelationReference(fields []coremodel.Field, ref coremodel.RelationReferenceField) {

	// fetch relation collection
	*q = append(*q, bson.M{
		"$lookup": bson.M{
			"from":         ref.GetRelationCollection(),
			"localField":   "_id",
			"foreignField": ref.GetNameInRelationCollection() + "._id",
			"as":           ref.GetID(),
		},
	})

	// filter relation collection
	if ref.GetRelationFilter() != nil {
		project := bson.M{
			"_id": 1,
		}

		for _, field := range fields {
			project[field.GetID()] = 1
		}
		project[ref.GetID()] = bson.M{
			"$filter": bson.M{
				"input": "$" + ref.GetID(),
				"as":    ref.GetID(),
				"cond":  ref.GetRelationFilter().ToFilterQuery(ref.GetID()),
			},
		}

		*q = append(*q, bson.M{
			"$project": project,
		})
	}

	// separate relation collections to convert them to object (from array) and fetch other document
	*q = append(*q, bson.M{
		"$unwind": bson.M{
			"path":                       "$" + ref.GetID(),
			"preserveNullAndEmptyArrays": true,
		},
	})

	// fetch other collection
	*q = append(*q, bson.M{
		"$lookup": bson.M{
			"from":         ref.GetOtherCollection(),
			"localField":   ref.GetID() + "." + ref.GetNameOfOtherCollectionInRelationCollection() + "._id",
			"foreignField": "_id",
			"as":           ref.GetID() + "." + ref.GetNameOfOtherCollectionInRelationCollection(),
		},
	})

	// filter other collection
	if ref.GetOtherCollectionFilter() != nil {
		project := bson.M{
			ref.GetID() + "._id": 1,
			ref.GetID() + "." + ref.GetNameInRelationCollection(): 1,
		}

		// keep parent fields
		for _, field := range fields {
			if ref.GetID() != field.GetID() {
				project[field.GetID()] = 1
			}
		}

		// keep relation fields
		project[ref.GetID()+"._id"] = 1
		if ref.GetFields() != nil {
			for _, f := range ref.GetFields() {
				//if f.GetID() != ref.GetNameOfOtherCollectionInRelationCollection() {
				project[ref.GetID()+"."+f.GetID()] = 1
				//}
			}
		}

		project[ref.GetID()+"."+ref.GetNameOfOtherCollectionInRelationCollection()] = bson.M{
			"$filter": bson.M{
				"input": "$" + ref.GetID() + "." + ref.GetNameOfOtherCollectionInRelationCollection(),
				"as":    ref.GetID() + "_" + ref.GetNameOfOtherCollectionInRelationCollection(),
				"cond":  ref.GetOtherCollectionFilter().ToFilterQuery(ref.GetID() + "_" + ref.GetNameOfOtherCollectionInRelationCollection()),
			},
		}

		*q = append(*q, bson.M{
			"$project": project,
		})
	}

	// separate other collections to convert them to object (from array)
	*q = append(*q, bson.M{
		"$unwind": bson.M{
			"path":                       "$" + ref.GetID() + "." + ref.GetNameOfOtherCollectionInRelationCollection(),
			"preserveNullAndEmptyArrays": true,
		},
	})

	// remove the relations that doesn't contain the other collection
	*q = append(*q, bson.M{
		"$match": bson.M{"$or": bson.A{
			bson.M{ref.GetID() + "._id": bson.M{"$exists": false}},
			bson.M{ref.GetID() + "." + ref.GetNameOfOtherCollectionInRelationCollection(): bson.M{"$exists": true}},
		}},
	})

	// filter fields
	////////////////////////////////////////////////////////////////
	project := bson.M{
		"_id": 1,
	}

	// keep parent fields
	for _, field := range fields {
		if ref.GetID() != field.GetID() {
			project[field.GetID()] = 1
		}
	}

	// keep relation fields
	project[ref.GetID()+"._id"] = 1
	if ref.GetFields() != nil {
		for _, f := range ref.GetFields() {
			if f.GetID() != ref.GetNameOfOtherCollectionInRelationCollection() {
				project[ref.GetID()+"."+f.GetID()] = 1
			}
		}
	}

	// keep other collection fields
	project[ref.GetID()+"."+ref.GetNameOfOtherCollectionInRelationCollection()+"._id"] = 1
	if ref.GetOtherCollectionFields() != nil {
		for _, f := range ref.GetOtherCollectionFields() {
			project[ref.GetID()+"."+ref.GetNameOfOtherCollectionInRelationCollection()+"."+f.GetID()] = 1
		}
	}

	*q = append(*q, bson.M{
		"$project": project,
	})
	////////////////////////////////////////////////////////////////

	// merge relation collection results
	////////////////////////////////////////////////////////////////
	group := bson.M{
		"_id":       "$_id",
		ref.GetID(): bson.M{"$push": "$" + ref.GetID()},
	}

	// add fields to include in group
	for _, fieldToInclude := range fields {

		// check if the field id is not the field we're generating.
		// otherwise this will select the first item from the array and the array will be single object
		if fieldToInclude.GetID() != ref.GetID() {
			group[fieldToInclude.GetID()] = bson.M{"$first": "$" + fieldToInclude.GetID()}
		}
	}

	*q = append(*q, bson.M{
		"$group": group,
	})
	////////////////////////////////////////////////////////////////

	// remove arrays with only one empty objects: [{}]
	////////////////////////////////////////////////////////////////
	removeEmptyObjectStep := bson.M{
		"_id": 1,
	}

	// keep parent fields
	for _, field := range fields {
		if ref.GetID() != field.GetID() {
			removeEmptyObjectStep[field.GetID()] = 1
		}
	}

	// sets field value as empty array with $cond if array is [{}] = []map[string]interface{}{{}}
	removeEmptyObjectStep[ref.GetID()] = bson.M{
		"$cond": bson.A{
			bson.M{"$eq": bson.A{"$" + ref.GetID(), []map[string]interface{}{{}}}},
			[]map[string]interface{}{},
			"$" + ref.GetID(),
		},
	}
	*q = append(*q, bson.M{
		"$project": removeEmptyObjectStep,
	})
	// remove arrays with only one empty objects END
	////////////////////////////////////////////////////////////////
}

func (q *QueryPipeline) AddCollectionLookup(fields []coremodel.Field, ref coremodel.CollectionLookupField) {
	*q = append(*q, bson.M{
		"$lookup": bson.M{
			"from":         ref.GetFrom(),
			"localField":   ref.GetLocalField(),
			"foreignField": ref.GetForeignField(),
			"as":           ref.GetID(),
		},
	})

	// filter fields START
	project := bson.M{
		"_id": 1,
	}

	// keep parent fields
	for _, field := range fields {
		if ref.GetID() != field.GetID() {
			project[field.GetID()] = 1
		}
	}

	// keep reference fields
	project[ref.GetID()+"._id"] = 1
	if ref.GetFields() != nil {
		for _, f := range ref.GetFields() {
			project[ref.GetID()+"."+f.GetID()] = 1
		}
	}

	*q = append(*q, bson.M{
		"$project": project,
	})
	// filter fields END

	if ref.IsSingle() {
		*q = append(*q, bson.M{
			"$unwind": bson.M{
				"path":                       "$" + ref.GetID(),
				"preserveNullAndEmptyArrays": true,
			},
		})
	}
}

func (service DatabaseService) Query(
	base, collection string, config *coremodel.QueryConfig,
) ([]*coremodel.DataModel, *coremodel.Meta, error) {

	pipeline := NewPipeline()

	if config.Filter != nil {
		pipeline.AddMatch(config)
	}

	for _, field := range config.Fields {
		if field.GetType() == coremodel.FieldTypeReference {
			pipeline.AddReference(config.Fields, coremodel.ReferenceFromField(field))
		} else if field.GetType() == coremodel.FieldTypeReverseReference {
			pipeline.AddReverseReference(config.Fields, coremodel.ReverseReferenceFromField(field))
		} else if field.GetType() == coremodel.FieldTypeRelationReference {
			pipeline.AddRelationReference(config.Fields, coremodel.SingleRelationReferenceFromField(field))
		} else if field.GetType() == coremodel.FieldTypeCollectionLookup {
			pipeline.AddCollectionLookup(config.Fields, coremodel.CollectionLookupFieldFromField(field))
		}
	}

	if config.Sort != nil {
		for _, sort := range config.Sort {
			pipeline.AddSort(sort.ID, sort.Order)
		}
	}

	pipeline.AddSkip(config.Skip)

	// filter the fields
	pipeline.AddProject(config.Fields)

	pipeline.AddLimit(config.Limit)

	var response AggregateResult
	err := service.Database.Aggregate(base, collection, pipeline, func(cur *mongo.Cursor) error {

		err := cur.Decode(&response)
		if err != nil {
			return err
		}
		return nil
	})

	if response.Results == nil {
		response.Results = []*coremodel.DataModel{}
	}

	var meta *coremodel.Meta
	if response.MetaWrapper != nil && len(response.MetaWrapper) > 0 {
		meta = response.MetaWrapper[0]
	}

	return response.Results, meta, err
}
