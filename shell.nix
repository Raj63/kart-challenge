{ pkgs ? import (builtins.fetchGit {
  name = "dev-go";
  url = "https://github.com/NixOS/nixpkgs";
  ref = "refs/heads/master";
  rev = "af9aa247e052c8eac01e055461d5b3626ff603a7";
}) {} }:

with pkgs;

mkShell {
  CGO_ENABLED = "0";

  buildInputs = [
    git
    go_1_23
    go-tools
    golangci-lint
    goreleaser
    gosec
    gotools
    gofumpt
    golint
    pre-commit
    awscli2
    act
    gitlint
  ];

  shellHook = ''
    export PATH="$(go env GOPATH)/bin:$PATH"
    export GOPROXY="https://proxy.golang.org,direct"
    export GOSUMDB="sum.golang.org"

    echo "Installing Go tools..."
    go install github.com/go-critic/go-critic/cmd/gocritic@v0.6.5
    go install github.com/sqs/goreturns@v0.0.0-20210930122003-39244a879540
    go install github.com/swaggo/swag/cmd/swag@v1.16.4
    go install go.uber.org/mock/mockgen@latest

    pre-commit install
    git config --local include.path ../.gitconfig

    clear
  '';
}
