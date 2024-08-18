.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/create/bootstrap aws/create/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/update/bootstrap aws/update/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/delete/bootstrap aws/delete/main.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/query/bootstrap aws/query/main.go

zip:
	zip -j bin/create.zip bin/create/bootstrap
	zip -j bin/update.zip bin/update/bootstrap
	zip -j bin/delete.zip bin/delete/bootstrap
	zip -j bin/query.zip bin/query/bootstrap

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy-devingen: clean build zip
	serverless deploy --stage prod --region eu-central-1 --verbose

teardown-devingen: clean
	serverless remove --stage prod --region eu-central-1 --verbose

deploy-devingen-dev: clean build
	serverless deploy --stage dev --region ca-central-1 --verbose

teardown-devingen-dev: clean
	serverless remove --stage dev --region ca-central-1 --verbose
