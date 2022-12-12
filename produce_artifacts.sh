#!/usr/bin/env bash

arch=$(uname -m)
if [[ "$arch" == "aarch64" ]]; then
    arch="arm64"
fi
plat=$(echo $(uname -s) | awk '{print tolower($0)}')

$(cd v5; cargo build --release)
$(cd v10; cargo build --release --target wasm32-wasi)

cp "v5/target/release/libvrl_bridge.a" "v5/deps/${plat}_${arch}/"
cp "v10/target/wasm32-wasi/release/vrl_bridge.wasm" "v10/deps/"
