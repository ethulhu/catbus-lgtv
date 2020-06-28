# SPDX-FileCopyrightText: 2020 Ethel Morgan
#
# SPDX-License-Identifier: MIT

{ pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoModule rec {
  name = "catbus-lgtv-${version}";
  version = "latest";
  goPackagePath = "go.eth.moe/catbus-lgtv";

  modSha256 = "0533if636x71g3bly3wshx2y5f7ziww0inngpnaq3fjmr0ksqzaf";

  src = ./.;

  meta = {
    homepage = "https://ethulhu.co.uk/catbus";
    licence = stdenv.lib.licenses.mit;
  };
}
