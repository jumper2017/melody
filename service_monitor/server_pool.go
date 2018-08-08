package service_monitor

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/jumper2017/melody/errdefine"
)

///////////////////////////////////////////////////////////////////////////////////////////
//对于启动service / 关闭service 如何同步 service_pool 的说明
//在应用层定义一个 service_conn 的map[string] peer ,
//key为 peerName(即sessionName), value为 PeerConnector(即GrpcPeerConnector)
//在watcher 中注入回调函数 defFunc 和 putFunc， 当加入新结点或者删除已有结点时，调用对应回调
//回调函数中操作 service_conn (注意加锁) 新建peer, 或者删除peer

//另外对于如何进行负载均衡，也是值得考虑的问题
///////////////////////////////////////////////////////////////////////////////////////////

type ServerPool struct {
	servers map[string]string
	sync.RWMutex
}

var serverPool ServerPool

func init() {
	serverPool.servers = make(map[string]string)
}

//使用说明：
//假设存在key
// /service/game/mahjong/mahjong_1
// /service/game/mahjong/mahjong_2
// GetServiceInfoWithName("mahjong_1") 获得该service 信息
// GetServiceInfoWithTypeOne("mahjong") 随机获得 mahjong_1 / mahjong_2 的信息
// GetServiceInfoWithTypeAll("mahjong") 获得 mahjong_1 和 mahjong_2 的信息
//需要定义如下接口：
//获取指定service信息
//name 为service name - mahjong_1
func GetServiceInfoWithName(name string) (string, error) {

	if name == "" {
		return "", errdefine.ERR_ETCD_INVALID_PARAM
	}

	serverPool.RLock()
	defer serverPool.RUnlock()

	for k := range serverPool.servers {
		if strings.HasSuffix(k, name) {
			return serverPool.servers[k], nil
		}
	}

	return "", errdefine.ERR_ETCD_SERVICE_NOT_FOUND
}

//获取指定类型service信息， 随机选取一个
func GetServiceInfoWithTypeOne(typ string) (string, string, error) {

	if typ == "" {
		return "", "", errdefine.ERR_ETCD_INVALID_PARAM
	}

	serverPool.RLock()
	defer serverPool.RUnlock()

	tmpKeys := make([]string, 0)

	for k := range serverPool.servers {
		if strings.Contains(k, typ) {
			tmpKeys = append(tmpKeys, k)
		}
	}

	if len(tmpKeys) == 0 {
		return "", "", errdefine.ERR_ETCD_SERVICE_NOT_FOUND
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	choose := r.Int31n(int32(len(tmpKeys)))

	tmpK := strings.Split(tmpKeys[choose], "/")
	name := tmpK[len(tmpK)-1]
	info := serverPool.servers[tmpKeys[choose]]

	return name, info, nil
}

//获取指定类型service信息， 所有
func GetServiceInfoWithTypeAll(typ string) (map[string]string, error) {

	if typ == "" {
		return nil, errdefine.ERR_ETCD_INVALID_PARAM
	}

	serverPool.RLock()
	defer serverPool.RUnlock()

	tmpInfo := make(map[string]string)

	for k, v := range serverPool.servers {
		if strings.Contains(k, typ) {
			tmpInfo[k] = v
		}
	}

	if len(tmpInfo) == 0 {
		return nil, errdefine.ERR_ETCD_SERVICE_NOT_FOUND
	}

	rstInfo := make(map[string]string)
	for k, v := range tmpInfo {
		tmpK := strings.Split(k, "/")
		name := tmpK[len(tmpK)-1]
		info := v
		rstInfo[name] = info
	}

	return rstInfo, nil
}
