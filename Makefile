.PHONY: test
test:
	go build -o build/main cmd/example/main.go
	go test -v ./tests/... -count=1 -run TestExample	


test-cover:
	go build -o build/main cmd/example/main.go
	go test -cover -coverpkg github.com/ixpectus/declarate/... -coverprofile cover.out -v ./tests/... -count=1 -run TestExample
	go tool cover -html=cover.out -o cover.html
