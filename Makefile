.PHONY: test
test:
	go build -o build/main cmd/example/main.go
	go test -v ./tests/... -count=1 -run TestExample	

