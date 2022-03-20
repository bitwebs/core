#!/usr/bin/env sh

BINARY=/iqd/${BINARY:-iqd}
ID=${ID:-0}
LOG=${LOG:-iqd.log}

if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'iqd'"
	exit 1
fi

BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"

if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

export IQDHOME="/iqd/node${ID}/iqd"

if [ -d "$(dirname "${IQDHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${IQDHOME}" "$@" | tee "${IQDHOME}/${LOG}"
else
  "${BINARY}" --home "${IQDHOME}" "$@"
fi
