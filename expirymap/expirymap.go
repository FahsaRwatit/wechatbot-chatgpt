package expirymap

import (
	"sync"
	"time"
)

type ExpiryMap struct {
	// 存储键值对key-value
	m map[string]string
	// 读写锁
	mutex sync.RWMutex
	// 存储每个键值对的过期时间
	expiryMap map[string]time.Time
}

func New() ExpiryMap {
	return ExpiryMap{
		m:         make(map[string]string),
		expiryMap: make(map[string]time.Time),
	}
}

func (em *ExpiryMap) Set(key string, value string, expiry time.Duration) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	em.m[key] = value
	em.expiryMap[key] = time.Now().Add(expiry)
}

func (em *ExpiryMap) Get(key string) (string, bool) {
	em.mutex.Lock()
	defer em.mutex.RUnlock()
	if value, ok := em.m[key]; ok {
		expiry := em.expiryMap[key]
		if time.Now().Before(expiry) {
			return value, true
		}
		delete(em.m, key)
		delete(em.expiryMap, key)
	}
	return "", false
}
func (em *ExpiryMap) Delete(key string) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	delete(em.m, key)
	delete(em.expiryMap, key)
}
