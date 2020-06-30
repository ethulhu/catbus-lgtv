# SPDX-FileCopyrightText: 2020 Ethel Morgan
#
# SPDX-License-Identifier: MIT

{ pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoModule rec {
  name = "catbus-lgtv-${version}";
  version = "latest";
  goPackagePath = "go.eth.moe/catbus-lgtv";

  modSha256 = "0fbrz71qh2rjnbdhj793v1ax3k2fzf7wayfw66j1wy8kwy0cq2l2";

  src = ./.;

  meta = {
    homepage = "https://ethulhu.co.uk/catbus";
    licence = stdenv.lib.licenses.mit;
  };
}
