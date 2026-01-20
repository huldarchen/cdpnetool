package pool

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cdpnetool/internal/logger"
)

// Pool 并发工作池，用于限制任务的并发处理数量
type Pool struct {
	sem         chan struct{}
	queue       chan func()
	queueCap    int
	log         logger.Logger
	totalSubmit int64
	totalDrop   int64
	mu          sync.Mutex
	stopMonitor chan struct{}
}

// New 创建工作池，size 为 0 表示无限制，queueCap 为任务排队队列容量
func New(size int, queueCap int) *Pool {
	if size <= 0 {
		return &Pool{}
	}

	if queueCap <= 0 {
		queueCap = size * 8
	}

	return &Pool{
		sem:      make(chan struct{}, size),
		queue:    make(chan func(), queueCap),
		queueCap: queueCap,
	}
}

// SetLogger 设置日志记录器
func (p *Pool) SetLogger(l logger.Logger) {
	p.log = l
}

// Start 启动工作池，创建固定数量 of worker 协程
func (p *Pool) Start(ctx context.Context) {
	if p.sem == nil {
		return
	}
	for i := 0; i < cap(p.sem); i++ {
		go p.worker(ctx)
	}
	p.stopMonitor = make(chan struct{})
	go p.monitor(ctx)
}

// Stop 停止监控协程
func (p *Pool) Stop() {
	if p.stopMonitor != nil {
		close(p.stopMonitor)
	}
}

// monitor 定期输出工作池状态监控日志
func (p *Pool) monitor(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopMonitor:
			return
		case <-ticker.C:
			qLen, qCap, submit, drop := p.Stats()
			if p.log != nil && submit > 0 {
				usage := float64(qLen) / float64(qCap) * 100
				dropRate := float64(drop) / float64(submit) * 100
				p.log.Info("工作池状态监控", "queueLen", qLen, "queueCap", qCap, "usage", fmt.Sprintf("%.1f%%", usage), "totalSubmit", submit, "totalDrop", drop, "dropRate", fmt.Sprintf("%.2f%%", dropRate))
			}
		}
	}
}

// worker 工作协程，从队列中取任务并执行
func (p *Pool) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case fn := <-p.queue:
			if fn != nil {
				fn()
			}
		}
	}
}

// Submit 提交任务到工作池，返回是否成功入队
func (p *Pool) Submit(fn func()) bool {
	if p.sem == nil {
		go fn()
		return true
	}
	p.mu.Lock()
	p.totalSubmit++
	p.mu.Unlock()
	select {
	case p.queue <- fn:
		return true
	default:
		p.mu.Lock()
		p.totalDrop++
		drop := p.totalDrop
		submit := p.totalSubmit
		p.mu.Unlock()
		if p.log != nil {
			p.log.Warn("工作池队列已满，任务被丢弃", "queueCap", p.queueCap, "totalSubmit", submit, "totalDrop", drop)
		}
		return false
	}
}

// Stats 返回工作池统计信息
func (p *Pool) Stats() (queueLen, queueCap, totalSubmit, totalDrop int64) {
	if p.sem == nil {
		return 0, 0, 0, 0
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return int64(len(p.queue)), int64(p.queueCap), p.totalSubmit, p.totalDrop
}

// GetQueueCap 返回队列容量
func (p *Pool) GetQueueCap() int {
	return p.queueCap
}

// IsEnabled 检查工作池是否已启用并发限制
func (p *Pool) IsEnabled() bool {
	return p.sem != nil
}
