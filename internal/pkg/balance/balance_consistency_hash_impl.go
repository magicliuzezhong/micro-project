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
type consistencyHashBalance struct {
	hashRingMap map[string]*hashRing                          //hash环映射
	nodeMap     map[string]map[string]*common.ServiceInstance //快速映射具体实例
}

//
// onceConsistencyHashBalance
// @Description: 一致性hash负载均衡once
//
var onceConsistencyHashBalance sync.Once

//
// balance
// @Description: 具体实例
//
var balance *consistencyHashBalance

//
// NewConsistencyHashBalance
// @Description: 一致性hash负载均衡算法由此调用
// @return ILoadBalance 负载均衡接口
//
func NewConsistencyHashBalance() ILoadBalance {
	onceConsistencyHashBalance.Do(func() {
		balance = &consistencyHashBalance{
			hashRingMap: make(map[string]*hashRing),
			nodeMap:     make(map[string]map[string]*common.ServiceInstance),
		}
	})
	return balance
}

//
// DoBalance
// @Description: 一致性hash负载均衡
// @receiver r ConsistencyHashBalance
// @param instances 实例
// @return *common.ServiceInstance 返回实例
// @return error 错误
//
func (c *consistencyHashBalance) DoBalance(instances []*common.ServiceInstance,
	keys ...string) (*common.ServiceInstance, error) {
	if c.hashRingMap == nil || c.nodeMap == nil {
		return nil, errors.New("参数实例化异常")
	}
	if len(keys) != 2 {
		return nil, errors.New("参数有误")
	}
	var serviceName = keys[0]
	var key = keys[1]
	var insLen = len(instances)
	if _, ok := c.hashRingMap[serviceName]; !ok { //如果不存在那么进行实例化
		c.hashRingMap[serviceName] = newhashRing(make([]string, 0), 32)
	}
	if _, ok := c.nodeMap[serviceName]; !ok { //如果node不存在，创建一个
		c.nodeMap[serviceName] = make(map[string]*common.ServiceInstance)
	}
	if insLen == 0 { //传入的实例为空，清空hashRing的数据
		delete(c.nodeMap, serviceName)
		delete(c.hashRingMap, serviceName)
		return nil, errors.New("传入实例为空")
	}
	originalNoeMap, ok := c.nodeMap[serviceName]     // 第一步 进行传入数据检查
	var addNode = make([]*common.ServiceInstance, 0) //新增的节点
	var delNode = make([]string, 0)                  //删除了的节点
	if !ok || len(originalNoeMap) == 0 {             //该步骤说明数据不存在，全部进行新增即可
		for _, instance := range instances {
			addNode = append(addNode, instance)
		}
	} else { //该步骤说明数据存在
		var newInstances = make(map[string]bool, 0) //映射所有的新节点
		for _, instance := range instances {
			newInstances[instance.GetUrl()] = true
		}
		for _, instance := range instances {
			var url = instance.GetUrl()
			if _, ok := originalNoeMap[url]; !ok { //如果实例不存在，那么说明是新增的节点
				addNode = append(addNode, instance)
			}
		}
		for key, _ := range originalNoeMap {
			if _, ok := newInstances[key]; !ok { //说明已经包含了
				delNode = append(delNode, key)
			}
		}
	}
	var hashRing = c.hashRingMap[serviceName]
	var recordNode = c.nodeMap[serviceName]
	if len(delNode) > 0 {
		for _, node := range delNode {
			hashRing.removeNode(node)
			delete(recordNode, node)
		}
	}
	if len(addNode) > 0 {
		for _, instance := range addNode {
			var url = instance.GetUrl()
			hashRing.addNode(url)
			recordNode[url] = instance
		}
	}
	if hashRing == nil { //没有被注册
		return nil, errors.New("服务实例未映射")
	}
	var serviceInstanceUrl = hashRing.getNode(key)
	if instance, ok := c.nodeMap[serviceName][serviceInstanceUrl]; ok {
		return instance, nil
	}
	return nil, errors.New("实例无法找到")
}

//
// hashRing
// @Description: hash环
//
type hashRing struct {
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
func (hr *hashRing) addNode(masterNode string) {
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

//
// addNodes
// @Description: 批量添加节点，实际还是循环调用的添加节点
// @receiver hr hash环
// @param masterNodes 节点
//
func (hr *hashRing) addNodes(masterNodes []string) {
	if len(masterNodes) > 0 {
		for _, node := range masterNodes {
			hr.addNode(node)
		}
	}
}

//
// removeNode
// @Description: 删除节点
// @receiver hr hash环
// @param masterNode 节点
//
func (hr *hashRing) removeNode(masterNode string) {
	for i := 0; i < hr.replicateCount; i++ {
		key := hr.hashKey(strconv.Itoa(i) + masterNode)
		delete(hr.nodes, key)
		if success, index := hr.getIndexForKey(key); success {
			hr.sortedNodes = append(hr.sortedNodes[:index], hr.sortedNodes[index+1:]...)
		}
	}
}

//
// getNode
// @Description: 获取节点
// @receiver hr hash环
// @param key key
// @return string 获取到的环节点
//
func (hr *hashRing) getNode(key string) string {
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

func (hr *hashRing) getIndexForKey(key uint32) (bool, int) {
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

//
// hashKey
// @Description: 获取hash值
// @receiver hr hash环
// @param key 待hash的key
// @return uint32 hashcode
//
func (hr *hashRing) hashKey(key string) uint32 {
	scratch := []byte(key)
	return crc32.ChecksumIEEE(scratch)
}

func newhashRing(nodes []string, replicateCount int) *hashRing {
	hr := new(hashRing)
	hr.replicateCount = replicateCount
	hr.nodes = make(map[uint32]string)
	hr.sortedNodes = []uint32{}
	hr.addNodes(nodes)
	return hr
}
