# Overview

Ominous Positivity Backend, written in Go meant to powered by AWS Lambda APIGW

## Local Dev
Requires [SAM-CLI, DynamoDB Local]
```shell
sam build
sam local start-api --docker-network host -p 8080 
```

## Testing

### Unit
> go test -v ./...

### Integration
These also require DynamoDB, but running on localhost
> go test -v -tags integration ./message


## Packaging
```shell
cd message
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go
zip function.zip bootstrap
aws s3 function.zip <destination>
```


