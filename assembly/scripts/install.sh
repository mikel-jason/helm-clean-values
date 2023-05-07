#!/usr/bin/env bash

set -eo pipefail

scriptDir=$(dirname "$0")
source "${scriptDir}/shlib.sh"

log_set_priority "${LOG_PRIO:-6}"

pushd "$HELM_PLUGIN_DIR" >/dev/null

version="$(grep "version" plugin.yaml | cut -d '"' -f 2)"

os=$(uname_os)
arch=$(uname_arch)

log_info "Installing helm-clean-values version v${version} for ${os} (${arch})"

mkdir -p "bin"
mkdir -p "releases/${version}"

url="https://github.com/sarcaustech/helm-clean-values/releases/download/v${version}/helm-clean-values_${version}_${os}_${arch}.tar.gz"

fileDir="./releases"
fileName="${version}.tar.gz"
filePath="${fileDir}/${fileName}"

http_download "${filePath}" "${url}"
(cd "${fileDir}" && untar "${fileName}")

binName="helm-clean-values"
if [ "${os}" == "windows" ]; then
	binName="${binName}.exe"
fi

log_debug "expected binary file name: ${binName}"

mv "${fileDir}/${binName}" "./bin/"

popd >/dev/null
