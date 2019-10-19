#!/bin/sh
while true; do
  ip=$(curl icanhazip.com)
  curl -XPOST ${HOSTIP}:9100/docker/container/pptpip?ip=$ip
  sleep 30;
done

