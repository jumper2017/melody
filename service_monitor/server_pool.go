package service_monitor

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/jumper2017/melody/errdefine"
)

type ServerPool struct {
	servers map[string]string
	sync.Mutex
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

	serverPool.Lock()
	defer serverPool.Unlock()

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

	serverPool.Lock()
	defer serverPool.Unlock()

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

	serverPool.Lock()
	defer serverPool.Unlock()

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
