default: test

test:
	go test ./...

update-changelog:
	conventional-changelog -p angular -i CHANGELOG.md -s