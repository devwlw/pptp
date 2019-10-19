#!/bin/bash
yum install -y net-tools
yum install -y wget
wget -O /etc/yum.repos.d/docker-ce.repo https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
yum install -y docker-ce-18.06.0.ce-3.el7
systemctl daemon-reload
systemctl start docker
systemctl enable docker
#启用docker restful api
sed -i 's/-H fd:\/\//-H fd:\/\/ -H tcp:\/\/0.0.0.0:2375/g' /lib/systemd/system/docker.service
systemctl daemon-reload
service docker restart
## todo 这里要放在文件里面,防止重启服务后失效
sudo modprobe nf_conntrack_pptp
sudo modprobe nf_conntrack_proto_gre
nic=$(ip route get 8.8.8.8 | awk '{printf $5}')
ip link set  $nic  promisc on

yum install -y epel-release
wget https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
yum install -y ./epel-release-latest-7.noarch.rpm
yum -y update
yum install -y git
yum -y install python-pip
pip install --upgrade pip
pip install --upgrade netifaces
pip install --upgrade requests
pip install --ignore-installed  netaddr==0.7.19
#
wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.10.3.linux-amd64.tar.gz
mkdir -p /root/go
mkdir -p /root/go/src
mkdir -p /root/go/bin
mkdir -p /root/go/pkg
echo "export PATH=$PATH:/usr/local/go/bin" >> /root/.bash_profile
echo "export GOPATH=/root/go" >> /root/.bash_profile
echo "export PATH=$PATH:$GOPATH/bin" >> /root/.bash_profile
#go get -u github.com/golang/dep/cmd/dep
source /root/.bash_profile
cd ../deliver
##todo  改成docker pull images
docker build -t pptp .
