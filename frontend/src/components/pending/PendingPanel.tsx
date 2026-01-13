import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { 
  Check, 
  X, 
  Clock, 
  ArrowRight, 
  ArrowLeft,
  ChevronDown,
  ChevronUp,
  AlertCircle
} from 'lucide-react'
import type { PendingItem, RewriteMutation } from '@/types/pending'

interface PendingPanelProps {
  items: PendingItem[]
  onApprove: (itemId: string, stage: 'request' | 'response', mutations?: RewriteMutation) => void
  onReject: (itemId: string) => void
}

export function PendingPanel({ items, onApprove, onReject }: PendingPanelProps) {
  const [expandedItem, setExpandedItem] = useState<string | null>(null)

  if (items.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center h-full text-muted-foreground">
        <Clock className="w-12 h-12 mb-4 opacity-50" />
        <p>暂无待审批项</p>
        <p className="text-sm mt-1">当规则触发 Pause 动作时，请求会在此等待审批</p>
      </div>
    )
  }

  return (
    <ScrollArea className="h-full">
      <div className="space-y-3 pr-4">
        {items.map((item) => (
          <PendingItemCard
            key={item.id}
            item={item}
            isExpanded={expandedItem === item.id}
            onToggleExpand={() => setExpandedItem(expandedItem === item.id ? null : item.id)}
            onApprove={(mutations) => onApprove(item.id, item.stage, mutations)}
            onReject={() => onReject(item.id)}
          />
        ))}
      </div>
    </ScrollArea>
  )
}

interface PendingItemCardProps {
  item: PendingItem
  isExpanded: boolean
  onToggleExpand: () => void
  onApprove: (mutations?: RewriteMutation) => void
  onReject: () => void
}

function PendingItemCard({ 
  item, 
  isExpanded, 
  onToggleExpand, 
  onApprove, 
  onReject 
}: PendingItemCardProps) {
  const [showMutations, setShowMutations] = useState(false)
  const [mutations, setMutations] = useState<RewriteMutation>({})

  const handleQuickApprove = () => {
    onApprove()
  }

  const handleApproveWithMutations = () => {
    onApprove(mutations)
    setShowMutations(false)
    setMutations({})
  }

  return (
    <div className="border rounded-lg bg-card overflow-hidden">
      {/* 头部 */}
      <div 
        className="flex items-center gap-3 p-3 cursor-pointer hover:bg-muted/50 transition-colors"
        onClick={onToggleExpand}
      >
        <Badge variant={item.stage === 'request' ? 'default' : 'secondary'}>
          {item.stage === 'request' ? (
            <><ArrowRight className="w-3 h-3 mr-1" />请求</>
          ) : (
            <><ArrowLeft className="w-3 h-3 mr-1" />响应</>
          )}
        </Badge>

        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <span className="font-mono text-sm font-medium">{item.method}</span>
            <span className="text-sm truncate text-muted-foreground">{item.url}</span>
          </div>
          {item.rule && (
            <div className="text-xs text-muted-foreground mt-0.5">
              规则: {item.rule}
            </div>
          )}
        </div>

        <div className="flex items-center gap-1">
          <Button
            variant="default"
            size="sm"
            onClick={(e) => { e.stopPropagation(); handleQuickApprove() }}
          >
            <Check className="w-4 h-4 mr-1" />
            放行
          </Button>
          <Button
            variant="destructive"
            size="sm"
            onClick={(e) => { e.stopPropagation(); onReject() }}
          >
            <X className="w-4 h-4 mr-1" />
            拒绝
          </Button>
          {isExpanded ? <ChevronUp className="w-4 h-4 ml-2" /> : <ChevronDown className="w-4 h-4 ml-2" />}
        </div>
      </div>

      {/* 展开详情 */}
      {isExpanded && (
        <div className="border-t p-4 space-y-4">
          {/* 基本信息 */}
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span className="text-muted-foreground">Target ID:</span>
              <span className="ml-2 font-mono">{item.target.slice(0, 16)}...</span>
            </div>
            <div>
              <span className="text-muted-foreground">Pending ID:</span>
              <span className="ml-2 font-mono">{item.id.slice(0, 16)}...</span>
            </div>
          </div>

          {/* URL */}
          <div className="space-y-1">
            <label className="text-sm text-muted-foreground">完整 URL</label>
            <div className="p-2 bg-muted rounded font-mono text-sm break-all">
              {item.url}
            </div>
          </div>

          {/* Headers（如果有） */}
          {item.headers && Object.keys(item.headers).length > 0 && (
            <div className="space-y-1">
              <label className="text-sm text-muted-foreground">Headers</label>
              <div className="p-2 bg-muted rounded font-mono text-xs max-h-32 overflow-auto">
                {Object.entries(item.headers).map(([key, value]) => (
                  <div key={key}>
                    <span className="text-primary">{key}:</span> {value}
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Body（如果有） */}
          {item.body && (
            <div className="space-y-1">
              <label className="text-sm text-muted-foreground">Body</label>
              <div className="p-2 bg-muted rounded font-mono text-xs max-h-40 overflow-auto whitespace-pre-wrap">
                {item.body}
              </div>
            </div>
          )}

          {/* 修改选项 */}
          <div className="border-t pt-4">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowMutations(!showMutations)}
            >
              {showMutations ? '隐藏修改选项' : '修改后放行'}
            </Button>

            {showMutations && (
              <MutationEditor
                mutations={mutations}
                onChange={setMutations}
                stage={item.stage}
                onApply={handleApproveWithMutations}
              />
            )}
          </div>
        </div>
      )}
    </div>
  )
}

// 修改编辑器
interface MutationEditorProps {
  mutations: RewriteMutation
  onChange: (mutations: RewriteMutation) => void
  stage: 'request' | 'response'
  onApply: () => void
}

function MutationEditor({ mutations, onChange, stage, onApply }: MutationEditorProps) {
  const [headersText, setHeadersText] = useState('')

  const parseHeaders = (text: string): Record<string, string | null> => {
    const result: Record<string, string | null> = {}
    text.split('\n').forEach(line => {
      const [key, ...valueParts] = line.split(':')
      if (key?.trim()) {
        const value = valueParts.join(':').trim()
        result[key.trim()] = value || null
      }
    })
    return result
  }

  return (
    <div className="mt-4 space-y-4 p-4 border rounded-lg bg-muted/30">
      <div className="flex items-center gap-2 text-sm text-amber-500">
        <AlertCircle className="w-4 h-4" />
        <span>修改将应用于{stage === 'request' ? '请求' : '响应'}后继续</span>
      </div>

      {stage === 'request' && (
        <>
          {/* URL 修改 */}
          <div className="space-y-1">
            <label className="text-sm font-medium">重写 URL</label>
            <Input
              value={mutations.url || ''}
              onChange={(e) => onChange({ ...mutations, url: e.target.value || undefined })}
              placeholder="留空不修改"
            />
          </div>

          {/* Method 修改 */}
          <div className="space-y-1">
            <label className="text-sm font-medium">重写 Method</label>
            <Input
              value={mutations.method || ''}
              onChange={(e) => onChange({ ...mutations, method: e.target.value || undefined })}
              placeholder="留空不修改 (GET/POST/PUT/DELETE...)"
            />
          </div>
        </>
      )}

      {/* Headers 修改 */}
      <div className="space-y-1">
        <label className="text-sm font-medium">修改 Headers (每行一个: Key: Value)</label>
        <Textarea
          value={headersText}
          onChange={(e) => {
            setHeadersText(e.target.value)
            const headers = parseHeaders(e.target.value)
            if (Object.keys(headers).length > 0) {
              onChange({ ...mutations, headers })
            } else {
              const { headers: _, ...rest } = mutations
              onChange(rest)
            }
          }}
          placeholder="Content-Type: application/json&#10;X-Custom-Header: value"
          rows={3}
          className="font-mono text-sm"
        />
      </div>

      {/* Body 修改（简化版：直接替换） */}
      <div className="space-y-1">
        <label className="text-sm font-medium">Body Base64 替换</label>
        <Textarea
          value={mutations.body?.base64?.value || ''}
          onChange={(e) => {
            if (e.target.value) {
              onChange({
                ...mutations,
                body: { base64: { value: e.target.value } }
              })
            } else {
              const { body: _, ...rest } = mutations
              onChange(rest)
            }
          }}
          placeholder="Base64 编码的内容（留空不修改）"
          rows={2}
          className="font-mono text-sm"
        />
      </div>

      <Button onClick={onApply} className="w-full">
        <Check className="w-4 h-4 mr-2" />
        应用修改并放行
      </Button>
    </div>
  )
}
