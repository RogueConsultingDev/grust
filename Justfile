lint:
    pre-commit run -a

test:
    go tool gotestsum -- -coverprofile=.coverage ./...

cov-render:
    go tool cover -html .coverage -o coverage.html
    xdg-open coverage.html
