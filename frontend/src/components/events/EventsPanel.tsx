import { useState, useMemo } from 'react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { 
  Search, 
  X,
  ChevronDown,
  ChevronUp,
  Trash2,
  Filter,
  ArrowRight,
  ArrowLeft,
  Copy,
  Check
} from 'lucide-react'
import type { InterceptEvent, EventType } from '@/types/events'
import { EVENT_TYPE_LABELS, EVENT_TYPE_COLORS } from '@/types/events'

interface EventsPanelProps {
  events: InterceptEvent[]
  onClear?: () => void
}

export function EventsPanel({ events, onClear }: EventsPanelProps) {
  const [search, setSearch] = useState('')
  const [typeFilter, setTypeFilter] = useState<EventType | 'all'>('all')
  const [expandedEvent, setExpandedEvent] = useState<string | null>(null)

  // 过滤事件
  const filteredEvents = useMemo(() => {
    return events.filter(evt => {
      // 类型过滤
      if (typeFilter !== 'all' && evt.type !== typeFilter) return false
      
      // 搜索过滤
      if (search) {
        const searchLower = search.toLowerCase()
        return (
          evt.url?.toLowerCase().includes(searchLower) ||
          evt.method?.toLowerCase().includes(searchLower) ||
          evt.rule?.toLowerCase().includes(searchLower) ||
          evt.target?.toLowerCase().includes(searchLower)
        )
      }
      return true
    })
  }, [events, search, typeFilter])

  // 统计各类型数量
  const typeCounts = useMemo(() => {
    const counts: Record<string, number> = { all: events.length }
    events.forEach(evt => {
      counts[evt.type] = (counts[evt.type] || 0) + 1
    })
    return counts
  }, [events])

  if (events.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-full text-muted-foreground">
        <div className="text-4xl mb-4 opacity-50">📡</div>
        <p>暂无拦截事件</p>
        <p className="text-sm mt-1">启用拦截后，事件将在此实时显示</p>
      </div>
    )
  }

  return (
    <div className="h-full flex flex-col">
      {/* 工具栏 */}
      <div className="flex items-center gap-2 mb-4">
        {/* 搜索框 */}
        <div className="relative flex-1 max-w-xs">
          <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="搜索 URL、方法、规则..."
            className="pl-9 pr-8"
          />
          {search && (
            <button 
              onClick={() => setSearch('')}
              className="absolute right-2.5 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
            >
              <X className="w-4 h-4" />
            </button>
          )}
        </div>

        {/* 类型过滤 */}
        <div className="flex items-center gap-1">
          <Filter className="w-4 h-4 text-muted-foreground" />
          <select
            value={typeFilter}
            onChange={(e) => setTypeFilter(e.target.value as EventType | 'all')}
            className="h-9 px-2 rounded-md border bg-background text-sm"
          >
            <option value="all">全部 ({typeCounts.all})</option>
            {Object.entries(EVENT_TYPE_LABELS).map(([type, label]) => (
              typeCounts[type] > 0 && (
                <option key={type} value={type}>
                  {label} ({typeCounts[type]})
                </option>
              )
            ))}
          </select>
        </div>

        {/* 清除按钮 */}
        {onClear && (
          <Button variant="outline" size="sm" onClick={onClear}>
            <Trash2 className="w-4 h-4 mr-1" />
            清除
          </Button>
        )}
      </div>

      {/* 统计栏 */}
      <div className="flex items-center gap-4 mb-3 text-sm text-muted-foreground">
        <span>共 {filteredEvents.length} 条</span>
        {search && <span>（搜索结果）</span>}
      </div>

      {/* 事件列表 */}
      <ScrollArea className="flex-1">
        <div className="space-y-2 pr-4">
          {filteredEvents.map((evt) => (
            <EventItem
              key={evt.id}
              event={evt}
              isExpanded={expandedEvent === evt.id}
              onToggleExpand={() => setExpandedEvent(expandedEvent === evt.id ? null : evt.id)}
            />
          ))}
        </div>
      </ScrollArea>
    </div>
  )
}

interface EventItemProps {
  event: InterceptEvent
  isExpanded: boolean
  onToggleExpand: () => void
}

function EventItem({ event, isExpanded, onToggleExpand }: EventItemProps) {
  const [copied, setCopied] = useState(false)
  const colors = EVENT_TYPE_COLORS[event.type] || EVENT_TYPE_COLORS.intercepted

  const handleCopyUrl = async () => {
    if (event.url) {
      await navigator.clipboard.writeText(event.url)
      setCopied(true)
      setTimeout(() => setCopied(false), 1500)
    }
  }

  const formatTime = (ts: number) => {
    return new Date(ts).toLocaleTimeString('zh-CN', { 
      hour: '2-digit', 
      minute: '2-digit', 
      second: '2-digit',
      hour12: false 
    })
  }

  return (
    <div className="border rounded-lg bg-card overflow-hidden">
      {/* 头部 */}
      <div 
        className="flex items-center gap-2 p-2.5 cursor-pointer hover:bg-muted/50 transition-colors"
        onClick={onToggleExpand}
      >
        {/* 类型标签 */}
        <Badge variant="outline" className={`${colors.bg} ${colors.text} border-0 text-xs`}>
          {EVENT_TYPE_LABELS[event.type]}
        </Badge>

        {/* 阶段标签 */}
        {event.stage && (
          <Badge variant="outline" className="text-xs">
            {event.stage === 'request' ? (
              <><ArrowRight className="w-3 h-3 mr-0.5" />REQ</>
            ) : (
              <><ArrowLeft className="w-3 h-3 mr-0.5" />RES</>
            )}
          </Badge>
        )}

        {/* Method */}
        {event.method && (
          <span className="font-mono text-xs font-medium px-1.5 py-0.5 rounded bg-muted">
            {event.method}
          </span>
        )}

        {/* URL */}
        <span className="flex-1 text-sm truncate text-muted-foreground font-mono">
          {event.url || '-'}
        </span>

        {/* 时间 */}
        <span className="text-xs text-muted-foreground shrink-0">
          {formatTime(event.timestamp)}
        </span>

        {/* 展开图标 */}
        {isExpanded ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
      </div>

      {/* 展开详情 */}
      {isExpanded && (
        <div className="border-t p-3 space-y-3 text-sm">
          {/* 基本信息 */}
          <div className="grid grid-cols-2 gap-2">
            <div>
              <span className="text-muted-foreground">Target:</span>
              <span className="ml-2 font-mono text-xs">{event.target.slice(0, 20)}...</span>
            </div>
            {event.rule && (
              <div>
                <span className="text-muted-foreground">Rule:</span>
                <span className="ml-2 font-mono text-xs">{event.rule}</span>
              </div>
            )}
            {event.statusCode && (
              <div>
                <span className="text-muted-foreground">Status:</span>
                <span className={`ml-2 font-mono ${
                  event.statusCode >= 400 ? 'text-red-500' : 
                  event.statusCode >= 300 ? 'text-yellow-500' : 'text-green-500'
                }`}>
                  {event.statusCode}
                </span>
              </div>
            )}
          </div>

          {/* URL */}
          {event.url && (
            <div className="space-y-1">
              <div className="flex items-center justify-between">
                <span className="text-muted-foreground">URL</span>
                <Button variant="ghost" size="sm" onClick={handleCopyUrl} className="h-6 px-2">
                  {copied ? <Check className="w-3 h-3" /> : <Copy className="w-3 h-3" />}
                </Button>
              </div>
              <div className="p-2 bg-muted rounded font-mono text-xs break-all">
                {event.url}
              </div>
            </div>
          )}

          {/* Request Headers */}
          {event.requestHeaders && Object.keys(event.requestHeaders).length > 0 && (
            <div className="space-y-1">
              <span className="text-muted-foreground">Request Headers</span>
              <div className="p-2 bg-muted rounded font-mono text-xs max-h-24 overflow-auto">
                {Object.entries(event.requestHeaders).map(([key, value]) => (
                  <div key={key}>
                    <span className="text-primary">{key}:</span> {value}
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Response Headers */}
          {event.responseHeaders && Object.keys(event.responseHeaders).length > 0 && (
            <div className="space-y-1">
              <span className="text-muted-foreground">Response Headers</span>
              <div className="p-2 bg-muted rounded font-mono text-xs max-h-24 overflow-auto">
                {Object.entries(event.responseHeaders).map(([key, value]) => (
                  <div key={key}>
                    <span className="text-primary">{key}:</span> {value}
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Body */}
          {event.body && (
            <div className="space-y-1">
              <span className="text-muted-foreground">Body</span>
              <div className="p-2 bg-muted rounded font-mono text-xs max-h-32 overflow-auto whitespace-pre-wrap">
                {event.body}
              </div>
            </div>
          )}

          {/* Error */}
          {event.error && (
            <div className="p-2 bg-red-500/10 text-red-500 rounded text-xs">
              Error: {event.error}
            </div>
          )}
        </div>
      )}
    </div>
  )
}
