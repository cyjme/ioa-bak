#!/bin/sh

/var/etcd-v3.3.12-linux-amd64/etcd &

sleep 3

cd /var/release
/var/release/ioa-httpServer &
/var/release/ioa-proxy
