#!/bin/bash -ex
#
# This script is used to download and install the v8 libraries on linux for
# travis-ci.
#

: "${V8_VERSION:?V8_VERSION must be set}"

V8_DIR=${HOME}/libv8gem
if [ -d "$V8_DIR" ]; then
    echo "Assume that the library has been downloaded."
else
    mkdir -p ${V8_DIR}
    pushd ${V8_DIR}

    unameOut="$(uname -s)"
    case "${unameOut}" in
        Linux*)
            curl https://rubygems.org/downloads/libv8-${V8_VERSION}-$(uname -m)-$(uname -s | tr '[:upper:]' '[:lower:]').gem | tar xv
            ;;
        Darwin*)
            : "${DARWIN_VERSION:?DARWIN_VERSION must be set}"
            curl https://rubygems.org/downloads/libv8-${V8_VERSION}-$(uname -m)-$(uname -s | tr '[:upper:]' '[:lower:]')-${DARWIN_VERSION}.gem | tar xv
            ;;
        *)
            machine="UNKNOWN:${unameOut}"
            echo "Not supported platform: ${machine}"
            exit -1
    esac

    tar xzvf data.tar.gz

    popd

    rm -rf libv8
    rm -rf include

    ln -s ${V8_DIR}/vendor/v8/out/x64.release libv8
    ln -s ${V8_DIR}/vendor/v8/include include
fi

go get ./...
