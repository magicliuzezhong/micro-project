//
// Package redis
// @Description：redis工具
// @Author：liuzezhong 2021/7/15 6:19 下午
// @Company cloud-ark.com
//
package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"micro-project/internal/pkg/util"
	"sync"
)

type IRedisClient interface {
	GetRedisCluster() *redis.ClusterClient
	GetRedisClient() *redis.Client
	Reset()
}

type currentRedisClient struct {
	clusterClient *redis.ClusterClient
	client        *redis.Client
	changeLock    sync.RWMutex
}

var currentRedisClientOnce sync.Once

var currentRedisClientInstance *currentRedisClient

func newCurrentRedisClient() IRedisClient {
	currentRedisClientOnce.Do(func() {
		currentRedisClientInstance = &currentRedisClient{}
		currentRedisClientInstance.initRedis() //初始化redis
	})
	return currentRedisClientInstance
}

func (r *currentRedisClient) GetRedisCluster() *redis.ClusterClient {
	r.changeLock.RLock()
	defer r.changeLock.RUnlock()
	return r.clusterClient
}

func (r *currentRedisClient) GetRedisClient() *redis.Client {
	r.changeLock.RLock()
	defer r.changeLock.RUnlock()
	return r.client
}

func (r *currentRedisClient) Reset() {
	r.changeLock.Lock()
	defer r.changeLock.Unlock()
	r.initRedis()
}

func (r *currentRedisClient) initRedis() {
	r.initRedisAlone()
	r.initRedisCluster()
}

func (r *currentRedisClient) initRedisAlone() {
	var redisAlone = util.GetApplication().DBInfo.Redis.RedisAlone
	if !redisAlone.Enabled || redisAlone.Url == "" { //单机版无法使用启用
		return
	}
	var client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisAlone.Url, redisAlone.Port),
		Password: redisAlone.Password,
	})
	var cmd = client.Ping()
	var err = cmd.Err()
	if err == nil {
		r.client = client
	} else {
		fmt.Println(err.Error())
	}
}

func (r *currentRedisClient) initRedisCluster() {
	var redisCluster = util.GetApplication().DBInfo.Redis.RedisCluster
	var redisClusters = util.GetApplication().DBInfo.Redis.RedisCluster.RedisClusterInfo
	if redisClusters == nil || !redisCluster.Enabled || len(redisClusters) == 0 {
		return
	}
	var redisClusterPassword = redisCluster.Password
	var addrArr = make([]string, 0)
	for _, value := range redisClusters {
		addrArr = append(addrArr, fmt.Sprintf("%s:%d", value.Url, value.Port))
	}
	var client = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrArr,              //set redis cluster url
		Password: redisClusterPassword, //set password
	})
	var cmd = client.Ping()
	var err = cmd.Err()
	if err == nil {
		r.clusterClient = client
	}
}

func GetRedisClient() *redis.Client {
	return newCurrentRedisClient().GetRedisClient()
}

func GetRedisCluster() *redis.ClusterClient {
	return newCurrentRedisClient().GetRedisCluster()
}

func Reset() {
	newCurrentRedisClient().Reset()
}
