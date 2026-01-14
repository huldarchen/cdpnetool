package cdp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/mafredri/cdp/protocol/fetch"
	"github.com/mafredri/cdp/protocol/network"

	"cdpnetool/pkg/rulespec"
)

// applyContinue 继续原请求或响应不做修改
func (m *Manager) applyContinue(ctx context.Context, ts *targetSession, ev *fetch.RequestPausedReply, stage string) {
	if ts == nil || ts.client == nil {
		return
	}
	if stage == "response" {
		_ = ts.client.Fetch.ContinueResponse(ctx, &fetch.ContinueResponseArgs{RequestID: ev.RequestID})
		m.log.Debug("继续原始响应")
	} else {
		_ = ts.client.Fetch.ContinueRequest(ctx, &fetch.ContinueRequestArgs{RequestID: ev.RequestID})
		m.log.Debug("继续原始请求")
	}
}

// applyFail 使请求失败并返回错误原因
func (m *Manager) applyFail(ctx context.Context, ts *targetSession, ev *fetch.RequestPausedReply, f *rulespec.Fail) {
	if ts == nil || ts.client == nil || f == nil {
		return
	}
	_ = ts.client.Fetch.FailRequest(ctx, &fetch.FailRequestArgs{RequestID: ev.RequestID, ErrorReason: network.ErrorReason(f.Reason)})
}

// applyRespond 返回自定义响应（可只改头或完整替换）
func (m *Manager) applyRespond(ctx context.Context, ts *targetSession, ev *fetch.RequestPausedReply, r *rulespec.Respond, stage string) {
	if ts == nil || ts.client == nil || r == nil {
		return
	}
	if stage == "response" && len(r.Body) == 0 {
		// 仅修改响应码/头，继续响应
		m.continueResponseWithModifications(ctx, ts, ev, r)
		return
	}
	// fulfill 完整响应
	m.fulfillRequest(ctx, ts, ev, r)
}

// continueResponseWithModifications 继续响应并修改状态码/头部
func (m *Manager) continueResponseWithModifications(ctx context.Context, ts *targetSession, ev *fetch.RequestPausedReply, r *rulespec.Respond) {
	args := &fetch.ContinueResponseArgs{RequestID: ev.RequestID}
	if r.Status != 0 {
		args.ResponseCode = &r.Status
	}
	if len(r.Headers) > 0 {
		args.ResponseHeaders = toHeaderEntries(r.Headers)
	}
	_ = ts.client.Fetch.ContinueResponse(ctx, args)
}

// fulfillRequest 完整响应请求
func (m *Manager) fulfillRequest(ctx context.Context, ts *targetSession, ev *fetch.RequestPausedReply, r *rulespec.Respond) {
	args := &fetch.FulfillRequestArgs{RequestID: ev.RequestID, ResponseCode: r.Status}
	if len(r.Headers) > 0 {
		args.ResponseHeaders = toHeaderEntries(r.Headers)
	}
	if len(r.Body) > 0 {
		args.Body = r.Body
	}
	_ = ts.client.Fetch.FulfillRequest(ctx, args)
}

// applyRewrite 根据规则对请求或响应进行重写
func (m *Manager) applyRewrite(ctx context.Context, ts *targetSession, ev *fetch.RequestPausedReply, rw *rulespec.Rewrite, stage string) {
	if ts == nil || ts.client == nil || rw == nil {
		return
	}
	if stage == "response" {
		m.applyResponseRewrite(ctx, ts, ev, rw)
	} else {
		m.applyRequestRewrite(ctx, ts, ev, rw)
	}
}

// applyResponseRewrite 处理响应阶段的重写
func (m *Manager) applyResponseRewrite(ctx context.Context, ts *targetSession, ev *fetch.RequestPausedReply, rw *rulespec.Rewrite) {
	if rw.Body == nil {
		// 仅修改头部，不需要获取 Body
		if rw.Headers != nil {
			cur := m.getCurrentResponseHeaders(ev)
			cur = applyHeaderPatch(cur, rw.Headers)
			_ = ts.client.Fetch.ContinueResponse(ctx, &fetch.ContinueResponseArgs{RequestID: ev.RequestID, ResponseHeaders: toHeaderEntries(cur)})
			return
		}
		_ = ts.client.Fetch.ContinueResponse(ctx, &fetch.ContinueResponseArgs{RequestID: ev.RequestID})
		return
	}

	// 需要修改 Body
	ctype, clen := m.extractResponseMetadata(ev)
	if !shouldGetBody(ctype, clen, m.bodySizeThreshold) {
		_ = ts.client.Fetch.ContinueResponse(ctx, &fetch.ContinueResponseArgs{RequestID: ev.RequestID})
		return
	}

	bodyText, ok := m.fetchResponseBody(ctx, ts, ev.RequestID)
	if !ok {
		_ = ts.client.Fetch.ContinueResponse(ctx, &fetch.ContinueResponseArgs{RequestID: ev.RequestID})
		return
	}

	newBody, ok := applyBodyPatch(bodyText, rw.Body)
	if !ok || len(newBody) == 0 {
		_ = ts.client.Fetch.ContinueResponse(ctx, &fetch.ContinueResponseArgs{RequestID: ev.RequestID})
		return
	}

	code := 200
	if ev.ResponseStatusCode != nil {
		code = *ev.ResponseStatusCode
	}
	cur := m.getCurrentResponseHeaders(ev)
	cur = applyHeaderPatch(cur, rw.Headers)
	args := &fetch.FulfillRequestArgs{
		RequestID:       ev.RequestID,
		ResponseCode:    code,
		ResponseHeaders: toHeaderEntries(cur),
		Body:            newBody,
	}
	_ = ts.client.Fetch.FulfillRequest(ctx, args)
}

// getCurrentResponseHeaders 获取当前响应头部映射
func (m *Manager) getCurrentResponseHeaders(ev *fetch.RequestPausedReply) map[string]string {
	cur := make(map[string]string, len(ev.ResponseHeaders))
	for i := range ev.ResponseHeaders {
		cur[strings.ToLower(ev.ResponseHeaders[i].Name)] = ev.ResponseHeaders[i].Value
	}
	return cur
}

// extractResponseMetadata 提取响应元数据（Content-Type, Content-Length）
func (m *Manager) extractResponseMetadata(ev *fetch.RequestPausedReply) (ctype string, clen int64) {
	for i := range ev.ResponseHeaders {
		k := ev.ResponseHeaders[i].Name
		v := ev.ResponseHeaders[i].Value
		if strings.EqualFold(k, "content-type") {
			ctype = v
		}
		if strings.EqualFold(k, "content-length") {
			if n, err := parseInt64(v); err == nil {
				clen = n
			}
		}
	}
	return
}

// fetchResponseBody 获取响应 Body 文本
func (m *Manager) fetchResponseBody(ctx context.Context, ts *targetSession, requestID fetch.RequestID) (string, bool) {
	if ts == nil || ts.client == nil {
		return "", false
	}
	ctx2, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	rb, err := ts.client.Fetch.GetResponseBody(ctx2, &fetch.GetResponseBodyArgs{RequestID: requestID})
	if err != nil || rb == nil {
		return "", false
	}
	if rb.Base64Encoded {
		if b, err := base64.StdEncoding.DecodeString(rb.Body); err == nil {
			return string(b), true
		}
		return "", false
	}
	return rb.Body, true
}

// applyRequestRewrite 处理请求阶段的重写
func (m *Manager) applyRequestRewrite(ctx context.Context, ts *targetSession, ev *fetch.RequestPausedReply, rw *rulespec.Rewrite) {
	if ts == nil || ts.client == nil || rw == nil {
		return
	}

	var url, method *string
	if rw.URL != nil {
		url = rw.URL
	}
	if rw.Method != nil {
		method = rw.Method
	}

	hdrs := m.buildRequestHeaders(rw, ev)
	post := m.buildRequestBody(rw, ev)

	args := &fetch.ContinueRequestArgs{
		RequestID: ev.RequestID,
		URL:       url,
		Method:    method,
		Headers:   hdrs,
	}

	if rw.Query != nil && url == nil {
		if u, err := urlParse(ev.Request.URL, rw.Query); err == nil {
			us := u.String()
			args.URL = &us
		}
	}

	if len(post) > 0 {
		args.PostData = post
	}

	_ = ts.client.Fetch.ContinueRequest(ctx, args)
}

// buildRequestHeaders 构建请求头部列表
func (m *Manager) buildRequestHeaders(rw *rulespec.Rewrite, ev *fetch.RequestPausedReply) []fetch.HeaderEntry {
	var hdrs []fetch.HeaderEntry
	if rw.Headers != nil {
		for k, v := range rw.Headers {
			if v != nil {
				hdrs = append(hdrs, fetch.HeaderEntry{Name: k, Value: *v})
			}
		}
	}

	if rw.Cookies != nil {
		h := map[string]string{}
		_ = json.Unmarshal(ev.Request.Headers, &h)
		var cookie string
		for k, v := range h {
			if strings.EqualFold(k, "cookie") {
				cookie = v
				break
			}
		}
		cm := parseCookie(cookie)
		for name, val := range rw.Cookies {
			if val == nil {
				delete(cm, name)
			} else {
				cm[name] = *val
			}
		}
		if len(cm) > 0 {
			var b strings.Builder
			first := true
			for k, v := range cm {
				if !first {
					b.WriteString("; ")
				}
				first = false
				b.WriteString(k)
				b.WriteString("=")
				b.WriteString(v)
			}
			hdrs = append(hdrs, fetch.HeaderEntry{Name: "Cookie", Value: b.String()})
		}
	}

	return hdrs
}

// buildRequestBody 构建请求 Body
func (m *Manager) buildRequestBody(rw *rulespec.Rewrite, ev *fetch.RequestPausedReply) []byte {
	if rw.Body == nil {
		return nil
	}
	var src string
	// 优先使用 PostDataEntries（新 API）
	if len(ev.Request.PostDataEntries) > 0 {
		for _, entry := range ev.Request.PostDataEntries {
			if entry.Bytes != nil {
				src += *entry.Bytes
			}
		}
	} else if ev.Request.PostData != nil {
		// 向下兼容，使用已弃用的 PostData
		src = *ev.Request.PostData
	}
	if b, ok := applyBodyPatch(src, rw.Body); ok && len(b) > 0 {
		return b
	}
	return nil
}

// toHeaderEntries 将头部映射转换为 CDP 头部条目
func toHeaderEntries(h map[string]string) []fetch.HeaderEntry {
	out := make([]fetch.HeaderEntry, 0, len(h))
	for k, v := range h {
		out = append(out, fetch.HeaderEntry{Name: k, Value: v})
	}
	return out
}
