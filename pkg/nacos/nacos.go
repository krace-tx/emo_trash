package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/zeromicro/go-zero/core/logx"
	"net"
	"strings"
	"time"
)

type NacosConf struct {
	Server struct {
		Host  string
		Port  uint64
		Group string
		Name  string
	}
	Discovery struct {
		ServerAddr string
		ServerPort uint64
		Namespace  string
		Group      string
	}
	Config struct {
		ServerAddr string
		ServerPort uint64
		Namespace  string
		Group      string
		DataId     string
	}
}

func Init(n *NacosConf) (config string, err error) {
	_ = InitNacosNamingClient(n)

	// 初始化 Nacos 配置客户端
	configClient, err := InitNacosConfigClient(n)
	if err != nil {
		logx.Errorf("配置客户端初始化失败: %v", err)
		return "", err
	}

	// 从 Nacos 获取动态配置，并更新到配置结构体中
	config, err = configClient.GetConfig(vo.ConfigParam{
		DataId: n.Config.DataId,
		Group:  n.Config.Group,
	})
	if err != nil {
		logx.Errorf("从 Nacos 获取配置失败: %v", err)
		return "", err
	}

	return config, nil
}

// InitNacosConfigClient 初始化 Nacos 配置客户端
func InitNacosConfigClient(conf *NacosConf) (config_client.IConfigClient, error) {
	// Nacos 服务配置
	serverConfig := []constant.ServerConfig{
		{
			IpAddr: conf.Config.ServerAddr,
			Port:   conf.Config.ServerPort,
		},
	}

	// 客户端配置
	clientConfig := constant.ClientConfig{
		NamespaceId:         conf.Config.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacosx/log",   // 可以根据需求进行调整
		CacheDir:            "/tmp/nacosx/cache", // 可以根据需求进行调整
	}

	// 创建配置客户端
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfig,
		},
	)

	if err != nil {
		logx.Errorf("初始化 Nacos 配置客户端失败: %v", err)
		return nil, err
	}

	return client, nil
}

const maxRetries = 5              // 最大重试次数
const retryInterval = time.Second // 重试间隔

// 初始化服务发现客户端并支持断线重连
func InitNacosNamingClient(conf *NacosConf) naming_client.INamingClient {
	var client naming_client.INamingClient
	var err error

	// Nacos 服务配置
	serverConfig := []constant.ServerConfig{
		{
			IpAddr: conf.Discovery.ServerAddr,
			Port:   conf.Discovery.ServerPort,
		},
	}

	// 客户端配置
	clientConfig := constant.ClientConfig{
		NamespaceId:         conf.Discovery.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",   // 根据需求调整
		CacheDir:            "/tmp/nacos/cache", // 根据需求调整
	}

	for i := 0; i < maxRetries; i++ {
		client, err = clients.NewNamingClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfig,
			},
		)
		if err == nil {
			break
		}
		logx.Errorf("尝试创建服务发现客户端失败 (%d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}
	if err != nil {
		logx.Errorf("初始化服务发现客户端失败: %v", err)
	}

	// 注册实例并启用自动重连机制
	go func() {
		for {
			err := registerService(client, conf)
			if err != nil {
				logx.Errorf("%s 服务注册失败，尝试重连: %v", conf.Server.Name, err)
				time.Sleep(retryInterval)
				continue
			}
			logx.Infof(" %s 服务注册成功", conf.Server.Name)
			return
		}
	}()
	return client
}

// registerService 注册服务实例
func registerService(client naming_client.INamingClient, conf *NacosConf) error {
	instanceParam := vo.RegisterInstanceParam{
		Ip:          getServiceIP(conf.Server.Host),
		Port:        conf.Server.Port,
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Metadata:    nil,
		ServiceName: conf.Server.Name,
		GroupName:   conf.Server.Group,
		Ephemeral:   true,
	}

	_, err := client.RegisterInstance(instanceParam)
	return err
}

const (
	localHost = "127.0.0.1;localhost;0.0.0.0"
)

// getServiceIP 获取服务 IP
func getServiceIP(host string) string {
	if strings.Contains(localHost, host) {
		ip, err := getActiveNetworkIP()
		if err != nil {
			logx.Errorf("无法读取 IP 地址: %v", err)
		}
		return ip
	}
	return host
}

// 获取当前活动的上网网卡的 IPv4 地址
func getActiveNetworkIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}

			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						return ipNet.IP.String(), nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("没有找到活动的上网网卡 IP 地址")
}
