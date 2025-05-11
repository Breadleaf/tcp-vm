#!/bin/bash

PROGRAM="stack.asm"

# I wrote this on the road trip back to California with my dad in the UHaul.
# Docker broke cause I had no wifi, thus couldn't test the full system.

base=$(git rev-parse --show-toplevel)
cd "$base"

mkdir -p "$base/build"

echo ">> building all binaries..."
go build -o "$base/build/router" "$base/router"
go build -o "$base/build/server" "$base/server"
go build -o "$base/build/client" "$base/client"

echo ">> starting router..."
./build/router &
router_pid=$!
sleep 1

echo ">> starting server (no args)..."
ROUTER_ID=localhost ./build/server &
server_pid=$!
sleep 1

echo ">> starting client (stack.asm)..."
ROUTER_ID=localhost ./build/client "$base/$PROGRAM" &
client_pid=$!

cleanup() {
    echo
    echo "@> cleaning up..."
    kill "$client_pid" "$server_pid" "$router_pid" 2>/dev/null || true
    wait
    echo "Done."
}
trap cleanup SIGINT SIGTERM

wait
