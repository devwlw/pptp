#!/bin/sh
exec pon vps debug dump logfd 2 nodetach persist "$@" >> /pptp.log 2>&1 &