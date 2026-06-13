lint:
    pre-commit run -a

test:
    go tool gotestsum -- -coverprofile=.coverage ./...

update:
    go get -u -t toolchain@XXX ./...

cov-render:
    go tool cover -html .coverage -o coverage.html
    xdg-open coverage.html
