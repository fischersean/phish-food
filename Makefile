.PHONY: test fmt lint showcover

showcover:
	@go tool cover -html=cp.out
	
lint:
	@make fmt
	@golangci-lint --skip-dirs cdk.out run

test:
	@make lint && go test -coverprofile cp.out $$(go list ./... | grep -v /cdk.out/)

run:
	@go run cmd/lambda/$(func)/main.go
			
fmt:
	@gofmt -s -w .

build:
	@cdk synth

deploy:
	@make lint
	@cdk deploy
	@find . -name 'asset.*' -print0 | xargs -0 rm -r

docker:
	@docker build -f Dockerfile.etl -t phishfoodetl .

sambuild:
	@cdk synth --no-staging > template.yaml
