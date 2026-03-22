package crx3

//go:generate go install go.uber.org/mock/mockgen@latest

// mcp mocks
//go:generate mockgen -source=./mcp/service.go -destination=./mcp/service_mock_test.go -package=mcp
