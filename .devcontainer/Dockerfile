# Make sure GO_VERSION matches the Go version you want to use
ARG GO_VERSION=1.25
FROM mcr.microsoft.com/devcontainers/go:$GO_VERSION

USER root

# --- Install goose + sqlc
# RUN go install github.com/pressly/goose/v3/cmd/goose@latest
# RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# ENV PATH="/go/bin:${PATH}"

# --- Custom setup
USER vscode

RUN git clone https://github.com/zsh-users/zsh-autosuggestions \
    ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions

RUN sed -ir 's/plugins=(git)/plugins=(git zsh-autosuggestions)/g' ~/.zshrc && \
    sudo chsh -s $(which zsh) $(whoami)
# ---
