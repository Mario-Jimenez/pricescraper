build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pricescraper-srv main.go

run:
	./pricescraper-srv

docker-build:
	docker build . -t docker.pkg.github.com/mario-jimenez/pricescraper/pricescraper-srv:0.0.1

docker-push:
	docker push docker.pkg.github.com/mario-jimenez/pricescraper/pricescraper-srv:0.0.1
