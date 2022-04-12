# archive reference is nixpkgs-unstable from 2022-03-25
{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/30d3d79b7d3607d56546dd2a6b49e156ba0ec634.tar.gz") {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.docker
    pkgs.docker-compose
    pkgs.go_1_18
    pkgs.terraform
  ];
}
