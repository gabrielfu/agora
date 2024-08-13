dry-release:
	goreleaser release --snapshot --clean

release:
	git tag -a $(TAG) -m $(TITLE)
	git push origin $(TAG)
	goreleaser release

test:
	go test ./...

test-cov:
	go test ./... -coverprofile=coverage.out
