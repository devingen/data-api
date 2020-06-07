package integrationtests

import (
	"github.com/devingen/api-core/database"
	coremodel "github.com/devingen/api-core/model"
	"github.com/devingen/api-core/util"
	"github.com/devingen/data-api/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

type QueriesTestSuite struct {
	suite.Suite
	service service.DataService
	base    string
}

func TestQueries(t *testing.T) {
	db, err := database.NewDatabaseWithURI("mongodb://localhost")
	if err != nil {
		log.Fatalf("Database connection failed %s", err.Error())
	}

	testSuite := &QueriesTestSuite{
		service: service.NewDatabaseService(db),
		base:    "dvn-data-api-integration-test",
	}

	util.InsertDataFromFile(db, testSuite.base, "users")
	util.InsertDataFromFile(db, testSuite.base, "spaces")
	util.InsertDataFromFile(db, testSuite.base, "memberships")
	util.InsertDataFromFile(db, testSuite.base, "groups")

	suite.Run(t, testSuite)
}

func (suite *QueriesTestSuite) TestQuery() {
	results, _, err := suite.service.Query(
		suite.base,
		"users",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "firstName"),
				coremodel.New(coremodel.FieldTypeText, "lastName"),
			},
		},
	)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 10, len(results))

	util.SaveResultFile("query", results)
}

func (suite *QueriesTestSuite) TestFields() {
	results, _, err := suite.service.Query(
		suite.base,
		"users",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "firstName"),
				coremodel.New(coremodel.FieldTypeText, "lastName"),
			},
		},
	)
	util.SaveResultFile("query-fields", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 10, len(results))

	// items should only contain _id, firstName and lastName fields
	for _, item := range results {
		assert.Equal(suite.T(), 3, item.GetFieldCount())
		assert.NotEqual(suite.T(), "", item.GetID())
		assert.NotEqual(suite.T(), "", item.GetString("firstName"))
		assert.NotEqual(suite.T(), "", item.GetString("lastName"))
	}
}

func (suite *QueriesTestSuite) TestFilter() {
	results, _, err := suite.service.Query(
		suite.base,
		"groups",
		&coremodel.QueryConfig{
			Filter: &coremodel.Filter{
				Comparison: "eq",
				FieldId:    "tags",
				FieldValue: "automotive",
			},
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "tags"),
				coremodel.New(coremodel.FieldTypeText, "title"),
			},
		},
	)
	util.SaveResultFile("query-filter", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2, len(results))

	// items should only contain _id and title fields
	for _, item := range results {
		assert.Equal(suite.T(), 3, item.GetFieldCount())
		assert.NotEqual(suite.T(), "", item.GetID())
		assert.NotEqual(suite.T(), "", item.GetString("title"))
	}
}

func (suite *QueriesTestSuite) TestReferenceSingle() {
	results, _, err := suite.service.Query(
		suite.base,
		"groups",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "title"),
				coremodel.New(coremodel.FieldTypeText, "status"),
				coremodel.NewReference(
					"space",
					"spaces",
					true,
				).ToField(),
			},
		},
	)
	util.SaveResultFile("query-reference-single", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 6, len(results))

	for _, item := range results {
		space := item.GetChildDataModel("space")

		// all groups except 6 must have space
		if item.GetID() == "6" {
			assert.Nil(suite.T(), space)
			continue
		}

		assert.NotNil(suite.T(), space)
		assert.NotEqual(suite.T(), "", space.GetID())
	}
}

func (suite *QueriesTestSuite) TestReferenceSingleFields() {
	results, _, err := suite.service.Query(
		suite.base,
		"groups",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "title"),
				coremodel.New(coremodel.FieldTypeText, "status"),
				coremodel.NewReference(
					"space",
					"spaces",
					true,
				).SetFields([]coremodel.Field{
					coremodel.New(coremodel.FieldTypeText, "title"),
				}).ToField(),
			},
		},
	)
	util.SaveResultFile("query-reference-single-fields", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 6, len(results))

	for _, item := range results {
		space := item.GetChildDataModel("space")

		// all groups except 6 must have space
		if item.GetID() == "6" {
			assert.Nil(suite.T(), space)
			continue
		}

		// spaces must have only _id and title fields
		assert.NotNil(suite.T(), space)
		assert.Equal(suite.T(), 2, space.GetFieldCount())
		assert.NotEqual(suite.T(), "", space.GetID())
		assert.NotEqual(suite.T(), "", space.GetString("title"))
	}
}

func (suite *QueriesTestSuite) TestReferenceSingleFilter() {
	results, _, err := suite.service.Query(
		suite.base,
		"groups",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "title"),
				coremodel.New(coremodel.FieldTypeText, "status"),
				coremodel.NewReference(
					"space",
					"spaces",
					true,
				).SetFields([]coremodel.Field{
					coremodel.New(coremodel.FieldTypeText, "title"),
				}).SetFilter(&coremodel.Filter{
					Comparison: "eq",
					FieldId:    "type",
					FieldValue: "incubator",
				}).ToField(),
			},
		},
	)
	util.SaveResultFile("query-reference-single-filter", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 3, len(results))

	for _, item := range results {
		space := item.GetChildDataModel("space")

		// spaces must have only _id, title and type fields
		assert.NotNil(suite.T(), space)
		assert.Equal(suite.T(), 2, space.GetFieldCount())
		assert.NotEqual(suite.T(), "", space.GetID())
		assert.NotEqual(suite.T(), "", space.GetString("title"))
	}
}

func (suite *QueriesTestSuite) TestReferenceMultiple() {
	results, _, err := suite.service.Query(
		suite.base,
		"groups",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "title"),
				coremodel.New(coremodel.FieldTypeText, "status"),
				coremodel.NewReference(
					"members",
					"users",
					false,
				).SetFields([]coremodel.Field{
					coremodel.New(coremodel.FieldTypeText, "city"),
				}).ToField(),
			},
		},
	)
	util.SaveResultFile("query-reference-multiple", results)

	assert.Nil(suite.T(), err)
	//assert.Equal(suite.T(), 5, len(results))

	for _, item := range results {
		members := item.GetChildDataModelArray("members")
		assert.NotNil(suite.T(), members)

		if item.GetID() == "1" {
			assert.Equal(suite.T(), 1, len(members))
		} else if item.GetID() == "2" {
			assert.Equal(suite.T(), 2, len(members))
		} else if item.GetID() == "3" {
			assert.Equal(suite.T(), 2, len(members))
		} else if item.GetID() == "4" {
			assert.Equal(suite.T(), 0, len(members))
		} else if item.GetID() == "5" {
			assert.Equal(suite.T(), 0, len(members))
		}
	}
}

func (suite *QueriesTestSuite) TestReferenceMultipleFields() {
	results, _, err := suite.service.Query(
		suite.base,
		"groups",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "title"),
				coremodel.New(coremodel.FieldTypeText, "status"),
				coremodel.NewReference(
					"members",
					"users",
					false,
				).SetFields([]coremodel.Field{
					coremodel.New(coremodel.FieldTypeText, "firstName"),
					coremodel.New(coremodel.FieldTypeText, "city"),
				}).ToField(),
			},
		},
	)
	util.SaveResultFile("query-reference-multiple-fields", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 6, len(results))

	for _, item := range results {
		members := item.GetChildDataModelArray("members")

		for _, member := range members {
			assert.NotEqual(suite.T(), "", member.GetID())
			assert.NotEqual(suite.T(), "", member.GetString("firstName"))

			if member.GetID() == "4" {
				// user 4 has no city
				assert.Equal(suite.T(), 2, member.GetFieldCount())
			} else {
				assert.Equal(suite.T(), 3, member.GetFieldCount())
				assert.NotEqual(suite.T(), "", member.GetString("city"))
			}
		}
	}
}

func (suite *QueriesTestSuite) TestReferenceMultipleFilter() {
	results, _, err := suite.service.Query(
		suite.base,
		"groups",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "title"),
				coremodel.New(coremodel.FieldTypeText, "status"),
				coremodel.NewReference(
					"members",
					"users",
					false,
				).SetFilter(&coremodel.Filter{
					Comparison: "eq",
					FieldId:    "city",
					FieldValue: "London",
				}).ToField(),
			},
		},
	)
	util.SaveResultFile("query-reference-multiple-filter", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 3, len(results))

	for _, item := range results {
		members := item.GetChildDataModelArray("members")
		assert.Equal(suite.T(), 1, len(members))
	}
}

func (suite *QueriesTestSuite) TestReverseReferenceSingle() {
	results, _, err := suite.service.Query(
		suite.base,
		"users",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "firstName"),
				coremodel.New(coremodel.FieldTypeText, "lastName"),
				coremodel.NewReverseReference(
					"group",
					"groups",
					"members",
					true,
				).ToField(),
			},
		},
	)
	util.SaveResultFile("query-reverse-reference-single", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 10, len(results))

	for _, item := range results {
		group := item.GetChildDataModel("group")

		// these users must have groups
		if item.GetID() == "3" || item.GetID() == "4" || item.GetID() == "5" || item.GetID() == "6" || item.GetID() == "7" {
			assert.NotNil(suite.T(), group)
		}
	}
}

// Sort can only be used on the singular fields.
// If the reference is saved in an array, this sort and limit doesn't work
func (suite *QueriesTestSuite) TestReverseReferenceSingleWithSort() {
	results, _, err := suite.service.Query(
		suite.base,
		"spaces",
		&coremodel.QueryConfig{
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
		},
	)
	util.SaveResultFile("query-reverse-reference-single-with-sort", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 3, len(results))

	for _, item := range results {
		group := item.GetChildDataModel("lastGroup")

		// these spaces must have last groups
		if item.GetID() == "1" || item.GetID() == "2" {
			assert.NotNil(suite.T(), group)
		}

		if item.GetID() == "1" {
			assert.NotNil(suite.T(), group)
			assert.Equal(suite.T(), "3", group.GetID())
		} else if item.GetID() == "2" {
			assert.NotNil(suite.T(), group)
			assert.Equal(suite.T(), "5", group.GetID())
		}
	}
}

func (suite *QueriesTestSuite) TestReverseReferenceSingleFields() {
	results, _, err := suite.service.Query(
		suite.base,
		"users",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "firstName"),
				coremodel.New(coremodel.FieldTypeText, "lastName"),
				coremodel.NewReverseReference(
					"group",
					"groups",
					"members",
					true,
				).SetFields([]coremodel.Field{
					coremodel.New(coremodel.FieldTypeText, "title"),
				}).ToField(),
			},
		},
	)
	util.SaveResultFile("query-reverse-reference-single-fields", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 10, len(results))

	// groups of users should only contain _id and title fields
	for _, item := range results {
		group := item.GetChildDataModel("group")

		if group != nil {
			assert.Equal(suite.T(), 2, group.GetFieldCount())
			assert.NotEqual(suite.T(), "", group.GetString("_id"))
			assert.NotEqual(suite.T(), "", group.GetString("title"))
		}
	}
}

func (suite *QueriesTestSuite) TestReverseReferenceSingleFilter() {
	results, _, err := suite.service.Query(
		suite.base,
		"users",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "firstName"),
				coremodel.New(coremodel.FieldTypeText, "lastName"),
				coremodel.NewReverseReference(
					"membership",
					"memberships",
					"user",
					true,
				).SetFilter(&coremodel.Filter{
					Comparison: "ne",
					FieldId:    "status",
					FieldValue: "ended",
				}).ToField(),
			},
		},
	)
	util.SaveResultFile("query-reverse-reference-single-filter", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 10, len(results))

	for _, item := range results {
		membership := item.GetChildDataModel("membership")

		// these users must have memberships
		if item.GetID() == "1" || item.GetID() == "2" || item.GetID() == "3" || item.GetID() == "4" || item.GetID() == "6" {
			assert.NotNil(suite.T(), membership)
		}
	}
}

func (suite *QueriesTestSuite) TestReverseReferenceMultiple() {
	results, _, err := suite.service.Query(
		suite.base,
		"spaces",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "title"),
				coremodel.NewReverseReference(
					"groups",
					"groups",
					"space",
					false,
				).SetFields([]coremodel.Field{
					coremodel.New(coremodel.FieldTypeText, "title"),
				}).ToField(),
			},
		},
	)
	util.SaveResultFile("query-reverse-reference-multiple", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 3, len(results))

	for _, item := range results {
		groups := item.GetChildDataModelArray("groups")

		if item.GetID() == "1" {
			// there should be three groups
			assert.NotNil(suite.T(), groups)
			assert.Equal(suite.T(), 3, len(groups))
		} else if item.GetID() == "2" {
			// there should be one groups
			assert.NotNil(suite.T(), groups)
			assert.Equal(suite.T(), 2, len(groups))
		} else if item.GetID() == "3" {
			// there should be no groups
			assert.NotNil(suite.T(), groups)
			assert.Equal(suite.T(), 0, len(groups))
		}
	}
}

func (suite *QueriesTestSuite) TestReverseReferenceMultipleFields() {
	results, _, err := suite.service.Query(
		suite.base,
		"spaces",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "title"),
				coremodel.NewReverseReference(
					"groups",
					"groups",
					"space",
					false,
				).SetFields([]coremodel.Field{
					coremodel.New(coremodel.FieldTypeText, "title"),
				}).ToField(),
			},
		},
	)
	util.SaveResultFile("query-reverse-reference-multiple-fields", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 3, len(results))

	// groups of spaces should only contain _id and title fields
	for _, item := range results {
		groups := item.GetChildDataModelArray("groups")

		for _, group := range groups {
			assert.NotNil(suite.T(), group)

			assert.Equal(suite.T(), 2, group.GetFieldCount())
			assert.NotEqual(suite.T(), "", group.GetString("_id"))
			assert.NotEqual(suite.T(), "", group.GetString("title"))
		}
	}
}

func (suite *QueriesTestSuite) TestReverseReferenceMultipleFilter() {
	results, _, err := suite.service.Query(
		suite.base,
		"spaces",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "title"),
				coremodel.NewReverseReference(
					"activeStartups",
					"groups",
					"space",
					false,
				).SetFilter(&coremodel.Filter{
					Filters: []coremodel.Filter{{
						Comparison: coremodel.ComparisonNe,
						FieldId:    "status",
						FieldValue: "passive",
					}},
					Operator: coremodel.OperatorAnd,
				}).SetFields([]coremodel.Field{
					coremodel.New(coremodel.FieldTypeText, "title"),
				}).ToField(),
			},
		},
	)
	util.SaveResultFile("query-reverse-reference-multiple-filter", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 3, len(results))

	for _, item := range results {
		groups := item.GetChildDataModelArray("activeStartups")

		if item.GetID() == "1" {
			// there should be two groups
			assert.NotNil(suite.T(), groups)
			assert.Equal(suite.T(), 2, len(groups))
		} else if item.GetID() == "2" {
			// there should be one group
			assert.NotNil(suite.T(), groups)
			assert.Equal(suite.T(), 1, len(groups))
		} else if item.GetID() == "3" {
			// there should be no groups
			assert.NotNil(suite.T(), groups)
			assert.Equal(suite.T(), 0, len(groups))
		}
	}

}

func (suite *QueriesTestSuite) TestRelationReference() {
	results, _, err := suite.service.Query(
		suite.base,
		"spaces",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "title"),
				coremodel.NewRelationReference(
					"members",
					"memberships",
					"space",
					"users",
					"user",
				).ToField(),
			},
		},
	)
	util.SaveResultFile("query-relation-reference", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 3, len(results))

	for _, item := range results {
		members := item.GetChildDataModelArray("members")

		if item.GetID() == "1" {
			// there should be five members
			assert.NotNil(suite.T(), members)
			assert.Equal(suite.T(), 5, len(members))
		} else if item.GetID() == "2" {
			// there should be two members
			assert.NotNil(suite.T(), members)
			assert.Equal(suite.T(), 2, len(members))
		} else if item.GetID() == "3" {
			// there should be no members
			assert.NotNil(suite.T(), members)
			assert.Equal(suite.T(), 0, len(members))
		}
	}
}

func (suite *QueriesTestSuite) TestRelationReferenceOtherCollectionFields() {
	results, _, err := suite.service.Query(
		suite.base,
		"spaces",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "title"),
				coremodel.NewRelationReference(
					"memberships",
					"memberships",
					"space",
					"users",
					"user",
				).SetFields([]coremodel.Field{
					coremodel.New(coremodel.FieldTypeText, "status"),
				}).SetOtherCollectionFields([]coremodel.Field{
					coremodel.New(coremodel.FieldTypeText, "firstName"),
					coremodel.New(coremodel.FieldTypeText, "lastName"),
					coremodel.New(coremodel.FieldTypeText, "city"),
				}).ToField(),
			},
		},
	)
	util.SaveResultFile("query-relation-reference-other-collection-fields", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 3, len(results))

	// users in memberships in spaces should only contain _id, firstName, lastName and city fields
	for _, item := range results {
		memberships := item.GetChildDataModelArray("memberships")

		for _, member := range memberships {
			assert.NotNil(suite.T(), member)
			user := member.GetChildDataModel("user")

			// user 4 has not city info
			if user.GetID() != "4" {
				assert.Equal(suite.T(), 4, user.GetFieldCount())
				assert.NotEqual(suite.T(), "", user.GetString("_id"))
				assert.NotEqual(suite.T(), "", user.GetString("firstName"))
				assert.NotEqual(suite.T(), "", user.GetString("lastName"))
				assert.NotEqual(suite.T(), "", user.GetString("city"))
			} else {
				assert.Equal(suite.T(), 3, user.GetFieldCount())
				assert.NotEqual(suite.T(), "", user.GetString("_id"))
				assert.NotEqual(suite.T(), "", user.GetString("firstName"))
				assert.NotEqual(suite.T(), "", user.GetString("lastName"))
			}
		}
	}
}

func (suite *QueriesTestSuite) TestRelationReferenceOtherCollectionFilter() {
	results, _, err := suite.service.Query(
		suite.base,
		"spaces",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "title"),
				coremodel.NewRelationReference(
					"members",
					"memberships",
					"space",
					"users",
					"user",
				).SetFields([]coremodel.Field{
					coremodel.New(coremodel.FieldTypeText, "status"),
					coremodel.New(coremodel.FieldTypeText, "user"),
				}).SetOtherCollectionFilter(&coremodel.Filter{
					Comparison: coremodel.ComparisonEq,
					FieldId:    "city",
					FieldValue: "Istanbul",
				}).ToField(),
			},
		},
	)
	util.SaveResultFile("query-relation-reference-other-collection-filter", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 3, len(results))

	for _, item := range results {
		members := item.GetChildDataModelArray("members")

		if item.GetID() == "1" {
			assert.NotNil(suite.T(), members)
			assert.Equal(suite.T(), 2, len(members))
		} else if item.GetID() == "3" {
			assert.NotNil(suite.T(), members)
			assert.Equal(suite.T(), 0, len(members))
		}
	}
}

func (suite *QueriesTestSuite) TestRelationReferenceRelationFields() {
	results, _, err := suite.service.Query(
		suite.base,
		"spaces",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "title"),
				coremodel.NewRelationReference(
					"memberships",
					"memberships",
					"space",
					"users",
					"user",
				).SetFields([]coremodel.Field{
					coremodel.New(coremodel.FieldTypeText, "status"),
				}).ToField(),
			},
		},
	)
	util.SaveResultFile("query-relation-reference-relation-fields", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 3, len(results))

	// memberships of spaces should only contain _id, status and user fields
	for _, item := range results {
		memberships := item.GetChildDataModelArray("memberships")

		for _, member := range memberships {
			assert.NotNil(suite.T(), member)

			assert.Equal(suite.T(), 3, member.GetFieldCount())
			assert.NotEqual(suite.T(), "", member.GetString("_id"))
			assert.NotEqual(suite.T(), "", member.GetString("status"))
			assert.NotNil(suite.T(), member.GetChildDataModel("user"))
		}
	}
}

func (suite *QueriesTestSuite) TestRelationReferenceRelationFilter() {
	results, _, err := suite.service.Query(
		suite.base,
		"spaces",
		&coremodel.QueryConfig{
			Limit: 1000,
			Fields: []coremodel.Field{
				coremodel.New(coremodel.FieldTypeText, "title"),
				coremodel.NewRelationReference(
					"memberships",
					"memberships",
					"space",
					"users",
					"user",
				).SetRelationFilter(&coremodel.Filter{
					Comparison: coremodel.ComparisonNe,
					FieldId:    "status",
					FieldValue: "ended",
				}).ToField(),
			},
		},
	)
	util.SaveResultFile("query-relation-reference-relation-filter", results)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 3, len(results))

	for _, item := range results {
		memberships := item.GetChildDataModelArray("memberships")

		if item.GetID() == "1" {
			// there should be four memberships
			assert.NotNil(suite.T(), memberships)
			assert.Equal(suite.T(), 4, len(memberships))
		} else if item.GetID() == "2" {
			// there should be one member
			assert.NotNil(suite.T(), memberships)
			assert.Equal(suite.T(), 1, len(memberships))
		} else if item.GetID() == "3" {
			// there should be no memberships
			assert.NotNil(suite.T(), memberships)
			assert.Equal(suite.T(), 0, len(memberships))
		}
	}
}
