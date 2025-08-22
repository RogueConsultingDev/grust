lint:
    pre-commit run -a

test:
    go test -coverprofile=.coverage ./...

cov-render:
    go tool cover -html .coverage -o coverage.html
    xdg-open coverage.html
