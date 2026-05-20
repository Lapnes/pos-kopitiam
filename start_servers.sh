#!/bin/bash
# Start backend
./main > /tmp/backend.log 2>&1 &
echo $! > /tmp/backend.pid

# Start Caddy
cd FE/kopitiam-pos-fixed
caddy run --config Caddyfile > /tmp/caddy.log 2>&1 &
echo $! > /tmp/caddy.pid

echo "Servers started!"
