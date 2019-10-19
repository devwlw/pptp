# -*- coding:utf8 -*-
import sys, getopt
import netifaces
import requests
import json
import math
import os
from netaddr import *


def main(argv):
    admin = ''
    try:
        opts, args = getopt.getopt(argv,"a:",["admin=",])
    except getopt.GetoptError:
        print("--admin -a")
        sys.exit(2)
    for opt, arg in opts:
        if opt in("-a", "--admin"):
            admin = arg
    cluster_info = check_cluster_status(admin)
    net_info = network_info()
    print(net_info)
    print(cluster_info)
    netmask = IPAddress(net_info["netmask"]).netmask_bits()
    ip_list = [ip for ip in IPNetwork('%s/%d' % (net_info["gateway"], netmask))]
    subnet = '%s/%d' % (net_info["gateway"], netmask)

    node_num = cluster_info["nodeNum"]
    max_instance = cluster_info["maxInstance"]
    max_machine = cluster_info["maxMachine"]  #
    start_ip_index = node_num * max_instance
    cidr_netmask = netmask + int(math.log(max_machine, 2))
    ip_range = '%s/%d' % (ip_list[start_ip_index], cidr_netmask)
    aux = ip_list[start_ip_index + 2]
    #创建docker网络
    cmd_str = "docker network create -d macvlan -o parent=%s --subnet %s --gateway %s --ip-range %s --aux-address " \
              "'host=%s' macvlan" % (net_info["netname"], subnet, net_info["gateway"], ip_range, aux)
    print(cmd_str)
    cmd_re = os.popen(cmd_str).read()
    print(cmd_re)
    #创建docker网桥
    #todo 这里应该单独提出来,每次重启网络,或者机器都会失效
    #ip link add mac-docker link ens33 type macvlan  mode bridge
    cmd = "ip link add mac-docker link %s type macvlan  mode bridge" % (net_info["netname"])
    print(cmd)
    cmd_re = os.popen(cmd)
    print(cmd_re)
    #ip addr add 192.168.101.200/32 dev mac-docker
    cmd = "ip addr add %s/32 dev mac-docker" % aux
    print(cmd)
    cmd_re = os.popen(cmd)
    print(cmd_re)
    cmd = "ip link set mac-docker up"
    print(cmd)
    cmd_re = os.popen("ip link set mac-docker up")
    print(cmd_re)
    cmd = "ip route add %s dev mac-docker" % ip_range
    print(cmd)
    cmd_re = os.popen(cmd)
    print(cmd_re)
    urlStr = "%s/docker/node/register/%s?aux=%s&ipRange=%s&subnet=%s&nic=%s" % (admin, net_info["addr"], aux, ip_range, subnet, net_info["netname"])
    print("url:", urlStr)
    re = requests.post("%s/docker/node/register/%s?aux=%s&ipRange=%s&subnet=%s&nic=%s" % (admin, net_info["addr"], aux, ip_range, subnet, net_info["netname"]))
    if re.status_code is not 200:
        print(re.content)
        exit(1)
    re = json.loads(re.content)
    if re["success"] is False:
        print(re["errMsg"])
        exit(1)
    else:
        print("节点注册成功")
    #todo 创建docker network
    #todo 向服务器注册
    #todo 为了最大限度利用ip地址,maxMachine必须为2的N次方,maxInstance必须根据ip计算, 找go的cidr程序
    #todo 实验bridge情况
    #todo 创建go的分发程序
    #todo 简单的admin
    #todo 邮件重写

def network_info():
    gw = netifaces.gateways()["default"][netifaces.AF_INET]
    gateway = gw[0]
    netname = gw[1]
    addrs = netifaces.ifaddresses(netname)[netifaces.AF_INET][0]
    addr = addrs["addr"]
    netmask = addrs["netmask"]
    return {
        "gateway": gateway,
        "netname": netname,
        "addr": addr,
        "netmask": netmask,
    }


def check_cluster_status(admin_addr):
    node = requests.get(admin_addr+"/docker/node") #设置超时时间,错误处理
    if node.status_code is not 200:
        print(node.content)
        exit(1)
    re = json.loads(node.content)
    if re["success"] is False:
        print(re["errMsg"])
        exit(1)
    nodeNum = 0
    if re["data"] is not None:
        node_info = re["data"]["nodes"]
        if node_info is not None:
            nodeNum = len(node_info)
    deploy = requests.get(admin_addr+"/docker/deploy")
    if deploy.status_code is not 200:
        print(deploy.content)
        exit(1)

    re = json.loads(deploy.content)
    if re["success"] is False:
        print(re["errMsg"])
        exit(1)
    re = re["data"]
    maxMachine = re["maxMachine"]
    maxInstance = re["maxInstance"]
    if nodeNum >= maxMachine:
        print('当前节点数%d,允许的最大节点数%d' % (nodeNum, maxMachine))
        exit(1)
    return {
        "nodeNum": nodeNum,
        "maxMachine": maxMachine,
        "maxInstance": maxInstance,
    }

#todo 判断已有服务器的ip信息,若冲突则提示

if __name__ == "__main__":
    main(sys.argv[1:])



#info = network_info()

#print(IPAddress("255.255.255.0").netmask_bits())
#c = IPAddress("255.255.0.0").netmask_bits()
#print(list(IPNetwork("192.168.101.0/24").subnet()))
#ip_list = [ip for ip in IPNetwork('192.168.101.0/24')]
#print(ip_list)
#print(len(ip_list))

# 获取admin已注册机器数和每台机器的最大docker实例数