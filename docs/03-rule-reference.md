# 规则配置参考

本章节详细介绍 cdpnetool 的规则配置格式、所有支持的匹配条件和执行行为。

---

## 配置文件结构

### Config 根对象

配置文件的根对象包含配置元信息和规则列表：

```json
{
  "id": "config-20260118-abc123",
  "name": "我的配置",
  "version": "1.0",
  "description": "配置说明",
  "settings": {},
  "rules": []
}
```

**字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `id` | string | 是 | 配置唯一标识符，格式：`config-YYYYMMDD-随机6位` |
| `name` | string | 是 | 配置名称 |
| `version` | string | 是 | 配置版本（当前为 1.0） |
| `description` | string | 否 | 配置描述 |
| `settings` | object | 否 | 预留设置项 |
| `rules` | array | 是 | 规则列表数组 |

---

### Rule 规则对象

每条规则包含匹配条件和执行行为：

```json
{
  "id": "rule-001",
  "name": "规则名称",
  "enabled": true,
  "priority": 0,
  "stage": "request",
  "match": {
    "allOf": [],
    "anyOf": []
  },
  "actions": []
}
```

**字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `id` | string | 是 | 规则唯一标识符，格式：`rule-XXX` |
| `name` | string | 是 | 规则名称 |
| `enabled` | boolean | 是 | 是否启用 |
| `priority` | number | 是 | 优先级，数值越大越先执行 |
| `stage` | string | 是 | 生命周期阶段（`request` 或 `response`） |
| `match` | object | 是 | 匹配条件对象 |
| `actions` | array | 是 | 执行行为数组 |

---

## 生命周期阶段（Stage）

### request - 请求阶段

- 在请求发送到服务器之前拦截
- 可以修改请求 URL、方法、请求头、请求体等
- 可以阻止请求发送（block 行为）

**适用场景：**
- 修改请求参数
- 添加或修改请求头
- 阻止特定请求

---

### response - 响应阶段

- 在响应返回到浏览器之前拦截
- 可以修改响应状态码、响应头、响应体等
- 不能阻止请求（请求已发送）

**适用场景：**
- Mock API 响应数据
- 修改响应内容
- 解决跨域问题

---

## 匹配条件（Match）完整参考

### 条件逻辑组合

匹配条件使用 `allOf` 和 `anyOf` 进行逻辑组合：

**AllOf（AND 逻辑）：**
- 所有条件都必须满足才匹配
- 条件数组，按顺序评估

**AnyOf（OR 逻辑）：**
- 任意一个条件满足即匹配
- 条件数组，按顺序评估

**示例：**
```json
{
  "match": {
    "allOf": [
      {"type": "urlContains", "value": "/api/"},
      {"type": "method", "values": ["POST"]}
    ],
    "anyOf": [
      {"type": "headerExists", "name": "Authorization"},
      {"type": "cookieExists", "name": "token"}
    ]
  }
}
```

> 💡 **规则**：当 `allOf` 和 `anyOf` 同时存在时，`allOf` 中的所有条件必须满足，且 `anyOf` 中至少有一个条件满足。

---

### URL 条件类型

#### urlEquals

**说明：** URL 精确匹配

**参数：**
- `value` (string) - 完整 URL 字符串

**示例：**
```json
{"type": "urlEquals", "value": "https://example.com/api/user"}
```

---

#### urlPrefix

**说明：** URL 前缀匹配

**参数：**
- `value` (string) - URL 前缀字符串

**示例：**
```json
{"type": "urlPrefix", "value": "https://example.com/api/"}
```

---

#### urlSuffix

**说明：** URL 后缀匹配

**参数：**
- `value` (string) - URL 后缀字符串

**示例：**
```json
{"type": "urlSuffix", "value": ".json"}
```

---

#### urlContains

**说明：** URL 包含指定字符串

**参数：**
- `value` (string) - 要包含的字符串

**示例：**
```json
{"type": "urlContains", "value": "/api/user"}
```

---

#### urlRegex

**说明：** URL 正则表达式匹配

**参数：**
- `pattern` (string) - 正则表达式字符串

**示例：**
```json
{"type": "urlRegex", "pattern": "^https://example\\.com/api/(user|order)/\\d+$"}
```

---

### HTTP 属性条件

#### method

**说明：** HTTP 方法匹配

**参数：**
- `values` (string[]) - HTTP 方法数组

**可选值：** `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `HEAD`, `OPTIONS`

**示例：**
```json
{"type": "method", "values": ["GET", "POST"]}
```

---

#### resourceType

**说明：** 资源类型匹配

**参数：**
- `values` (string[]) - 资源类型数组

**可选值：**

| 值 | 说明 |
|---|---|
| `document` | HTML 文档 |
| `script` | JavaScript |
| `stylesheet` | CSS |
| `image` | 图片 |
| `media` | 音视频 |
| `font` | 字体 |
| `xhr` | XMLHttpRequest |
| `fetch` | Fetch API |
| `websocket` | WebSocket |
| `other` | 其他 |

**示例：**
```json
{"type": "resourceType", "values": ["xhr", "fetch"]}
```

---

### Header 条件类型

#### headerExists

**说明：** 判断 Header 是否存在

**参数：**
- `name` (string) - Header 名称（不区分大小写）

**示例：**
```json
{"type": "headerExists", "name": "Authorization"}
```

---

#### headerNotExists

**说明：** 判断 Header 是否不存在

**参数：**
- `name` (string) - Header 名称

**示例：**
```json
{"type": "headerNotExists", "name": "X-Custom-Header"}
```

---

#### headerEquals

**说明：** Header 值精确匹配

**参数：**
- `name` (string) - Header 名称
- `value` (string) - 匹配值

**示例：**
```json
{"type": "headerEquals", "name": "Content-Type", "value": "application/json"}
```

---

#### headerContains

**说明：** Header 值包含指定字符串

**参数：**
- `name` (string) - Header 名称
- `value` (string) - 包含的字符串

**示例：**
```json
{"type": "headerContains", "name": "User-Agent", "value": "Chrome"}
```

---

#### headerRegex

**说明：** Header 值正则表达式匹配

**参数：**
- `name` (string) - Header 名称
- `pattern` (string) - 正则表达式

**示例：**
```json
{"type": "headerRegex", "name": "Authorization", "pattern": "^Bearer\\s+[A-Za-z0-9\\-_]+$"}
```

---

### Query 参数条件类型

Query 参数条件与 Header 条件类似，用于匹配 URL 查询参数：

- `queryExists` - Query 参数存在
- `queryNotExists` - Query 参数不存在
- `queryEquals` - Query 参数值精确匹配
- `queryContains` - Query 参数值包含匹配
- `queryRegex` - Query 参数值正则匹配

**参数：** `name` - 参数名称，其他参数与 Header 条件相同

**示例：**
```json
{"type": "queryEquals", "name": "page", "value": "1"}
```

---

### Cookie 条件类型

Cookie 条件与 Header 条件类似，用于匹配 HTTP Cookie：

- `cookieExists` - Cookie 存在
- `cookieNotExists` - Cookie 不存在
- `cookieEquals` - Cookie 值精确匹配
- `cookieContains` - Cookie 值包含匹配
- `cookieRegex` - Cookie 值正则匹配

**参数：** `name` - Cookie 名称，其他参数与 Header 条件相同

**示例：**
```json
{"type": "cookieExists", "name": "sessionId"}
```

---

### Body 条件类型

#### bodyContains

**说明：** Body 包含指定字符串

**参数：**
- `value` (string) - 要包含的字符串

**示例：**
```json
{"type": "bodyContains", "value": "username"}
```

---

#### bodyRegex

**说明：** Body 正则表达式匹配

**参数：**
- `pattern` (string) - 正则表达式

**示例：**
```json
{"type": "bodyRegex", "pattern": "\"userId\":\\s*\\d+"}
```

---

#### bodyJsonPath

**说明：** 使用 JSON Path 匹配 Body 内容

**参数：**
- `path` (string) - JSON Path 表达式
- `value` (string) - 期望值

**示例：**
```json
{"type": "bodyJsonPath", "path": "$.user.id", "value": "123"}
```

---

## 执行行为（Actions）完整参考

### 请求阶段专用行为

以下行为仅在 `stage: "request"` 时可用：

#### setUrl

**说明：** 设置请求 URL

**参数：**
- `value` (string) - 新的 URL 字符串

**示例：**
```json
{"type": "setUrl", "value": "https://example.com/api/v2/user"}
```

---

#### setMethod

**说明：** 设置请求方法

**参数：**
- `value` (string) - HTTP 方法（GET/POST/PUT/DELETE/等）

**示例：**
```json
{"type": "setMethod", "value": "POST"}
```

---

#### setQueryParam

**说明：** 设置 URL 查询参数

**参数：**
- `name` (string) - 参数名称
- `value` (any) - 参数值

**示例：**
```json
{"type": "setQueryParam", "name": "page", "value": "1"}
```

---

#### removeQueryParam

**说明：** 移除 URL 查询参数

**参数：**
- `name` (string) - 参数名称

**示例：**
```json
{"type": "removeQueryParam", "name": "debug"}
```

---

#### setCookie

**说明：** 设置 Cookie

**参数：**
- `name` (string) - Cookie 名称
- `value` (any) - Cookie 值

**示例：**
```json
{"type": "setCookie", "name": "token", "value": "abc123"}
```

---

#### removeCookie

**说明：** 移除 Cookie

**参数：**
- `name` (string) - Cookie 名称

**示例：**
```json
{"type": "removeCookie", "name": "tracking_id"}
```

---

#### setFormField

**说明：** 设置表单字段（适用于 `application/x-www-form-urlencoded`）

**参数：**
- `name` (string) - 字段名称
- `value` (any) - 字段值

**示例：**
```json
{"type": "setFormField", "name": "username", "value": "testuser"}
```

---

#### removeFormField

**说明：** 移除表单字段

**参数：**
- `name` (string) - 字段名称

**示例：**
```json
{"type": "removeFormField", "name": "csrf_token"}
```

---

#### block

**说明：** 拦截请求并返回自定义响应（终结性行为，后续行为不再执行）

**参数：**
- `statusCode` (number, 可选) - HTTP 状态码，默认 200
- `headers` (object, 可选) - 响应头对象
- `body` (string, 可选) - 响应体字符串
- `bodyEncoding` (string, 可选) - Body 编码方式（`text` 或 `base64`），默认 `text`

**示例：**
```json
{
  "type": "block",
  "statusCode": 403,
  "headers": {"Content-Type": "application/json"},
  "body": "{\"error\": \"Blocked by rule\"}",
  "bodyEncoding": "text"
}
```

---

### 响应阶段专用行为

以下行为仅在 `stage: "response"` 时可用：

#### setStatus

**说明：** 设置响应状态码

**参数：**
- `value` (number) - HTTP 状态码

**示例：**
```json
{"type": "setStatus", "value": 200}
```

---

### 通用行为（请求/响应均可用）

以下行为在两个阶段均可使用：

#### setHeader

**说明：** 设置头部（请求头或响应头）

**参数：**
- `name` (string) - Header 名称
- `value` (any) - Header 值

**示例：**
```json
{"type": "setHeader", "name": "X-Custom-Header", "value": "custom-value"}
```

---

#### removeHeader

**说明：** 移除头部

**参数：**
- `name` (string) - Header 名称

**示例：**
```json
{"type": "removeHeader", "name": "X-Frame-Options"}
```

---

#### setBody

**说明：** 完全替换 Body 内容

**参数：**
- `value` (string) - Body 内容
- `encoding` (string, 可选) - 编码方式（`text` 或 `base64`），默认 `text`

**示例：**
```json
{
  "type": "setBody",
  "value": "{\"code\": 0, \"data\": {}}",
  "encoding": "text"
}
```

---

#### replaceBodyText

**说明：** 字符串替换 Body 内容

**参数：**
- `search` (string) - 搜索字符串
- `replace` (string) - 替换字符串
- `replaceAll` (boolean, 可选) - 是否全部替换，默认 false

**示例：**
```json
{
  "type": "replaceBodyText",
  "search": "old_value",
  "replace": "new_value",
  "replaceAll": true
}
```

---

#### patchBodyJson

**说明：** 使用 JSON Patch 修改 Body（遵循 RFC 6902 标准）

**参数：**
- `patches` (array) - JSON Patch 操作数组

**示例：**
```json
{
  "type": "patchBodyJson",
  "patches": [
    {"op": "replace", "path": "/user/name", "value": "newName"},
    {"op": "add", "path": "/user/age", "value": 25}
  ]
}
```

---

## JSON Patch 操作详解

`patchBodyJson` 行为支持以下 JSON Patch 操作（RFC 6902 标准）：

### add - 添加

在指定路径添加值（如果路径已存在则替换）

**参数：** `op`, `path`, `value`

**示例：**
```json
{"op": "add", "path": "/user/email", "value": "test@example.com"}
```

---

### remove - 删除

删除指定路径的值

**参数：** `op`, `path`

**示例：**
```json
{"op": "remove", "path": "/user/age"}
```

---

### replace - 替换

替换指定路径的值

**参数：** `op`, `path`, `value`

**示例：**
```json
{"op": "replace", "path": "/user/name", "value": "newName"}
```

---

### move - 移动

将值从一个路径移动到另一个路径

**参数：** `op`, `from`, `path`

**示例：**
```json
{"op": "move", "from": "/user/oldField", "path": "/user/newField"}
```

---

### copy - 复制

将值从一个路径复制到另一个路径

**参数：** `op`, `from`, `path`

**示例：**
```json
{"op": "copy", "from": "/user/name", "path": "/user/displayName"}
```

---

### test - 测试

测试指定路径的值是否等于期望值（常用于条件操作）

**参数：** `op`, `path`, `value`

**示例：**
```json
{"op": "test", "path": "/user/status", "value": "active"}
```

> 💡 **提示**：JSON Path 使用 `/` 分隔路径层级，数组索引使用数字表示，例如 `/users/0/name` 表示 `users` 数组第一个元素的 `name` 字段。

---

## Body 编码方式

在 `setBody` 和 `block` 行为中，可以通过 `encoding` 或 `bodyEncoding` 参数指定 Body 的编码方式：

### text（默认）

纯文本编码，适用于 JSON、HTML、XML 等文本格式。

**示例：**
```json
{
  "type": "setBody",
  "value": "{\"message\": \"Hello\"}",
  "encoding": "text"
}
```

---

### base64

Base64 编码，适用于二进制内容或需要编码传输的场景。

**示例：**
```json
{
  "type": "setBody",
  "value": "eyJtZXNzYWdlIjogIkhlbGxvIn0=",
  "encoding": "base64"
}
```

---

## 完整配置示例

以下是一个包含多条规则的完整配置示例：

```json
{
  "id": "config-20260118-demo",
  "name": "演示配置",
  "version": "1.0",
  "description": "包含常见规则的演示配置",
  "settings": {},
  "rules": [
    {
      "id": "rule-001",
      "name": "Mock 用户信息接口",
      "enabled": true,
      "priority": 10,
      "stage": "response",
      "match": {
        "allOf": [
          {"type": "urlContains", "value": "/api/user/info"}
        ],
        "anyOf": []
      },
      "actions": [
        {
          "type": "setBody",
          "value": "{\"code\":0,\"data\":{\"id\":1,\"name\":\"测试用户\",\"email\":\"test@example.com\"}}",
          "encoding": "text"
        }
      ]
    },
    {
      "id": "rule-002",
      "name": "解决跨域问题",
      "enabled": true,
      "priority": 5,
      "stage": "response",
      "match": {
        "allOf": [],
        "anyOf": []
      },
      "actions": [
        {"type": "setHeader", "name": "Access-Control-Allow-Origin", "value": "*"},
        {"type": "setHeader", "name": "Access-Control-Allow-Methods", "value": "GET, POST, PUT, DELETE"},
        {"type": "setHeader", "name": "Access-Control-Allow-Headers", "value": "Content-Type, Authorization"}
      ]
    },
    {
      "id": "rule-003",
      "name": "修改用户权限",
      "enabled": true,
      "priority": 0,
      "stage": "response",
      "match": {
        "allOf": [
          {"type": "urlContains", "value": "/api/user/profile"}
        ],
        "anyOf": []
      },
      "actions": [
        {
          "type": "patchBodyJson",
          "patches": [
            {"op": "replace", "path": "/data/role", "value": "admin"}
          ]
        }
      ]
    }
  ]
}
```

---

## 实用规则示例

### 示例 1：Mock API 响应

**场景：** 前端开发时，后端接口尚未完成，需要模拟接口返回数据

```json
{
  "id": "rule-mock-api",
  "name": "Mock 商品列表接口",
  "enabled": true,
  "priority": 0,
  "stage": "response",
  "match": {
    "allOf": [
      {"type": "urlContains", "value": "/api/products"},
      {"type": "method", "values": ["GET"]}
    ],
    "anyOf": []
  },
  "actions": [
    {
      "type": "setBody",
      "value": "{\"code\":0,\"data\":[{\"id\":1,\"name\":\"商品A\",\"price\":99.99},{\"id\":2,\"name\":\"商品B\",\"price\":199.99}]}",
      "encoding": "text"
    }
  ]
}
```

---

### 示例 2：精确修改响应数据

**场景：** 使用 JSON Patch 精确修改响应中的某个字段，而不替换整个响应体

```json
{
  "id": "rule-modify-response",
  "name": "修改用户权限字段",
  "enabled": true,
  "priority": 0,
  "stage": "response",
  "match": {
    "allOf": [
      {"type": "urlContains", "value": "/api/user/profile"}
    ],
    "anyOf": []
  },
  "actions": [
    {
      "type": "patchBodyJson",
      "patches": [
        {"op": "replace", "path": "/data/role", "value": "admin"},
        {"op": "add", "path": "/data/permissions", "value": ["read", "write", "delete"]}
      ]
    }
  ]
}
```

---

## 下一步

现在你已经掌握了规则配置的完整语法，可以：

- 📖 查看 [常见问题](./04-faq.md) 了解规则调试技巧
- 🔙 返回 [快速开始](./02-quick-start.md) 创建你的第一个规则
