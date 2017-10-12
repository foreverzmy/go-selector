all: test vet cover

spell:
	@misspell -locale US *.go *.md

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

tools:
	go get -u github.com/client9/misspell