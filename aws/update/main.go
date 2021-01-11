package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	coreaws "github.com/devingen/api-core/aws"
	"github.com/devingen/api-core/util"
	"github.com/devingen/data-api/aws"
	"github.com/devingen/data-api/controller"
	"github.com/devingen/data-api/dto"
	"github.com/devingen/data-api/service"
	"net/http"
)

func main() {

	db := aws.GetDatabase()
	databaseService := service.NewDatabaseService(db)
	serviceController := controller.NewServiceController(databaseService)

	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		var config dto.UpdateConfig
		err := json.Unmarshal([]byte(req.Body), &config)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		result, err := serviceController.Update(&dto.UpdateRequest{
			Base:                req.PathParameters["base"],
			Collection:          req.PathParameters["collection"],
			ID:                  req.PathParameters["id"],
			UpdateConfig:        &config,
			AuthorizationHeader: req.Headers["Authorization"],
		})
		response, err := util.BuildResponse(http.StatusOK, result, err)
		return coreaws.AdaptResponse(response, err)
	})
}
