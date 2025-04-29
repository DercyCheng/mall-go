package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdClient Etcd客户端
type EtcdClient struct {
	client *clientv3.Client
	prefix string
}

// NewEtcdClient 创建Etcd客户端
func NewEtcdClient(endpoints []string, prefix string) (*EtcdClient, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Printf("初始化Etcd客户端失败: %v", err)
		return nil, err
	}

	return &EtcdClient{
		client: cli,
		prefix: prefix,
	}, nil
}

// Close 关闭Etcd客户端
func (c *EtcdClient) Close() error {
	return c.client.Close()
}

// LoadConfig 从配置中心加载配置
func (c *EtcdClient) LoadConfig(key string, configStruct interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fullKey := fmt.Sprintf("%s/%s", c.prefix, key)
	resp, err := c.client.Get(ctx, fullKey)
	if err != nil {
		log.Printf("从Etcd获取配置失败: %v", err)
		return err
	}

	if len(resp.Kvs) == 0 {
		return fmt.Errorf("配置不存在: %s", fullKey)
	}

	value := resp.Kvs[0].Value
	err = json.Unmarshal(value, configStruct)
	if err != nil {
		log.Printf("解析配置失败: %v", err)
		return err
	}

	return nil
}

// SaveConfig 保存配置到配置中心
func (c *EtcdClient) SaveConfig(key string, configStruct interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	value, err := json.Marshal(configStruct)
	if err != nil {
		log.Printf("序列化配置失败: %v", err)
		return err
	}

	fullKey := fmt.Sprintf("%s/%s", c.prefix, key)
	_, err = c.client.Put(ctx, fullKey, string(value))
	if err != nil {
		log.Printf("保存配置到Etcd失败: %v", err)
		return err
	}

	return nil
}

// WatchConfig 监听配置变更
func (c *EtcdClient) WatchConfig(key string, onChange func([]byte)) {
	fullKey := fmt.Sprintf("%s/%s", c.prefix, key)
	watchChan := c.client.Watch(context.Background(), fullKey)

	// 在新的goroutine中处理监听事件
	go func() {
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				if event.Type == clientv3.EventTypePut {
					log.Printf("配置发生变更: %s", fullKey)
					onChange(event.Kv.Value)
				}
			}
		}
	}()
}
