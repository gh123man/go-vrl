#!/usr/bin/env bash

arch=$(uname -m)
if [[ "$arch" == "aarch64" ]]; then
    arch="arm64"
fi
plat=$(echo $(uname -s) | awk '{print tolower($0)}')

CARGO_TARGET_DIR="deps/${plat}_${arch}" cargo build --release
