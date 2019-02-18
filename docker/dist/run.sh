#!/bin/sh

/var/etcd-v3.3.12-linux-amd64/etcd --data-dir="/data/etcd" &

sleep 3

cd /var/dashboard
http-server -p 9993 &

cd /var/release
/var/release/ioa-httpServer &
/var/release/ioa-proxy
