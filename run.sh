#!/bin/bash 

set -e 

cargo build --release
RUST_BACKTRACE=full go run .