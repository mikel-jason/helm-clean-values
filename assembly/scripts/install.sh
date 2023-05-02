#!/usr/bin/env bash

set -eo pipefail

basePath="$1"
pushd $1

version="$(cat plugin.yaml | grep "version" | cut -d '"' -f 2)"

uname="Linux"
arch="amd64"

mkdir -p "{bin,releases/v${version}}"

url="https://github.com/sarcaustech/helm-clean-values/releases/download/v${version}/helm-clean-values_${version}_${uname}_${arch}.tar.gz"

curl -sSL "${url}" -o "releases/v${version}.tar.gz"
tar xzf "releases/v${version}.tar.gz" -C "releases/v${version}"

mv "releases/v${version}/plugin.yaml" .
mv "releases/v${version}/helm-clean-values" .

popd
