{
  pkgs,
  lib,
  ...
}: {
  languages.go = {
    enable = true;
    package = pkgs.go;
  };

  packages = with pkgs; [
    lefthook
    golangci-lint
    gotestsum
    go-task
    gomarkdoc
    python313Packages.mkdocs-material
    python313Packages.mkdocs
  ];
}
