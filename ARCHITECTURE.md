# 架构设计指南（阶段一）

- 项目名称：cdpnetool（Chrome DevTools Protocol Network Tools）
- 目标：基于CDP与浏览器通信，拦截并可控篡改网页的网络请求与响应
- GUI策略：GUI与核心逻辑分离，核心以Go库形式提供稳定API与事件流；GUI后续通过进程内或IPC集成
- 依赖：CDP客户端库采用 github.com/mafredri/cdp（Chromium/Chrome/Edge优先）

## 阶段划分与范围

- 阶段一（核心功能）：CDP连接与会话管理、请求/响应拦截、规则引擎（复杂匹配器与内容修改器）、并发与背压、基础可观测性（结构化日志与基础指标）、基础配置（JSON）
- 暂缓实现：权限/鉴权、热更新与回滚、插件与脚本扩展、IPC服务、复杂GUI交互与远程控制、审计与合规的高级特性

## 架构总览

- 分层模块
  - cdp/session：DevTools发现与连接、会话生命周期、Target管理、Fetch/Network启停
  - intercept/handler：Fetch.requestPaused处理、执行continue/fulfill/fail等动作
  - rules/engine：复杂匹配器与内容修改器、优先级与冲突处理、命中统计与缓存
  - payload/store：大Body外存分块与临时文件管理、对象池与零拷贝优化
  - observability：结构化日志、基础指标、事件追踪ID与关联
  - config/runtime：规则与运行配置的JSON加载与校验、默认值管理
- 事件流
  - CDP事件→入队→规则匹配→执行动作→调用CDP返回→上报结果事件与指标
  - 背压与降级：队列超限或超时自动继续原请求并记录降级事件

## 两阶段拦截模型

- 请求阶段：浏览器→拦截与篡改→服务器
  - Fetch.requestPaused(stage=Request)→匹配→continueRequest/fulfillRequest/failRequest
- 响应阶段：服务器→拦截与篡改→浏览器
  - Fetch.requestPaused(stage=Response)→读取响应头与必要时体→匹配→continueResponse/fulfillRequest/failRequest
  - 必要时结合Network.responseReceived与Network.getResponseBody获取体；超出体量阈值走摘要与截断策略
- RequestPattern默认启用
  - 资源类型：XHR、Fetch；可选扩展到Document/Script/Stylesheet/Media（按需）
  - 阶段：Request与Response均启用
  - 回退：若不支持Response阶段拦截，则以Network事件观察与有限替换策略降级

## 规则引擎设计

- 匹配器（复杂匹配）
  - 组合逻辑：allOf/anyOf/noneOf，支持嵌套；默认allOf
  - 条件类型：
    - URL：glob/regex/prefix/exact
    - 方法：集合匹配
    - 头部：键存在/equals/contains/regex（键大小写不敏感）
    - Query/Cookie：键值equals/contains/regex
    - Body：
      - JSON Pointer：存在/equals/regex
      - 文本：mime为text/*时支持contains/regex
      - 体量与类型：size范围、mime前缀
    - 时间窗口与概率：支持窗口化规则与采样
- 内容修改器（请求/响应）
  - 请求修改：URL/方法替换；头/Query/Cookie的set/remove/append；Body支持JSON Patch（RFC6902）、文本regex替换、二进制base64
  - 响应修改：状态码、头部与Body替换；响应体通过fulfillRequest返回
  - 组合与顺序：修改器按声明顺序执行；失败策略可配置（继续原值/失败）
- 优先级与冲突
  - 显式priority与稳定排序；默认短路（命中后停止后续规则）
  - 可选聚合执行（将多个动作组合应用）
  - 命中上限与采样策略可配置（阶段二扩展）

## 手动阻塞修改机制

- 目标：请求前或响应前暂停，等待用户手动编辑；超时自动放行（默认继续原请求或应用默认动作）
- 规则触发：
  - 动作pause：包含stage=request/response、timeoutMS、defaultAction（continue_original/continue_mutated/fulfill/fail）
- 事件与审批：
  - 生成PendingItem事件（URL/方法/头/Query/Body预览、会话/目标ID），进入待审批队列
  - 接口（为GUI/CLI预留）：SubscribePending()、ApproveRequest(itemID, mutations)、ApproveResponse(itemID, mutations)、Reject(itemID)
  - 超时与降级：timeoutMS到达未审批，执行defaultAction并记录事件（完整审计在阶段二开启）

## 公共API（最小稳定集草案）

```go
package api

type SessionID string
type TargetID string
type RuleID string

type SessionConfig struct {
    DevToolsURL       string
    Concurrency       int           // worker数量
    BodySizeThreshold int64         // 外存分块阈值，默认4MB
    PendingCapacity   int           // 待审批队列容量，默认64
    ProcessTimeoutMS  int           // 拦截处理最大耗时，默认200ms
}

type RuleSet struct {
    Version string          // "1.0"
    Rules   []Rule
}

type Rule struct {
    ID       RuleID
    Priority int
    Mode     string         // "short_circuit" | "aggregate"
    Match    Match          // 复杂匹配器
    Action   Action         // 内容修改器 + 控制动作
}

type Match struct {
    AllOf      []Condition
    AnyOf      []Condition
    NoneOf     []Condition
}

type Condition struct {
    Type   string           // "url"|"method"|"header"|"query"|"cookie"|"json_pointer"|"text"|"mime"|"size"|"time_window"|"probability"
    // 依据Type包含相应字段：pattern/mode/values/key/op/value/pointer等
}

type Action struct {
    Rewrite   *Rewrite
    Respond   *Respond
    Fail      *Fail
    DelayMS   int
    DropRate  float64
    Pause     *Pause
}

type Rewrite struct {
    URL     *string
    Method  *string
    Headers map[string]*string
    Query   map[string]*string
    Cookies map[string]*string
    Body    *BodyPatch
}

type BodyPatch struct {
    Type string               // "json_patch"|"text_regex"|"base64"
    Ops  []any                // 根据Type解析
}

type Respond struct {
    Status  int
    Headers map[string]string
    Body    []byte
    Base64  bool
}

type Fail struct {
    Reason string             // 例如 "ConnectionFailed"
}

type Pause struct {
    Stage        string       // "request"|"response"
    TimeoutMS    int          // 默认3000
    DefaultAction struct {
        Type   string         // "continue_original"|"continue_mutated"|"fulfill"|"fail"
        Status int            // 当Type=fulfill
        Reason string         // 当Type=fail
    }
}

type EngineStats struct {
    Total   int64
    Matched int64
    ByRule  map[RuleID]int64
}

type Event struct {
    Type    string            // "intercepted"|"mutated"|"fulfilled"|"failed"|"degraded"|"pending"|"approved"|"timeout"
    Session SessionID
    Target  TargetID
    Rule    *RuleID
    Error   error
}

type Service interface {
    StartSession(cfg SessionConfig) (SessionID, error)
    StopSession(id SessionID) error
    AttachTarget(id SessionID, target TargetID) error
    DetachTarget(id SessionID, target TargetID) error

    EnableInterception(id SessionID) error
    DisableInterception(id SessionID) error

    LoadRules(id SessionID, rs RuleSet) error
    GetRuleStats(id SessionID) (EngineStats, error)

    SubscribeEvents(id SessionID) (<-chan Event, error)
    SubscribePending(id SessionID) (<-chan any, error) // PendingItem
    ApproveRequest(itemID string, mutations Rewrite) error
    ApproveResponse(itemID string, mutations Rewrite) error
    Reject(itemID string) error
}
```

## 规则JSON示例

```json
{
  "version": "1.0",
  "rules": [
    {
      "id": "manual_edit_payment",
      "priority": 200,
      "mode": "short_circuit",
      "match": {
        "allOf": [
          { "type": "url", "mode": "prefix", "pattern": "https://api.example.com/payment" },
          { "type": "method", "values": ["POST"] }
        ]
      },
      "action": {
        "pause": {
          "stage": "request",
          "timeoutMS": 5000,
          "defaultAction": { "type": "continue_original" }
        }
      }
    },
    {
      "id": "inject_500",
      "priority": 90,
      "mode": "short_circuit",
      "match": {
        "anyOf": [
          { "type": "url", "mode": "regex", "pattern": "/payment" },
          { "type": "header", "key": "X-Route", "op": "equals", "value": "pay" }
        ],
        "noneOf": [
          { "type": "probability", "value": 0.0 }
        ]
      },
      "action": {
        "respond": {
          "status": 500,
          "headers": { "Content-Type": "application/json" },
          "body": "eyJlcnJvciI6ICJmYWlsZWQifQ==",
          "base64": true
        },
        "delayMS": 50
      }
    }
  ]
}
```

## CDP集成与调用序列（mafredri/cdp）

- 启动与启用
  - 使用DevTools端点发现Target→建立WebSocket→创建cdp.Client
  - Network.enable→Fetch.enable({patterns: [{stage=Request, resourceTypes=[XHR,Fetch]}, {stage=Response, resourceTypes=[XHR,Fetch]}]})
- 请求阶段处理
  - on Fetch.requestPaused(stage=Request)→构造上下文→规则匹配→continueRequest/fulfillRequest/failRequest
- 响应阶段处理
  - on Fetch.requestPaused(stage=Response)→读取响应头→必要时Network.getResponseBody→规则匹配→continueResponse/fulfillRequest/failRequest
- 回退
  - 若Response拦截缺失，则以Network.responseReceived观察并有限替换（主要头与错误注入）

## 并发与背压策略

- 队列：优先级队列；fast-path（无暂停的小改写）与pending队列（需人工审批）
- Worker Pool：固定或弹性并发数，按SessionConfig.Concurrency控制
- 超时与降级：每拦截事件设最大处理时限（ProcessTimeoutMS，默认200ms）；超时自动continue原请求并记录“degraded”
- pending容量：PendingCapacity（默认64）；满载时直接执行defaultAction并上报降级

## 性能与可靠性目标

- 延迟预算：单次拦截处理开销控制在毫秒级（典型<5ms，不含手动暂停）
- 吞吐：高并发下维持稳定处理；具备限流与降级路径
- 资源管理：BodySizeThreshold（默认4MB）以上走外存分块与摘要；对象池用于临时结构降低GC压力
- 错误处理：区分可恢复/不可恢复；重试与回退策略针对CDP连接与Fetch启停

## 可观测性与日志

- 结构化日志：控制台与文件输出；默认遮盖Cookie、Authorization等敏感字段
- 指标：QPS、规则命中率、延迟分布、错误统计、队列长度、pending超时比例、外存使用
- 事件追踪：为每拦截事件生成TraceID并贯穿生命周期

## 配置管理（阶段一）

- 规则与运行配置均采用JSON；提供Schema校验与明确错误提示
- 默认参数：
  - Concurrency：CPU核数或配置值
  - BodySizeThreshold：4MB
  - ProcessTimeoutMS：200
  - PendingCapacity：64
  - Pause默认timeoutMS：3000

## 兼容与支持范围

- 浏览器：Chromium系稳定版优先；针对CDP能力进行检测与优雅回退
- 平台：Windows/macOS/Linux；Go 1.21+（建议1.22）

## 非目标与边界（阶段一）

- 不做系统级代理/MITM；仅作用于受控CDP会话
- 不持久化敏感数据；日志默认遮盖
- 权限/鉴权、审计与合规高级特性将在后续阶段规划

## 里程碑

- M1：最小可用（连接、两阶段拦截、基本规则与修改器、日志与指标）
- M2：规则引擎完善（优先级、聚合、复杂Body处理与优化）
- M3：可观测性与稳定性增强（更全面指标、事件追踪、降级策略打磨）
- M4：扩展能力（热更新、IPC、GUI对接、插件/脚本）

## 验收准则（阶段一）

- 稳定拦截：在Chromium稳定版上持续拦截与处理请求/响应
- 规则正确性：表驱动测试覆盖复杂匹配器与修改器核心路径
- 性能与资源：典型场景延迟毫秒级；高并发具备背压与降级；外存与对象池生效
- 可观测性：日志与指标完整；支持定位问题与分析命中
- 配置与错误提示：规则与配置加载可靠，错误信息清晰

## 后续工作

- 提供规则JSON Schema草案，保证配置校验与版本化
- 完善API与伪实现流程图，开始搭建代码框架
