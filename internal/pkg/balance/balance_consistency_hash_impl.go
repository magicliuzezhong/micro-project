//
// Package balance
// @Description：一致性hash负载均衡实现
// @Author：liuzezhong 2021/6/30 7:18 下午
// @Company cloud-ark.com
//
package balance

import (
	"errors"
	"hash/crc32"
	"micro-project/internal/pkg/common"
	"sort"
	"strconv"
	"sync"
)

//
// ConsistencyHashBalance
// @Description: 一致性hash负载均衡实现
//
type ConsistencyHashBalance struct {
	hashRingMap map[string]*HashRing                          //hash环映射
	nodeMap     map[string]map[string]*common.ServiceInstance //快速映射具体实例
}

var onceConsistencyHashBalance sync.Once

var balance *ConsistencyHashBalance

//
// NewConsistencyHashBalance
// @Description: 一致性hash负载均衡算法由此调用
// @param serviceName 服务名称
// @param instances 实例
// @return ILoadBalance 负载均衡接口
//
func NewConsistencyHashBalance(serviceName string, instances []*common.ServiceInstance) ILoadBalance {
	onceConsistencyHashBalance.Do(func() {
		balance = &ConsistencyHashBalance{
			hashRingMap: make(map[string]*HashRing),
			nodeMap:     make(map[string]map[string]*common.ServiceInstance),
		}
	})
	_ = balance.RegisterServiceInstance(serviceName, instances)
	return balance
}

func (c *ConsistencyHashBalance) RegisterServiceInstance(serviceName string, instances []*common.ServiceInstance) error {
	if _, ok := c.hashRingMap[serviceName]; ok { //已经注册过了
		return errors.New("该实例已注册")
	}
	var instanceLen = len(instances)
	if instanceLen == 0 {
		return errors.New("服务实例为空")
	}
	var nodes = make([]string, instanceLen)
	var serviceInstanceMap = make(map[string]*common.ServiceInstance)
	for i := 0; i < instanceLen; i++ {
		var url = instances[i].GetUrl()
		nodes[i] = url
		serviceInstanceMap[url] = instances[i]
	}
	c.nodeMap[serviceName] = serviceInstanceMap
	c.hashRingMap[serviceName] = newHashRing(nodes, 32) //实例化节点，并且使用32个虚拟节点
	return nil
}

//
// DoBalance
// @Description: 一致性hash负载均衡
// @receiver r ConsistencyHashBalance
// @param instances 实例
// @return *common.ServiceInstance 返回实例
// @return error 错误
//
func (c *ConsistencyHashBalance) DoBalance(instances []*common.ServiceInstance,
	keys ...string) (*common.ServiceInstance, error) {
	if len(keys) != 2 {
		return nil, errors.New("参数有误")
	}
	var serviceName = keys[0]
	var key = keys[1]
	if _, ok := c.hashRingMap[serviceName]; !ok { //已经注册过了
		return nil, errors.New("服务实例未映射")
	}
	var instanceLen = len(instances)
	if instanceLen == 0 { //清空之前个实例映射，解除关系引用等待垃圾回收
		delete(c.hashRingMap, serviceName)
		delete(c.nodeMap, serviceName)
		return nil, errors.New("服务实例为空")
	}
	var hashRing = c.hashRingMap[serviceName]
	var serviceInstanceUrl = hashRing.getNode(key)
	if instance, ok := c.nodeMap[serviceName][serviceInstanceUrl]; ok {
		return instance, nil
	}
	return nil, errors.New("实例无法找到")
}

type HashRing struct {
	replicateCount int // 虚拟节点.
	nodes          map[uint32]string
	sortedNodes    []uint32
}

//
// addNode
// @Description: 作用：在哈希环上添加单个服务器节点（包含虚拟节点）的方法
// @receiver hr hash环
// @param masterNode 入参：服务器地址
//
func (hr *HashRing) addNode(masterNode string) {
	// 为每台服务器生成数量为 replicateCount-1 个虚拟节点
	// 并将其与服务器的实际节点一同添加到哈希环中
	for i := 0; i < hr.replicateCount; i++ {
		// 获取节点的哈希值，其中节点的字符串为 i+address
		key := hr.hashKey(strconv.Itoa(i) + masterNode)
		// 设置该节点所对应的服务器（建立节点与服务器地址的映射）
		hr.nodes[key] = masterNode
		// 将节点的哈希值添加到哈希环中
		hr.sortedNodes = append(hr.sortedNodes, key)
	}
	// 按照值从大到小的排序函数
	sort.Slice(hr.sortedNodes, func(i, j int) bool {
		return hr.sortedNodes[i] < hr.sortedNodes[j]
	})
}

func (hr *HashRing) addNodes(masterNodes []string) {
	if len(masterNodes) > 0 {
		for _, node := range masterNodes {
			hr.addNode(node)
		}
	}
}

func (hr *HashRing) removeNode(masterNode string) {

	for i := 0; i < hr.replicateCount; i++ {
		key := hr.hashKey(strconv.Itoa(i) + masterNode)
		delete(hr.nodes, key)

		if success, index := hr.getIndexForKey(key); success {
			hr.sortedNodes = append(hr.sortedNodes[:index], hr.sortedNodes[index+1:]...)
		}
	}
}

func (hr *HashRing) getNode(key string) string {
	if len(hr.nodes) == 0 {
		return ""
	}
	hashKey := hr.hashKey(key)
	nodes := hr.sortedNodes
	masterNode := hr.nodes[nodes[0]]

	for _, node := range nodes {
		if hashKey < node {
			masterNode = hr.nodes[node]
			break
		}
	}

	return masterNode
}

func (hr *HashRing) getIndexForKey(key uint32) (bool, int) {

	index := -1
	success := false

	for i, v := range hr.sortedNodes {
		if v == key {
			index = i
			success = true
			break
		}
	}

	return success, index
}

func (hr *HashRing) hashKey(key string) uint32 {
	scratch := []byte(key)
	return crc32.ChecksumIEEE(scratch)
}

func newHashRing(nodes []string, replicateCount int) *HashRing {
	hr := new(HashRing)
	hr.replicateCount = replicateCount
	hr.nodes = make(map[uint32]string)
	hr.sortedNodes = []uint32{}
	hr.addNodes(nodes)
	return hr
}
