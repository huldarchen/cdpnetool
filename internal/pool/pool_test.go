package pool_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"cdpnetool/internal/pool"
)

// TestPool_Basic 验证任务能正常执行
func TestPool_Basic(t *testing.T) {
	p := pool.New(2, 50) // 调大队列确保不丢弃
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p.Start(ctx)

	var count int32
	wg := sync.WaitGroup{}
	numTasks := 20

	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		ok := p.Submit(func() {
			atomic.AddInt32(&count, 1)
			wg.Done()
		})
		if !ok {
			t.Errorf("任务 %d 提交失败", i)
			wg.Done()
		}
	}

	wg.Wait()
	if atomic.LoadInt32(&count) != int32(numTasks) {
		t.Errorf("期望执行 %d 个任务, 实际执行 %d", numTasks, count)
	}
}

// TestPool_ConcurrencyLimit 验证并发数限制
func TestPool_ConcurrencyLimit(t *testing.T) {
	size := 3
	p := pool.New(size, 20)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p.Start(ctx)

	var activeWorkers int32
	var maxActive int32
	wg := sync.WaitGroup{}
	numTasks := 10
	block := make(chan struct{})

	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		p.Submit(func() {
			defer wg.Done()
			current := atomic.AddInt32(&activeWorkers, 1)
			// 更新历史最高并发数
			for {
				prevMax := atomic.LoadInt32(&maxActive)
				if current <= prevMax || atomic.CompareAndSwapInt32(&maxActive, prevMax, current) {
					break
				}
			}
			<-block // 阻塞任务
			atomic.AddInt32(&activeWorkers, -1)
		})
	}

	// 给一点时间让 worker 运行
	time.Sleep(100 * time.Millisecond)

	actualMax := atomic.LoadInt32(&maxActive)
	if actualMax != int32(size) {
		t.Errorf("期望最大并发数为 %d, 实际为 %d", size, actualMax)
	}

	close(block) // 释放所有任务
	wg.Wait()
}

// TestPool_Drop 验证队列满时的丢弃策略
func TestPool_Drop(t *testing.T) {
	// size=1 (1个正在跑), queueCap=1 (1个在排队), 总容量=2
	p := pool.New(1, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p.Start(ctx)

	block := make(chan struct{})
	defer close(block)

	// 1. 提交任务 A (占用 worker)
	if !p.Submit(func() { <-block }) {
		t.Fatal("任务 A 提交失败")
	}

	// 给一点点时间让 worker 把任务 A 从队列领走
	for i := 0; i < 50; i++ {
		qLen, _, _, _ := p.Stats()
		if qLen == 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// 2. 提交任务 B (进入队列)
	if !p.Submit(func() { <-block }) {
		t.Fatal("任务 B 提交失败")
	}

	// 3. 提交任务 C (队列已满，应该被丢弃)
	if p.Submit(func() { <-block }) {
		t.Error("任务 C 应该提交失败，但成功了")
	}

	_, _, submit, drop := p.Stats()
	if submit != 3 {
		t.Errorf("期望提交计数为 3, 实际为 %d", submit)
	}
	if drop != 1 {
		t.Errorf("期望丢弃计数为 1, 实际为 %d", drop)
	}
}

// TestPool_Unbounded 验证 size=0 情况下的无限制模式
func TestPool_Unbounded(t *testing.T) {
	p := pool.New(0, 0)
	if p.IsEnabled() {
		t.Error("size=0 时 IsEnabled 应该返回 false")
	}

	var count int32
	wg := sync.WaitGroup{}
	numTasks := 50

	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		p.Submit(func() {
			atomic.AddInt32(&count, 1)
			wg.Done()
		})
	}

	wg.Wait()
	if atomic.LoadInt32(&count) != int32(numTasks) {
		t.Errorf("期望执行 %d 个任务, 实际执行 %d", numTasks, count)
	}
}

// TestPool_ContextCancel 验证 worker 随 context 取消而退出
func TestPool_ContextCancel(t *testing.T) {
	p := pool.New(2, 10)
	ctx, cancel := context.WithCancel(context.Background())
	p.Start(ctx)

	cancel() // 取消 context

	// 理论上此时 worker 应该退出了
	// 我们通过观察 Stats 无法直接确认协程退出，
	// 但我们可以提交一个永远不会被处理的任务来验证

	// 给点时间让 context 传播
	time.Sleep(50 * time.Millisecond)

	taskRan := make(chan struct{})
	p.Submit(func() {
		close(taskRan)
	})

	select {
	case <-taskRan:
		t.Error("Worker 应该在 context 取消后停止处理任务")
	case <-time.After(100 * time.Millisecond):
		// 正常：任务没有被处理
	}
}
