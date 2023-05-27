let
  pkgs = import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/feda52be1d59.tar.gz") { };

in
pkgs.mkShell {
  buildInputs = with pkgs; [
    git-ignore
  ];
}
