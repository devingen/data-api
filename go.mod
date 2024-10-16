module github.com/devingen/data-api

go 1.12

//replace github.com/devingen/api-core => ../api-core

require (
	github.com/aws/aws-lambda-go v1.46.0
	github.com/devingen/api-core v0.1.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-resty/resty/v2 v2.12.0
	github.com/gorilla/mux v1.7.4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.7.2
	go.mongodb.org/mongo-driver v1.13.2
)
