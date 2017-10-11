all: test vet cover

vet:
	@go vet

test:
	@go test -v 

cover: 
	@go test -cover

cover-report:
	@go test -coverprofile=coverage.out
	@go tool cover -html=coverage.out
	@rm coverage.out
