#!/bin/sh
mknod /dev/ppp c 108 0
cat > /etc/ppp/peers/vps <<_EOF_
pty "pptp ${SERVER} --nolaunchpppd"
name "${USERNAME}"
password "${PASSWORD}"
remotename PPTP
file /etc/ppp/options.pptp
ipparam vps
_EOF_

cat > /etc/ppp/ip-up <<"_EOF_"
#!/bin/sh
ip route add 0.0.0.0/1 dev $1
ip route add 128.0.0.0/1 dev $1
_EOF_

cat > /etc/ppp/ip-down <<"_EOF_"
#!/bin/sh
ip route del 0.0.0.0/1 dev $1
ip route del 128.0.0.0/1 dev $1
_EOF_
/tick.sh > tick.log 2>&1 &
java -jar iem-server-163.jar
#exec pon vps debug dump logfd 2 nodetach persist "$@"
