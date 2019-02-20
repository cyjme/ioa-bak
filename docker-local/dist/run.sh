#!/bin/sh

/var/etcd-v3.3.12-linux-amd64/etcd --data-dir="/data/etcd" &

sleep 3

cd /dashboard
http-server -p 9993 &

ioa-httpServer --config $CONFIG_LOCATION &
ioa-proxy --config $CONFIG_LOCATION
