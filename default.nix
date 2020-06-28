# SPDX-FileCopyrightText: 2020 Ethel Morgan
#
# SPDX-License-Identifier: MIT

{ pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoModule rec {
  name = "catbus-lgtv-${version}";
  version = "latest";
  goPackagePath = "go.eth.moe/catbus-lgtv";

  modSha256 = "1xy5zm1g9w5v28mrk1i06971slfawiq7lgr7na5v3lcvqz2r4mf5";

  src = ./.;

  meta = {
    homepage = "https://ethulhu.co.uk/catbus";
    licence = stdenv.lib.licenses.mit;
  };
}
