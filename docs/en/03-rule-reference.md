# Rule Configuration Reference

This section provides detailed information about cdpnetool's rule configuration format, all supported matching conditions, and action types.

---

## Configuration File Structure

### Config Root Object

The configuration root object contains metadata and a list of rules:

```json
{
  "id": "config-20260118-abc123",
  "name": "My Configuration",
  "version": "1.0",
  "description": "Configuration description",
  "settings": {},
  "rules": []
}
```

**Field Description:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique configuration identifier, format: `config-YYYYMMDD-random6` |
| `name` | string | Yes | Configuration name |
| `version` | string | Yes | Configuration version (currently 1.0) |
| `description` | string | No | Configuration description |
| `settings` | object | No | Reserved settings |
| `rules` | array | Yes | Array of rules |

---

### Rule Object

Each rule contains matching conditions and actions:

```json
{
  "id": "rule-001",
  "name": "Rule Name",
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

**Field Description:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique rule identifier, format: `rule-XXX` |
| `name` | string | Yes | Rule name |
| `enabled` | boolean | Yes | Whether enabled |
| `priority` | number | Yes | Priority, higher values execute first |
| `stage` | string | Yes | Lifecycle stage (`request` or `response`) |
| `match` | object | Yes | Match condition object |
| `actions` | array | Yes | Array of actions |

---

## Lifecycle Stages

### request - Request Stage

- Intercepts before request is sent to server
- Can modify request URL, method, headers, body
- Can block request from being sent (block action)

**Use Cases:**
- Modify request parameters
- Add or modify request headers
- Block specific requests

---

### response - Response Stage

- Intercepts before response returns to browser
- Can modify response status code, headers, body
- Cannot block request (already sent)

**Use Cases:**
- Mock API response data
- Modify response content
- Resolve CORS issues

---

## Match Conditions Reference

### Condition Logic Combination

Match conditions use `allOf` and `anyOf` for logical combination:

**AllOf (AND Logic):**
- All conditions must be satisfied to match
- Array of conditions, evaluated in order

**AnyOf (OR Logic):**
- Any one condition satisfied triggers match
- Array of conditions, evaluated in order

**Example:**
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

> 💡 **Rule**: When both `allOf` and `anyOf` exist, all conditions in `allOf` must be satisfied AND at least one condition in `anyOf` must be satisfied.

---

## URL Condition Types

| Condition Type | Description | Parameters | Example Value |
|---------------|-------------|------------|---------------|
| `urlEquals` | Exact URL match | `value` (string) | `"https://example.com/api/user"` |
| `urlPrefix` | URL prefix match | `value` (string) | `"https://example.com/api/"` |
| `urlSuffix` | URL suffix match | `value` (string) | `".json"` |
| `urlContains` | URL contains string | `value` (string) | `"/api/user"` |
| `urlRegex` | URL regex match | `pattern` (string) | `"^https://example\\.com/api/(user|order)/\\d+$"` |

---

## HTTP Property Conditions

### method

**Description:** HTTP method match

**Parameters:**
- `values` (string[]) - Array of HTTP methods

**Supported Values:** `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `HEAD`, `OPTIONS`

**Example:**
```json
{"type": "method", "values": ["GET", "POST"]}
```

---

### resourceType

**Description:** Resource type match

**Parameters:**
- `values` (string[]) - Array of resource types

**Supported Values:**

| Value | Description |
|-------|-------------|
| `document` | HTML document |
| `script` | JavaScript |
| `stylesheet` | CSS |
| `image` | Images |
| `media` | Audio/Video |
| `font` | Fonts |
| `xhr` | XMLHttpRequest |
| `fetch` | Fetch API |
| `websocket` | WebSocket |
| `other` | Other |

**Example:**
```json
{"type": "resourceType", "values": ["xhr", "fetch"]}
```

---

## Header Condition Types

| Condition Type | Description | Parameters | Example |
|---------------|-------------|------------|---------|
| `headerExists` | Check if header exists | `name` (string) | `{"type": "headerExists", "name": "Authorization"}` |
| `headerNotExists` | Check if header not exists | `name` (string) | `{"type": "headerNotExists", "name": "X-Custom"}` |
| `headerEquals` | Header value exact match | `name`, `value` (string) | `{"type": "headerEquals", "name": "Content-Type", "value": "application/json"}` |
| `headerContains` | Header value contains string | `name`, `value` (string) | `{"type": "headerContains", "name": "User-Agent", "value": "Chrome"}` |
| `headerRegex` | Header value regex match | `name`, `pattern` (string) | `{"type": "headerRegex", "name": "Authorization", "pattern": "^Bearer\\s+[A-Za-z0-9\\-_]+$"}` |

---

## Query Parameter Condition Types

Query parameter conditions are similar to header conditions, used for matching URL query parameters:

- `queryExists` - Query parameter exists
- `queryNotExists` - Query parameter not exists
- `queryEquals` - Query parameter value exact match
- `queryContains` - Query parameter value contains match
- `queryRegex` - Query parameter value regex match

**Parameters:** `name` - Parameter name, other parameters same as header conditions

**Example:**
```json
{"type": "queryEquals", "name": "page", "value": "1"}
```

---

## Cookie Condition Types

Cookie conditions are similar to header conditions, used for matching HTTP Cookies:

- `cookieExists` - Cookie exists
- `cookieNotExists` - Cookie not exists
- `cookieEquals` - Cookie value exact match
- `cookieContains` - Cookie value contains match
- `cookieRegex` - Cookie value regex match

**Parameters:** `name` - Cookie name, other parameters same as header conditions

**Example:**
```json
{"type": "cookieExists", "name": "sessionId"}
```

---

## Body Condition Types

### bodyContains

**Description:** Body contains specified string

**Parameters:**
- `value` (string) - String to contain

**Example:**
```json
{"type": "bodyContains", "value": "username"}
```

---

### bodyRegex

**Description:** Body regex match

**Parameters:**
- `pattern` (string) - Regular expression

**Example:**
```json
{"type": "bodyRegex", "pattern": "\"userId\":\\s*\\d+"}
```

---

### bodyJsonPath

**Description:** Match body content using JSON Path

**Parameters:**
- `path` (string) - JSON Path expression
- `value` (string) - Expected value

**Example:**
```json
{"type": "bodyJsonPath", "path": "$.user.id", "value": "123"}
```

---

## Actions Reference

### Request Stage Only Actions

The following actions are only available when `stage: "request"`:

| Action Type | Description | Parameters | Example |
|-------------|-------------|------------|---------|
| `setUrl` | Set request URL | `value` (string) | `{"type": "setUrl", "value": "https://example.com/api/v2/user"}` |
| `setMethod` | Set request method | `value` (string) | `{"type": "setMethod", "value": "POST"}` |
| `setQueryParam` | Set URL query parameter | `name`, `value` | `{"type": "setQueryParam", "name": "page", "value": "1"}` |
| `removeQueryParam` | Remove URL query parameter | `name` (string) | `{"type": "removeQueryParam", "name": "debug"}` |
| `setCookie` | Set Cookie | `name`, `value` | `{"type": "setCookie", "name": "token", "value": "abc123"}` |
| `removeCookie` | Remove Cookie | `name` (string) | `{"type": "removeCookie", "name": "tracking_id"}` |
| `setFormField` | Set form field | `name`, `value` | `{"type": "setFormField", "name": "username", "value": "test"}` |
| `removeFormField` | Remove form field | `name` (string) | `{"type": "removeFormField", "name": "csrf_token"}` |

---

### block Action

**Description:** Block request and return custom response (terminal action, subsequent actions not executed)

**Parameters:**
- `statusCode` (number, optional) - HTTP status code, default 200
- `headers` (object, optional) - Response headers object
- `body` (string, optional) - Response body string
- `bodyEncoding` (string, optional) - Body encoding (`text` or `base64`), default `text`

**Example:**
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

### Response Stage Only Actions

The following actions are only available when `stage: "response"`:

| Action Type | Description | Parameters | Example |
|-------------|-------------|------------|---------|
| `setStatus` | Set response status code | `value` (number) | `{"type": "setStatus", "value": 200}` |

---

### Common Actions (Available in Both Stages)

The following actions can be used in both stages:

| Action Type | Description | Parameters | Example |
|-------------|-------------|------------|---------|
| `setHeader` | Set header | `name`, `value` | `{"type": "setHeader", "name": "X-Custom-Header", "value": "custom-value"}` |
| `removeHeader` | Remove header | `name` (string) | `{"type": "removeHeader", "name": "X-Frame-Options"}` |
| `setBody` | Completely replace body | `value` (string), `encoding` (optional) | `{"type": "setBody", "value": "{\"code\": 0}", "encoding": "text"}` |
| `replaceBodyText` | String replace body content | `search`, `replace`, `replaceAll` (optional) | `{"type": "replaceBodyText", "search": "old", "replace": "new", "replaceAll": true}` |
| `patchBodyJson` | Modify body using JSON Patch | `patches` (array) | See JSON Patch section below |

---

## JSON Patch Operations

The `patchBodyJson` action supports the following JSON Patch operations (RFC 6902 standard):

| Operation | Description | Parameters | Example |
|-----------|-------------|------------|---------|
| `add` | Add value at path | `op`, `path`, `value` | `{"op": "add", "path": "/user/email", "value": "test@example.com"}` |
| `remove` | Remove value at path | `op`, `path` | `{"op": "remove", "path": "/user/age"}` |
| `replace` | Replace value at path | `op`, `path`, `value` | `{"op": "replace", "path": "/user/name", "value": "newName"}` |
| `move` | Move value from one path to another | `op`, `from`, `path` | `{"op": "move", "from": "/user/oldField", "path": "/user/newField"}` |
| `copy` | Copy value from one path to another | `op`, `from`, `path` | `{"op": "copy", "from": "/user/name", "path": "/user/displayName"}` |
| `test` | Test if value at path equals expected value | `op`, `path`, `value` | `{"op": "test", "path": "/user/status", "value": "active"}` |

> 💡 **Tip**: JSON Path uses `/` to separate path hierarchy, array indices use numbers, e.g., `/users/0/name` represents the `name` field of the first element in the `users` array.

---

## Body Encoding Methods

In `setBody` and `block` actions, you can specify body encoding method via `encoding` or `bodyEncoding` parameter:

### text (Default)

Plain text encoding, suitable for JSON, HTML, XML and other text formats.

**Example:**
```json
{
  "type": "setBody",
  "value": "{\"message\": \"Hello\"}",
  "encoding": "text"
}
```

---

### base64

Base64 encoding, suitable for binary content or scenarios requiring encoded transmission.

**Example:**
```json
{
  "type": "setBody",
  "value": "eyJtZXNzYWdlIjogIkhlbGxvIn0=",
  "encoding": "base64"
}
```

---

## Complete Configuration Example

The following is a complete configuration example containing multiple rules:

```json
{
  "id": "config-20260118-demo",
  "name": "Demo Configuration",
  "version": "1.0",
  "description": "Demo configuration with common rules",
  "settings": {},
  "rules": [
    {
      "id": "rule-001",
      "name": "Mock User Info API",
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
          "value": "{\"code\":0,\"data\":{\"id\":1,\"name\":\"Test User\",\"email\":\"test@example.com\"}}",
          "encoding": "text"
        }
      ]
    },
    {
      "id": "rule-002",
      "name": "Resolve CORS Issue",
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
      "name": "Modify User Role",
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

## Practical Rule Examples

### Example 1: Mock API Response

**Scenario:** Frontend development when backend API not yet completed, need to mock interface return data

```json
{
  "id": "rule-mock-api",
  "name": "Mock Product List API",
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
      "value": "{\"code\":0,\"data\":[{\"id\":1,\"name\":\"Product A\",\"price\":99.99},{\"id\":2,\"name\":\"Product B\",\"price\":199.99}]}",
      "encoding": "text"
    }
  ]
}
```

---

### Example 2: Precisely Modify Response Data

**Scenario:** Use JSON Patch to precisely modify a field in response without replacing entire response body

```json
{
  "id": "rule-modify-response",
  "name": "Modify User Role Field",
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

## Next Steps

Now that you have mastered the complete rule configuration syntax, you can:

- 📖 View [FAQ](./04-faq.md) to learn rule debugging tips
- 🔙 Return to [Quick Start](./02-quick-start.md) to create your first rule
