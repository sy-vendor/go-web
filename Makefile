.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o apiserver -trimpath -ldflags "-s -w" go-web/cmd/apiserver

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: pack
pack-%: % build-%
	tar -czf $<.tar.gz $<

.PHONY: clean
clean-%: %
	rm $<

.PHONY: docker
docker: build
	docker build . -t go-web:latest

.PHONY: run-docker
run-docker:
	docker build . -t go-web:latest
	docker rm -f go-web | echo "remove ok"
	docker run -d --name go-web go-web
	docker ps
