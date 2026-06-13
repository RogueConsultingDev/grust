lint:
    pre-commit run -a

test:
    go tool gotestsum -- -coverprofile=.coverage ./...

update:
    #!/usr/bin/env bash
    TOOLCHAIN="$(cat go.mod | grep '^go [0-9.]\+$' | cut -d ' ' -f2)"
    go get -u -t "toolchain@${TOOLCHAIN}" ./...

cov-render:
    go tool cover -html .coverage -o coverage.html
    xdg-open coverage.html
