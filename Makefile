run-postgres:
	docker run --name pg -d -e POSTGRES_HOST_AUTH_METHOD=trust -p 5440:5432 postgres:10.21

allure-results:
	mkdir allure-results

start-postgres:
	docker start pg

stop-postgres:
	docker stop pg


test: allure-results
	go build -o build/main cmd/example/main.go
	go test -v ./tests/... -count=1 -run TestSuite


test-cover:
	go build -o build/main cmd/example/main.go
	go test -cover -coverpkg github.com/ixpectus/declarate/... -coverprofile cover.out -v ./tests/... -count=1 -run TestExample
	go tool cover -html=cover.out -o cover.html


run-polling: 
	go run cmd/example/main.go -dir ./tests/yaml_poll/poll_long.yaml -progress_bar

