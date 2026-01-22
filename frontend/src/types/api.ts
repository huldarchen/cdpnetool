/**
 * 统一后端响应格式
 */
export interface ApiResponse<T = any> {
  success: boolean;
  code?: string;
  message?: string;
  data?: T;
}

/**
 * 空数据类型
 */
export interface EmptyData {}

/**
 * 会话数据
 */
export interface SessionData {
  sessionId: string;
}

/**
 * 浏览器数据
 */
export interface BrowserData {
  devToolsUrl: string;
}

/**
 * 目标信息 (对应 Go 的 domain.TargetInfo)
 */
export interface TargetInfo {
  id: string;
  type: string;
  url: string;
  title: string;
  isCurrent: boolean;
  isUser: boolean;
}

/**
 * 目标列表数据
 */
export interface TargetListData {
  targets: TargetInfo[];
}
