# SPDX-FileCopyrightText: 2020 Ethel Morgan
#
# SPDX-License-Identifier: MIT

{ pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoModule rec {
  name = "catbus-lgtv-${version}";
  version = "latest";
  goPackagePath = "go.eth.moe/catbus-lgtv";

  modSha256 = "1s34h495h40vs90c79vdzk6pmpzskbmq0dqn4nw4iq7880i5s0z2";

  src = ./.;

  meta = {
    homepage = "https://ethulhu.co.uk/catbus";
    licence = stdenv.lib.licenses.mit;
  };
}
