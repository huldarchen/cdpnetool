package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	api "cdpnetool/pkg/api"
	"cdpnetool/pkg/model"
	"cdpnetool/pkg/rulespec"
)

// Server 提供给 GUI 的 HTTP 接口入口
type Server struct {
	svc api.Service
}

// NewServer 创建 HTTP 接口服务
func NewServer(svc api.Service) *Server {
	return &Server{svc: svc}
}

// ServeHTTP 处理所有 GUI HTTP 请求
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, ErrInvalidRequest.withError(err))
		return
	}
	res := s.dispatch(r.Context(), &req)
	writeResponse(w, res)
}

// Request 表示通用请求结构
type Request struct {
	Method string          `json:"method"`
	ID     string          `json:"id,omitempty"`
	Params json.RawMessage `json:"params"`
}

// Response 表示通用响应结构
type Response struct {
	ID     string       `json:"id,omitempty"`
	Result interface{}  `json:"result,omitempty"`
	Error  *ErrorObject `json:"error,omitempty"`
}

// ErrorObject 表示错误信息
type ErrorObject struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ApiError 表示内部错误类型
type ApiError struct {
	Code string
	Err  error
}

func (e ApiError) withError(err error) ApiError {
	return ApiError{Code: e.Code, Err: err}
}

var (
	// ErrInvalidRequest 无效请求
	ErrInvalidRequest = ApiError{Code: "invalid_request"}
	// ErrMethodNotFound 方法不存在
	ErrMethodNotFound = ApiError{Code: "method_not_found"}
	// ErrInvalidParams 参数错误
	ErrInvalidParams = ApiError{Code: "invalid_params"}
	// ErrInternal 内部错误
	ErrInternal = ApiError{Code: "internal"}
)

// sessionStartParams 会话创建参数
type sessionStartParams struct {
	DevToolsURL       string `json:"devToolsURL"`
	Concurrency       int    `json:"concurrency"`
	BodySizeThreshold int64  `json:"bodySizeThreshold"`
	PendingCapacity   int    `json:"pendingCapacity"`
	ProcessTimeoutMS  int    `json:"processTimeoutMS"`
}

// sessionOnlyParams 仅包含会话标识的参数
type sessionOnlyParams struct {
	SessionID string `json:"sessionId"`
}

// targetAttachParams 目标附加参数
type targetAttachParams struct {
	SessionID string `json:"sessionId"`
	TargetID  string `json:"targetId,omitempty"`
}

// targetDetachParams 目标移除参数
type targetDetachParams struct {
	SessionID string `json:"sessionId"`
	TargetID  string `json:"targetId"`
}

// rulesLoadParams 规则装载参数
type rulesLoadParams struct {
	SessionID string           `json:"sessionId"`
	Rules     rulespec.RuleSet `json:"rules"`
}

// sessionStartResult 会话创建结果
type sessionStartResult struct {
	SessionID string `json:"sessionId"`
}

// targetView 目标视图结构
type targetView struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	Title    string `json:"title"`
	Attached bool   `json:"attached"`
	IsUser   bool   `json:"isUser"`
}

// statsRulesResult 规则统计结果
type statsRulesResult struct {
	Total   int64            `json:"total"`
	Matched int64            `json:"matched"`
	ByRule  map[string]int64 `json:"byRule"`
}

// dispatch 根据 method 分发请求
func (s *Server) dispatch(ctx context.Context, req *Request) *Response {
	var (
		result interface{}
		err    *ErrorObject
	)
	switch req.Method {
	case "session.start":
		result, err = s.handleSessionStart(ctx, req.Params)
	case "session.stop":
		result, err = s.handleSessionStop(ctx, req.Params)
	case "session.enable":
		result, err = s.handleSessionEnable(ctx, req.Params)
	case "session.disable":
		result, err = s.handleSessionDisable(ctx, req.Params)
	case "target.list":
		result, err = s.handleTargetList(ctx, req.Params)
	case "target.attach":
		result, err = s.handleTargetAttach(ctx, req.Params)
	case "target.detach":
		result, err = s.handleTargetDetach(ctx, req.Params)
	case "rules.load":
		result, err = s.handleRulesLoad(ctx, req.Params)
	case "stats.rules":
		result, err = s.handleStatsRules(ctx, req.Params)
	default:
		err = toErrorObject(ErrMethodNotFound)
	}
	return &Response{ID: req.ID, Result: result, Error: err}
}

// writeResponse 写出统一响应
func writeResponse(w http.ResponseWriter, res *Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	_ = enc.Encode(res)
}

// writeError 写出错误响应
func writeError(w http.ResponseWriter, apiErr ApiError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	_ = enc.Encode(&Response{Error: toErrorObject(apiErr)})
}

// toErrorObject 转换错误为响应错误对象
func toErrorObject(e ApiError) *ErrorObject {
	msg := e.Code
	if e.Err != nil {
		msg = e.Err.Error()
	}
	return &ErrorObject{Code: e.Code, Message: msg}
}

// handleSessionStart 处理会话创建
func (s *Server) handleSessionStart(ctx context.Context, params json.RawMessage) (interface{}, *ErrorObject) {
	_ = ctx
	var p sessionStartParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, toErrorObject(ErrInvalidParams.withError(err))
	}
	if p.DevToolsURL == "" {
		return nil, toErrorObject(ErrInvalidParams.withError(errors.New("devToolsURL is required")))
	}
	cfg := model.SessionConfig{
		DevToolsURL:       p.DevToolsURL,
		Concurrency:       defaultInt(p.Concurrency, 4),
		BodySizeThreshold: defaultInt64(p.BodySizeThreshold, 4*1024*1024),
		PendingCapacity:   defaultInt(p.PendingCapacity, 64),
		ProcessTimeoutMS:  defaultInt(p.ProcessTimeoutMS, 200),
	}
	id, err := s.svc.StartSession(cfg)
	if err != nil {
		return nil, toErrorObject(ErrInternal.withError(err))
	}
	return &sessionStartResult{SessionID: string(id)}, nil
}

// handleSessionStop 处理会话停止
func (s *Server) handleSessionStop(ctx context.Context, params json.RawMessage) (interface{}, *ErrorObject) {
	_ = ctx
	var p sessionOnlyParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, toErrorObject(ErrInvalidParams.withError(err))
	}
	if p.SessionID == "" {
		return nil, toErrorObject(ErrInvalidParams.withError(errors.New("sessionId is required")))
	}
	if err := s.svc.StopSession(model.SessionID(p.SessionID)); err != nil {
		return nil, toErrorObject(ErrInternal.withError(err))
	}
	return nil, nil
}

// handleSessionEnable 处理会话启用拦截
func (s *Server) handleSessionEnable(ctx context.Context, params json.RawMessage) (interface{}, *ErrorObject) {
	_ = ctx
	var p sessionOnlyParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, toErrorObject(ErrInvalidParams.withError(err))
	}
	if p.SessionID == "" {
		return nil, toErrorObject(ErrInvalidParams.withError(errors.New("sessionId is required")))
	}
	if err := s.svc.EnableInterception(model.SessionID(p.SessionID)); err != nil {
		return nil, toErrorObject(ErrInternal.withError(err))
	}
	return nil, nil
}

// handleSessionDisable 处理会话停用拦截
func (s *Server) handleSessionDisable(ctx context.Context, params json.RawMessage) (interface{}, *ErrorObject) {
	_ = ctx
	var p sessionOnlyParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, toErrorObject(ErrInvalidParams.withError(err))
	}
	if p.SessionID == "" {
		return nil, toErrorObject(ErrInvalidParams.withError(errors.New("sessionId is required")))
	}
	if err := s.svc.DisableInterception(model.SessionID(p.SessionID)); err != nil {
		return nil, toErrorObject(ErrInternal.withError(err))
	}
	return nil, nil
}

// handleTargetList 处理目标列表查询
func (s *Server) handleTargetList(ctx context.Context, params json.RawMessage) (interface{}, *ErrorObject) {
	_ = ctx
	var p sessionOnlyParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, toErrorObject(ErrInvalidParams.withError(err))
	}
	if p.SessionID == "" {
		return nil, toErrorObject(ErrInvalidParams.withError(errors.New("sessionId is required")))
	}
	targets, err := s.svc.ListTargets(model.SessionID(p.SessionID))
	if err != nil {
		return nil, toErrorObject(ErrInternal.withError(err))
	}
	views := make([]targetView, 0, len(targets))
	for _, t := range targets {
		views = append(views, targetView{
			ID:       string(t.ID),
			Type:     t.Type,
			URL:      t.URL,
			Title:    t.Title,
			Attached: t.IsCurrent,
			IsUser:   t.IsUser,
		})
	}
	return views, nil
}

// handleTargetAttach 处理目标附加
func (s *Server) handleTargetAttach(ctx context.Context, params json.RawMessage) (interface{}, *ErrorObject) {
	_ = ctx
	var p targetAttachParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, toErrorObject(ErrInvalidParams.withError(err))
	}
	if p.SessionID == "" {
		return nil, toErrorObject(ErrInvalidParams.withError(errors.New("sessionId is required")))
	}
	if err := s.svc.AttachTarget(model.SessionID(p.SessionID), model.TargetID(p.TargetID)); err != nil {
		return nil, toErrorObject(ErrInternal.withError(err))
	}
	return nil, nil
}

// handleTargetDetach 处理目标移除
func (s *Server) handleTargetDetach(ctx context.Context, params json.RawMessage) (interface{}, *ErrorObject) {
	_ = ctx
	var p targetDetachParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, toErrorObject(ErrInvalidParams.withError(err))
	}
	if p.SessionID == "" || p.TargetID == "" {
		return nil, toErrorObject(ErrInvalidParams.withError(errors.New("sessionId and targetId are required")))
	}
	if err := s.svc.DetachTarget(model.SessionID(p.SessionID), model.TargetID(p.TargetID)); err != nil {
		return nil, toErrorObject(ErrInternal.withError(err))
	}
	return nil, nil
}

// handleRulesLoad 处理规则装载
func (s *Server) handleRulesLoad(ctx context.Context, params json.RawMessage) (interface{}, *ErrorObject) {
	_ = ctx
	var p rulesLoadParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, toErrorObject(ErrInvalidParams.withError(err))
	}
	if p.SessionID == "" {
		return nil, toErrorObject(ErrInvalidParams.withError(errors.New("sessionId is required")))
	}
	if p.Rules.Version == "" {
		p.Rules.Version = "1.0"
	}
	if err := s.svc.LoadRules(model.SessionID(p.SessionID), p.Rules); err != nil {
		return nil, toErrorObject(ErrInternal.withError(err))
	}
	return nil, nil
}

// handleStatsRules 处理规则统计查询
func (s *Server) handleStatsRules(ctx context.Context, params json.RawMessage) (interface{}, *ErrorObject) {
	_ = ctx
	var p sessionOnlyParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, toErrorObject(ErrInvalidParams.withError(err))
	}
	if p.SessionID == "" {
		return nil, toErrorObject(ErrInvalidParams.withError(errors.New("sessionId is required")))
	}
	st, err := s.svc.GetRuleStats(model.SessionID(p.SessionID))
	if err != nil {
		return nil, toErrorObject(ErrInternal.withError(err))
	}
	res := statsRulesResult{
		Total:   st.Total,
		Matched: st.Matched,
		ByRule:  make(map[string]int64, len(st.ByRule)),
	}
	for k, v := range st.ByRule {
		res.ByRule[string(k)] = v
	}
	return res, nil
}

// defaultInt 整型默认值
func defaultInt(v, d int) int {
	if v == 0 {
		return d
	}
	return v
}

// defaultInt64 整型默认值（int64）
func defaultInt64(v, d int64) int64 {
	if v == 0 {
		return d
	}
	return v
}
