.PHONY: deps test mocks cover sync-coveralls proto

deps:
	go mod download

covertest: deps
	go test  -coverprofile=coverage.out `go list ./... | grep -v pb`
	go tool cover -html=coverage.out

sync-coveralls: deps
	go test  -coverprofile=coverage.out `go list ./... | grep -v pb`
	goveralls -coverprofile=coverage.out -reponame=go-webdriver -repotoken=${COVERALLS_GO_CRX3_TOKEN} -service=local

proto:
	@protoc  --go_out=. ./pb/crx3.proto