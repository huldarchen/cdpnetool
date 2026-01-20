package regexutil

import (
	"regexp"
	"sync"
)

// Cache 正则表达式缓存
type Cache struct {
	cache sync.Map
}

// New 创建新的正则缓存
func New() *Cache {
	return &Cache{}
}

// Get 返回缓存中的正则或编译后加入缓存
func (c *Cache) Get(p string) (*regexp.Regexp, error) {
	// 1. 尝试从缓存中读取
	if val, ok := c.cache.Load(p); ok {
		return val.(*regexp.Regexp), nil
	}

	// 2. 编译正则
	compiled, err := regexp.Compile(p)
	if err != nil {
		return nil, err
	}

	// 3. 存入缓存
	c.cache.Store(p, compiled)
	return compiled, nil
}
