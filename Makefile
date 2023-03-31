.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build-%: %
	GOOS=linux GOARCH=amd64 go build -o $< -trimpath -ldflags "-s -w" go-web/cmd/$<

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: pack
pack-%: % build-%
	tar -czf $<.tar.gz $<

.PHONY: clean
clean-%: %
	rm $<
