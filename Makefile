.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/query aws/query/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy-mentornity: clean build
	serverless deploy --stage mentornity --region eu-west-1 --verbose
