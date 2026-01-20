package regexutil_test

import (
	"sync"
	"testing"

	"cdpnetool/internal/regexutil"
)

// TestCache_Hit 验证缓存命中逻辑：相同的 pattern 应该返回同一个对象指针
func TestCache_Hit(t *testing.T) {
	c := regexutil.New()
	pattern := `^https?://.*`

	// 第一次获取
	re1, err := c.Get(pattern)
	if err != nil {
		t.Fatalf("第一次获取失败: %v", err)
	}

	// 第二次获取
	re2, err := c.Get(pattern)
	if err != nil {
		t.Fatalf("第二次获取失败: %v", err)
	}

	// 验证指针地址是否一致
	if re1 != re2 {
		t.Errorf("缓存失效：两次获取相同 pattern 返回了不同的对象指针")
	}
}

// TestCache_InvalidRegex 验证非法正则表达式的处理
func TestCache_InvalidRegex(t *testing.T) {
	c := regexutil.New()
	invalidPattern := `[` // 非法正则

	_, err := c.Get(invalidPattern)
	if err == nil {
		t.Error("期望非法正则返回错误，但实际未返回")
	}
}

// TestCache_Concurrency 验证并发安全性
func TestCache_Concurrency(t *testing.T) {
	c := regexutil.New()
	pattern := `[a-z]+`

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// 启动 100 个协程同时获取同一个正则
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_, err := c.Get(pattern)
			if err != nil {
				t.Errorf("并发获取失败: %v", err)
			}
		}()
	}

	wg.Wait()
}

// TestCache_MultiplePatterns 验证多个不同正则的缓存
func TestCache_MultiplePatterns(t *testing.T) {
	c := regexutil.New()
	patterns := []string{`^abc`, `\d+`, `.*\.js$`}

	for _, p := range patterns {
		re1, _ := c.Get(p)
		re2, _ := c.Get(p)
		if re1 != re2 {
			t.Errorf("Pattern %s 缓存失效", p)
		}
	}
}
