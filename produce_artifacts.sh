#!/usr/bin/env bash

arch=$(uname -m)
plat=$(echo $(uname -s) | awk '{print tolower($0)}')

CARGO_TARGET_DIR="deps/${plat}_${arch}" cargo build --release
