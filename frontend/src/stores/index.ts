import { create } from 'zustand'
import type { 
  NetworkEvent,
  MatchedEventWithId, 
  UnmatchedEventWithId 
} from '@/types/events'

// 类型定义
export interface TargetInfo {
  id: string
  type: string
  url: string
  title: string
  isCurrent: boolean
  isUser: boolean
}

// Session 状态
interface SessionState {
  currentSessionId: string | null
  devToolsURL: string
  isConnected: boolean
  isIntercepting: boolean
  collectUnmatched: boolean // 是否收集未匹配的请求
  activeConfigId: number | null
  targets: TargetInfo[]
  attachedTargetId: string | null
  matchedEvents: MatchedEventWithId[]    // 匹配的事件（会存入数据库）
  unmatchedEvents: UnmatchedEventWithId[] // 未匹配的事件（仅内存）
  
  // Actions
  setDevToolsURL: (url: string) => void
  setCurrentSession: (id: string | null) => void
  setConnected: (connected: boolean) => void
  setIntercepting: (intercepting: boolean) => void
  setCollectUnmatched: (enabled: boolean) => void
  setActiveConfigId: (id: number | null) => void
  setTargets: (targets: TargetInfo[]) => void
  setAttachedTargetId: (targetId: string | null) => void
  resetSession: () => void
  
  // 事件操作
  addInterceptEvent: (event: NetworkEvent) => void
  clearMatchedEvents: () => void
  clearUnmatchedEvents: () => void
  clearAllEvents: () => void
}

// 生成事件 ID
function generateEventId(timestamp: number): string {
  return `${timestamp}_${Math.random().toString(36).slice(2, 10)}`
}

export const useSessionStore = create<SessionState>((set) => ({
  currentSessionId: null,
  devToolsURL: 'http://localhost:9222',
  isConnected: false,
  isIntercepting: false,
  collectUnmatched: false,
  activeConfigId: null,
  targets: [],
  attachedTargetId: null,
  matchedEvents: [],
  unmatchedEvents: [],
  
  setDevToolsURL: (url) => set({ devToolsURL: url }),
  setCurrentSession: (id) => set({ currentSessionId: id }),
  setConnected: (connected) => set({ isConnected: connected }),
  setIntercepting: (intercepting) => set({ isIntercepting: intercepting }),
  setCollectUnmatched: (enabled) => set({ collectUnmatched: enabled }),
  setActiveConfigId: (id) => set({ activeConfigId: id }),
  setTargets: (targets) => set({ targets }),
  setAttachedTargetId: (targetId) => set({ attachedTargetId: targetId }),

  resetSession: () => set({
    attachedTargetId: null,
    activeConfigId: null,
    isIntercepting: false,
    targets: [],
  }),
  
  // 添加事件（根据 isMatched 分开存储）
  addInterceptEvent: (event) => set((state) => {
    console.log('[Store] 处理拦截事件:', event)
    
    // 后端现在直接发送 NetworkEvent 扁平对象
    // 我们根据 isMatched 字段进行分发，并为 UI 组件构造预期的嵌套结构
    if (event.isMatched) {
      const eventWithId: MatchedEventWithId = {
        networkEvent: event,
        id: generateEventId(event.timestamp),
      }
      return {
        matchedEvents: [eventWithId, ...state.matchedEvents].slice(0, 200)
      }
    } else {
      const eventWithId: UnmatchedEventWithId = {
        networkEvent: event,
        id: generateEventId(event.timestamp),
      }
      return {
        unmatchedEvents: [eventWithId, ...state.unmatchedEvents].slice(0, 100)
      }
    }
  }),
  
  clearMatchedEvents: () => set({ matchedEvents: [] }),
  clearUnmatchedEvents: () => set({ unmatchedEvents: [] }),
  clearAllEvents: () => set({ matchedEvents: [], unmatchedEvents: [] }),
}))

// 主题状态
interface ThemeState {
  isDark: boolean
  toggle: () => void
}

export const useThemeStore = create<ThemeState>((set) => ({
  isDark: true,
  toggle: () => set((state) => {
    const newIsDark = !state.isDark
    if (newIsDark) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
    return { isDark: newIsDark }
  }),
}))

// 初始化主题
if (typeof window !== 'undefined') {
  document.documentElement.classList.add('dark')
}
