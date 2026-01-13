// Pending 审批相关类型

export interface PendingItem {
  id: string
  stage: 'request' | 'response'
  url: string
  method: string
  target: string
  rule?: string
  // 扩展字段（未来可从后端获取更多详情）
  headers?: Record<string, string>
  body?: string
  timestamp?: number
}

export interface RewriteMutation {
  url?: string
  method?: string
  headers?: Record<string, string | null>
  query?: Record<string, string | null>
  cookies?: Record<string, string | null>
  body?: {
    jsonPatch?: Array<{
      op: 'add' | 'remove' | 'replace' | 'move' | 'copy'
      path: string
      from?: string
      value?: any
    }>
    textRegex?: {
      pattern: string
      replace: string
    }
    base64?: {
      value: string
    }
  }
}
