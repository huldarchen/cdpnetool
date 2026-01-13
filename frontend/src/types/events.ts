// 拦截事件相关类型

export type EventType = 
  | 'intercepted'  // 请求被拦截
  | 'continued'    // 请求已放行
  | 'mutated'      // 请求/响应被修改
  | 'failed'       // 请求被阻断
  | 'paused'       // 请求进入 Pending
  | 'approved'     // Pending 被审批
  | 'rejected'     // Pending 被拒绝
  | 'error'        // 发生错误

export interface InterceptEvent {
  id: string
  type: EventType
  session: string
  target: string
  rule?: string
  url?: string
  method?: string
  stage?: 'request' | 'response'
  timestamp: number
  error?: string
  
  // 详情字段（可选）
  statusCode?: number
  requestHeaders?: Record<string, string>
  responseHeaders?: Record<string, string>
  body?: string
}

export const EVENT_TYPE_LABELS: Record<EventType, string> = {
  intercepted: '拦截',
  continued: '放行',
  mutated: '修改',
  failed: '阻断',
  paused: '暂停',
  approved: '审批通过',
  rejected: '审批拒绝',
  error: '错误',
}

export const EVENT_TYPE_COLORS: Record<EventType, { bg: string; text: string }> = {
  intercepted: { bg: 'bg-blue-500/20', text: 'text-blue-500' },
  continued: { bg: 'bg-green-500/20', text: 'text-green-500' },
  mutated: { bg: 'bg-yellow-500/20', text: 'text-yellow-500' },
  failed: { bg: 'bg-red-500/20', text: 'text-red-500' },
  paused: { bg: 'bg-purple-500/20', text: 'text-purple-500' },
  approved: { bg: 'bg-emerald-500/20', text: 'text-emerald-500' },
  rejected: { bg: 'bg-orange-500/20', text: 'text-orange-500' },
  error: { bg: 'bg-rose-500/20', text: 'text-rose-500' },
}

export function createMockEvent(partial: Partial<InterceptEvent> = {}): InterceptEvent {
  return {
    id: Math.random().toString(36).slice(2),
    type: 'intercepted',
    session: '',
    target: '',
    timestamp: Date.now(),
    ...partial,
  }
}
