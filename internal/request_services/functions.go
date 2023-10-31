package request_services

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func convertInterfaceToStruct(m interface{}, s interface{}) error {
	// convert map to json
	jsonString, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// convert json to struct
	err = json.Unmarshal(jsonString, &s)
	if err != nil {
		return err
	}

	return nil
}

type DBAPIInfoK8s struct {
	ServicePort string
	ServiceName string
	Namespace   string
	K8sDNSIp    string
}

func (c *DBAPIInfoK8s) getServiceIP() string {
	dnsName := c.ServiceName + "." + c.Namespace + ".svc.cluster.local"

	// 使用指定的 DNS 服务器地址进行查询
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", c.K8sDNSIp+":53")
		},
	}

	ips, err := resolver.LookupIPAddr(context.Background(), dnsName)
	if err != nil {
		fmt.Printf("无法解析 DNS 名称：%v\n", err)
		os.Exit(1)
	}

	// fmt.Printf("服务 %s 的 IP 地址列表：\n", c.MongodbInfo.ServiceName)

	return ips[0].String()
}
