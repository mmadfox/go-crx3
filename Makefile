.PHONY: deps test mocks cover sync-coveralls docker-protoc proto

deps:
	go mod download

test/cover: deps
	go test  -coverprofile=coverage.out `go list ./... | grep -v pb`
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

coveralls: deps
	go test  -coverprofile=coverage.out `go list ./... | grep -v pb`
	goveralls -coverprofile=coverage.out -reponame=go-webdriver -repotoken=${COVERALLS_GO_CRX3_TOKEN} -service=local

docker-protoc:
	@docker build -t github.com/mediabuyerbot/go-crx3/protoc:latest -f   \
           ./docker/protoc.dockerfile .

proto: docker-protoc
	@bash ./scripts/protoc.bash