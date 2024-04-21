{ pkgs, ... }:

{
  packages = [ pkgs.git ];

  enterShell = ''
    git --version
    go version
  '';

  enterTest = ''
    go test gocog/processor
  '';

  languages.go.enable = true;

  pre-commit.hooks.golangci-lint.enable = true;
  pre-commit.hooks.gofmt.enable = true;
}
