# SPDX-FileCopyrightText: 2020 Ethel Morgan
#
# SPDX-License-Identifier: MIT

{ pkgs ? import <nixpkgs> {} }:
with pkgs;

buildGoModule rec {
  name = "catbus-lgtv-${version}";
  version = "latest";
  goPackagePath = "go.eth.moe/catbus-lgtv";

  modSha256 = "1ww9azrgbkcavybx97y1x0ajck8l60g6n7spv4zjnxj52cra85p4";

  src = ./.;

  meta = {
    homepage = "https://ethulhu.co.uk/catbus";
    licence = stdenv.lib.licenses.mit;
  };
}
