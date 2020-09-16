.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/query aws/query/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy-mentornity: clean build
	serverless deploy --stage mentornity --region eu-central-1 --verbose

teardown-mentornity: clean
	serverless remove --stage mentornity --region eu-central-1 --verbose

deploy-devingen: clean build
	serverless deploy --stage devingen --region eu-central-1 --verbose

teardown-devingen: clean
	serverless remove --stage devingen --region eu-central-1 --verbose
