{
  pkgs,
  lib,
  config,
  inputs,
  ...
}: {
  languages.go = {
    enable = true;
    package = pkgs.go;
  };

  packages = [
    pkgs.lefthook
    pkgs.golangci-lint
    pkgs.gotestsum
    pkgs.go-task
  ];

  enterShell = ''
    go version
  '';
}
