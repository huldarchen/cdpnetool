package inspector

import (
	"context"
	"sync"
	"time"

	"cdpnetool/internal/executor"
	"cdpnetool/internal/logger"
	"cdpnetool/pkg/domain"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/fetch"
)

// Inspector 流量观察者，负责全量网络请求的镜像采集与状态聚合
type Inspector struct {
	mu      sync.RWMutex
	enabled bool
	events  chan domain.NetworkEvent
	log     logger.Logger

	// 在途事务池: map[RequestID]*domain.NetworkEvent
	sessions sync.Map
}

// Config 观察者配置
type Config struct {
	Events chan domain.NetworkEvent
	Logger logger.Logger
}

// New 创建流量观察者
func New(cfg Config) *Inspector {
	if cfg.Logger == nil {
		cfg.Logger = logger.NewNop()
	}
	i := &Inspector{
		events: cfg.Events,
		log:    cfg.Logger,
	}
	go i.cleanupLoop()
	return i
}

// SetEnabled 动态开关
func (i *Inspector) SetEnabled(enabled bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.enabled = enabled
	if !enabled {
		// 关闭时清空内存池
		i.sessions.Range(func(key, value any) bool {
			i.sessions.Delete(key)
			return true
		})
	}
}

// IsEnabled 是否启用
func (i *Inspector) IsEnabled() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.enabled
}

// RecordRequest 镜像采集请求阶段
func (i *Inspector) RecordRequest(
	targetID domain.TargetID,
	ev *fetch.RequestPausedReply,
) {
	if !i.IsEnabled() {
		return
	}

	// 采集请求快照（只读模式）
	exec := executor.New(i.log, ev, executor.Options{})
	reqInfo := exec.CaptureRequestSnapshot()

	event := domain.NetworkEvent{
		ID:        string(ev.RequestID),
		Target:    targetID,
		Timestamp: time.Now().UnixMilli(),
		Request:   reqInfo,
		IsMatched: false, // Inspector 记录的默认视为未命中（由 Handler 标记命中的除外）
	}

	// 入池暂存，等待响应
	i.sessions.Store(ev.RequestID, &event)

	// 立即向前端推送“请求中”状态
	i.pushEvent(event)
}

// RecordResponse 镜像采集响应阶段
func (i *Inspector) RecordResponse(
	client *cdp.Client,
	ctx context.Context,
	targetID domain.TargetID,
	ev *fetch.RequestPausedReply,
) {
	if !i.IsEnabled() {
		return
	}

	val, ok := i.sessions.Load(ev.RequestID)
	if !ok {
		return
	}
	event := val.(*domain.NetworkEvent)
	defer i.sessions.Delete(ev.RequestID)

	// 采集响应快照
	exec := executor.New(i.log, ev, executor.Options{})

	// 注意：Inspector 仅做展示，默认不拉取大响应体以节省性能
	// 这里可以后续根据需求决定是否拉取 Body
	resInfo := exec.CaptureResponseSnapshot("")

	event.Response = resInfo
	if event.FinalResult == "" {
		event.FinalResult = "passed"
	}

	// 推送完整事件
	i.pushEvent(*event)
}

func (i *Inspector) pushEvent(evt domain.NetworkEvent) {
	if i.events == nil {
		return
	}
	select {
	case i.events <- evt:
	default:
		// 通道满时丢弃，防止阻塞主流程
	}
}

// cleanupLoop 定期清理内存池中的孤儿请求
func (i *Inspector) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now().UnixMilli()
		i.sessions.Range(func(key, value any) bool {
			evt := value.(*domain.NetworkEvent)
			if now-evt.Timestamp > 120000 { // 2分钟超时
				i.sessions.Delete(key)
			}
			return true
		})
	}
}
