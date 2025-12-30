build:
	go build -o .out/restore-vercel-deployments main.go
run:
	make build
	./.out/restore-vercel-deployments
run-dev:
	go run main.go
build-docker:
	docker build -t restore-vercel-deployments .