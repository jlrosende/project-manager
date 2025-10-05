FROM mcr.microsoft.com/devcontainers/go:1.25-bookworm

ENV SHELL=/bin/zsh

RUN curl -fsSL https://opencode.ai/install | bash && \
    mv $HOME/.opencode/bin/opencode /usr/local/bin/opencode && \
    curl -fsSL https://dl.dagger.io/dagger/install.sh | BIN_DIR=/usr/local/bin sh
