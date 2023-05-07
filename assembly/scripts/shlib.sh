#!/bin/sh
cat /dev/null <<EOF
------------------------------------------------------------------------
https://github.com/client9/shlib - portable posix shell functions

Public domain - http://unlicense.org
https://github.com/client9/shlib/blob/master/LICENSE.md

but credit (and pull requests) appreciated.
------------------------------------------------------------------------
EOF
# write message to stderr
echoerr() {
	echo "$@" 1>&2
}
http_download_curl() {
	local_file=$1
	source_url=$2
	header=$3
	if [ -z "$header" ]; then
		code=$(curl -w '%{http_code}' -sL -o "$local_file" "$source_url")
	else
		code=$(curl -w '%{http_code}' -sL -H "$header" -o "$local_file" "$source_url")
	fi
	if [ "$code" != "200" ]; then
		log_debug "http_download_curl received HTTP status $code"
		return 1
	fi
	return 0
}

# http_download_wget
#
# unable to get server response code in a portable manner
# busybox wget (used on alpine linux) does not support "--server-response"
#
http_download_wget() {
	local_file=$1
	source_url=$2
	header=$3
	if [ -z "$header" ]; then
		wget -q -O "$local_file" "$source_url"
	else
		wget -q --header "$header" -O "$local_file" "$source_url"
	fi
}
#
# http_download [local-file] [url] [optional extra header]
#
# if arg3 is not empty it will add it as an extra HTTP header
# must be in the form "foo: bar"
#
http_download() {
	log_debug "http_download $2"
	if is_command curl; then
		http_download_curl "$@"
		return
	elif is_command wget; then
		http_download_wget "$@"
		return
	fi
	log_crit "http_download unable to find wget or curl"
	return 1
}

# http_copy - copies contents of a URL to stdout or fail
#
# needed since curl is broken
#
http_copy() {
	tmp=$(mktemp)
	http_download "${tmp}" "$1" "$2" || return 1
	body=$(cat "$tmp")
	rm -f "${tmp}"
	echo "$body"
}
#
# is_command: returns true if command exists
#
# `which` is not portable, in particular is often
# not available on RedHat/CentOS systems.
#
# `type` is implemented in many shells but technically not
# part of the posix spec.
#
# `command -v`
#
is_command() {
	command -v "$1" >/dev/null
	#type "$1" > /dev/null 2> /dev/null
}
# function to prefix each log output
#  over-ride to add custom output or format
#
# by default prints the script name ($0)
log_prefix() {
	echo "$0"
}

# default priority
_logp=6

# set the log priority
#  todo: fancy turn string into number
log_set_priority() {
	_logp="$1"
}

# if no args, return the priority
# if arg, then test if greater than or equals to priority
log_priority() {
	if test -z "$1"; then
		echo "$_logp"
		return
	fi
	[ "$1" -le "$_logp" ]
}

log_tag() {
	case $1 in
	0) echo "emerg" ;;
	1) echo "alert" ;;
	2) echo "crit" ;;
	3) echo "err" ;;
	4) echo "warning" ;;
	5) echo "notice" ;;
	6) echo "info" ;;
	7) echo "debug" ;;
	*) echo "$1" ;;
	esac
}

log_debug() {
	log_priority 7 || return 0
	echoerr "$(log_prefix)" "$(log_tag 7)" "$@"
}

log_info() {
	log_priority 6 || return 0
	echoerr "$(log_prefix)" "$(log_tag 6)" "$@"
}

log_err() {
	log_priority 3 || return 0
	echoerr "$(log_prefix)" "$(log_tag 3)" "$@"
}

# log_crit is for platform problems
log_crit() {
	log_priority 2 || return 0
	echoerr "$(log_prefix)" "$(log_tag 2)" "$@"
}
# uname_arch converts `uname -m` back into standardized golang
# OS types.
#
# See also `uname_arch_check` for a self-check
#
# ## NOTES
#
# Notes on ARM:
# arm 5,6,7: uname is of form `armv6l`, ` armv7l` where a letter
# or something else is after the number. Has examples:
# https://github.com/golang/go/wiki/GoArm
# https://en.wikipedia.org/wiki/List_of_ARM_microarchitectures
#
# arm 8 is know as arm64, and aarch64
#
# more notes: https://github.com/golang/go/issues/13669
#
# ## EXAMPLE
#
# ```bash
# ARCH=$(uname_arch)
# ```
#
#
uname_arch() {
	arch=$(uname -m)
	case $arch in
	x86_64) arch="amd64" ;;
	x86) arch="386" ;;
	i686) arch="386" ;;
	i386) arch="386" ;;
	aarch64) arch="arm64" ;;
	armv5*) arch="armv5" ;;
	armv6*) arch="armv6" ;;
	armv7*) arch="armv7" ;;
	esac
	echo "${arch}"
}
#!/bin/sh

# uname_os converts `uname -s` into standard golang OS types
# golang types are used since they cover
# most platforms and are standardized while raw uname values vary
# wildly.  See complete list of values by running
# "go tool dist list"
#
# ## EXAMPLE
#
# ```bash
# OS=$(uname_os)
# ```
#
uname_os() {
	os=$(uname -s | tr '[:upper:]' '[:lower:]')

	# fixed up for https://github.com/client9/shlib/issues/3
	case "$os" in
	msys*) os="windows" ;;
	mingw*) os="windows" ;;
	cygwin*) os="windows" ;;
	win*) os="windows" ;; # for windows busybox and like # https://frippery.org/busybox/
	esac

	# other fixups here
	echo "$os"
}
#
# untar: untar or unzip $1
#
# if you need to unpack in specific directory use a
# subshell and cd
#
# (cd /foo && untar mytarball.gz)
#
untar() {
	tarball=$1
	case "${tarball}" in
	*.tar.gz | *.tgz) tar -xzf "${tarball}" ;;
	*.tar) tar -xf "${tarball}" ;;
	*.zip) unzip "${tarball}" ;;
	*)
		log_err "untar unknown archive format for ${tarball}"
		return 1
		;;
	esac
}
#!/bin/sh
cat /dev/null <<EOF
------------------------------------------------------------------------
End of functions from https://github.com/client9/shlib
------------------------------------------------------------------------
EOF
