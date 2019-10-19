package tool

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os/exec"
)

//获取当前网络的最大可用网络数
func GetMaxNetworks() (int, error){
	nic, err := GetNicName()
	if err != nil{
		return 0, err
	}
	it, err := net.InterfaceByName(nic)
	if err != nil{
		return 0, err
	}
	addrs, err := it.Addrs()
	if err != nil{
		return 0, err
	}
	if len(addrs) < 1{
		return 0,fmt.Errorf("请检查网络配置,网卡:%s",nic)
	}
	maxN := 0
	for _,v := range addrs{
		log.Println(v.String())
		ipInfo,_,err := net.ParseCIDR(v.String())
		if err != nil{
			return 0,err
		}
		netMask, size := ipInfo.DefaultMask().Size()
		log.Printf("netmask:%d, size:%d",netMask,size)
		if size != 32{ //只支持IPV4
			continue
		}
		maxN = int(math.Pow(2, float64(size - netMask)))
		break
	}
	if maxN == 0{
		return 0,fmt.Errorf("获取网络信息失败,请检查服务器日志")
	}
	return maxN, nil
}

//获取网卡名称
func GetNicName() (string, error){
	cmd1 := exec.Command( "ip","route","get","8.8.8.8")
	cmd2 := exec.Command( "awk","{ printf $5; exit }")
	nic, err := runTwoCmd(cmd1,cmd2)
	if err != nil{
		return "", err
	}
	if nic == ""{
		return "",errors.New("nic is empty")
	}
	return nic, nil
}

func runTwoCmd(cmd1, cmd2 *exec.Cmd)(string, error){
	reader, writer := io.Pipe()
	var buf bytes.Buffer
	cmd1.Stdout = writer
	cmd2.Stdin = reader
	cmd2.Stdout = &buf
	err := cmd1.Start()
	if err != nil{
		return "", nil
	}
	err = cmd2.Start()
	if err != nil{
		return "", nil
	}
	err = cmd1.Wait()
	if err != nil{
		return "", nil
	}
	err = writer.Close()
	if err != nil{
		return "", nil
	}
	err= cmd2.Wait()
	if err != nil{
		return "", nil
	}
	err = reader.Close()
	if err != nil{
		return "", nil
	}
	return buf.String(), nil
}