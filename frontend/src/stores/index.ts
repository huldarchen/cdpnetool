import { create } from 'zustand'
import type { 
  NetworkEvent,
  MatchedEventWithId 
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
  activeConfigId: number | null
  targets: TargetInfo[]
  attachedTargetId: string | null
  matchedEvents: MatchedEventWithId[]    // 匹配的事件（会存入数据库）
  isTrafficCapturing: boolean           // 是否正在捕获全量流量
  trafficEvents: NetworkEvent[]         // 全量流量列表（仅内存，最近100条）
  
  // Actions
  setDevToolsURL: (url: string) => void
  setCurrentSession: (id: string | null) => void
  setConnected: (connected: boolean) => void
  setIntercepting: (intercepting: boolean) => void
  setTrafficCapturing: (capturing: boolean) => void
  setActiveConfigId: (id: number | null) => void
  setTargets: (targets: TargetInfo[]) => void
  setAttachedTargetId: (targetId: string | null) => void
  resetSession: () => void
  
  // 事件操作
  addInterceptEvent: (event: NetworkEvent) => void
  addTrafficEvent: (event: NetworkEvent) => void
  clearMatchedEvents: () => void
  clearTrafficEvents: () => void
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
  activeConfigId: null,
  targets: [],
  attachedTargetId: null,
  matchedEvents: [],
  isTrafficCapturing: false,
  trafficEvents: [],
  
  setDevToolsURL: (url) => set({ devToolsURL: url }),
  setCurrentSession: (id) => set({ currentSessionId: id }),
  setConnected: (connected) => set({ isConnected: connected }),
  setIntercepting: (intercepting) => set({ isIntercepting: intercepting }),
  setTrafficCapturing: (capturing) => set({ isTrafficCapturing: capturing }),
  setActiveConfigId: (id) => set({ activeConfigId: id }),
  setTargets: (targets) => set({ targets }),
  setAttachedTargetId: (targetId) => set({ attachedTargetId: targetId }),

  resetSession: () => set({
    attachedTargetId: null,
    activeConfigId: null,
    isIntercepting: false,
    isTrafficCapturing: false,
    targets: [],
    trafficEvents: [],
  }),
  
  // 添加事件
  addInterceptEvent: (event) => set((state) => {
    console.log('[Store] 处理拦截事件:', event)
    
    // 后端现在只发送匹配成功的 NetworkEvent
    if (event.isMatched) {
      const eventWithId: MatchedEventWithId = {
        networkEvent: event,
        id: generateEventId(event.timestamp),
      }
      return {
        matchedEvents: [eventWithId, ...state.matchedEvents].slice(0, 200)
      }
    }
    return state
  }),

  // 添加全量流量事件 (Upsert 逻辑)
  addTrafficEvent: (event) => set((state) => {
    const existingIndex = state.trafficEvents.findIndex(e => e.id === event.id)
    if (existingIndex > -1) {
      const newList = [...state.trafficEvents]
      // 合并数据：响应阶段会补全 response 和 finalResult
      newList[existingIndex] = { ...newList[existingIndex], ...event }
      return { trafficEvents: newList }
    } else {
      // 新请求，限制 100 条
      return { trafficEvents: [event, ...state.trafficEvents].slice(0, 100) }
    }
  }),
  
  clearMatchedEvents: () => set({ matchedEvents: [] }),
  clearTrafficEvents: () => set({ trafficEvents: [] }),
  clearAllEvents: () => set({ matchedEvents: [], trafficEvents: [] }),
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
