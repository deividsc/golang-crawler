
test:
	go test ./... --coverprofile cover.out


craw:
	URL_VISIT=$(URL) go run cmd/main.go