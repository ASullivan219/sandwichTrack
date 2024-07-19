docker: docker_build docker_compose
run:
	go run ./cmd/main/
docker_build:
	@echo "building container"
	docker build -t sandwhich-tracker .
docker_compose:
	@echo "starting docker container"
	docker-compose up -d
