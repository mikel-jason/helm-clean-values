#!/usr/bin/env bash

set -eo pipefail

pushd "$HELM_PLUGIN_DIR"

version="$(grep "version" plugin.yaml | cut -d '"' -f 2)"

uname="linux"
arch="amd64"

mkdir -p "bin"
mkdir -p "releases/${version}"

url="https://github.com/sarcaustech/helm-clean-values/releases/download/v${version}/helm-clean-values_${version}_${uname}_${arch}.tar.gz"

ls -l
ls -l releases

echo "Downloading ${url} to ./releases/${version}.tar.gz ($1)"

curl -sSL "${url}" -o "./releases/${version}.tar.gz"
tar xzf "./releases/${version}.tar.gz" -C "./releases/${version}"

mv "./releases/${version}/helm-clean-values" "./bin/"

popd
