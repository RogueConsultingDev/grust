ARG GO_VERSION=tip

FROM golang:${GO_VERSION}

RUN set -eu; \
    apt update; \
    apt install -y --no-install-recommends curl git openssh-client zsh; \
    echo 'github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl' >> /etc/ssh/ssh_known_hosts;

RUN useradd --create-home --shell /bin/zsh dev;

USER dev

COPY zshrc /home/dev/.zshrc
