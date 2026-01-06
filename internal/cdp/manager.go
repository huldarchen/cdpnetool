package cdp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ilog "cdpnetool/internal/log"
	"cdpnetool/internal/rules"
	"cdpnetool/pkg/model"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/fetch"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/rpcc"
)

type Manager struct {
	devtoolsURL string
	conn        *rpcc.Conn
	client      *cdp.Client
	ctx         context.Context
	cancel      context.CancelFunc
	events      chan model.Event
	pending     chan any
	engine      *rules.Engine
	approvals   map[string]chan model.Rewrite
}

func New(devtoolsURL string, events chan model.Event, pending chan any) *Manager {
	return &Manager{devtoolsURL: devtoolsURL, events: events, pending: pending, approvals: make(map[string]chan model.Rewrite)}
}

func (m *Manager) AttachTarget(target model.TargetID) error {
	ctx, cancel := context.WithCancel(context.Background())
	m.ctx = ctx
	m.cancel = cancel
	dt := devtool.New(m.devtoolsURL)
	targets, err := dt.List(ctx)
	if err != nil {
		return err
	}
	var sel *devtool.Target
	for i := range targets {
		if string(targets[i].ID) == string(target) || target == "" {
			sel = targets[i]
			if target == "" {
				break
			}
		}
	}
	if sel == nil {
		return fmt.Errorf("no target")
	}
	conn, err := rpcc.DialContext(ctx, sel.WebSocketDebuggerURL)
	if err != nil {
		return err
	}
	m.conn = conn
	m.client = cdp.NewClient(conn)
	return nil
}

func (m *Manager) Detach() error {
	if m.cancel != nil {
		m.cancel()
	}
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}

func (m *Manager) Enable() error {
	if m.client == nil {
		return fmt.Errorf("not attached")
	}
	err := m.client.Network.Enable(m.ctx, nil)
	if err != nil {
		return err
	}
	p := "*"
	patterns := []fetch.RequestPattern{
		{URLPattern: &p, RequestStage: fetch.RequestStageRequest},
		{URLPattern: &p, RequestStage: fetch.RequestStageResponse},
	}
	err = m.client.Fetch.Enable(m.ctx, &fetch.EnableArgs{Patterns: patterns})
	if err != nil {
		return err
	}
	go m.consume()
	return nil
}

func (m *Manager) Disable() error {
	if m.client == nil {
		return fmt.Errorf("not attached")
	}
	return m.client.Fetch.Disable(m.ctx)
}

func (m *Manager) consume() {
	rp, err := m.client.Fetch.RequestPaused(m.ctx)
	if err != nil {
		return
	}
	defer rp.Close()
	for {
		ev, err := rp.Recv()
		if err != nil {
			return
		}
		m.handle(ev)
	}
}

func (m *Manager) handle(ev *fetch.RequestPausedReply) {
	ctx, cancel := context.WithTimeout(m.ctx, 3*time.Second)
	defer cancel()
	m.events <- model.Event{Type: "intercepted"}
	stg := "request"
	if ev.ResponseStatusCode != nil {
		stg = "response"
	}
	a := m.decide(ev, stg)
	if a == nil {
		m.applyContinue(ctx, ev, stg)
		return
	}
	if a.DelayMS > 0 {
		time.Sleep(time.Duration(a.DelayMS) * time.Millisecond)
	}
	if a.Pause != nil {
		m.applyPause(ctx, ev, a.Pause, stg)
		return
	}
	if a.Fail != nil {
		m.applyFail(ctx, ev, a.Fail)
		return
	}
	if a.Respond != nil {
		m.applyRespond(ctx, ev, a.Respond)
		return
	}
	if a.Rewrite != nil {
		m.applyRewrite(ctx, ev, a.Rewrite)
		return
	}
	m.applyContinue(ctx, ev, stg)
}

func (m *Manager) decide(ev *fetch.RequestPausedReply, stage string) *model.Action {
	if m.engine == nil {
		return nil
	}
	h := map[string]string{}
	_ = json.Unmarshal(ev.Request.Headers, &h)
	res := m.engine.Eval(rules.Ctx{URL: ev.Request.URL, Method: ev.Request.Method, Headers: h, Stage: stage})
	if res == nil {
		return nil
	}
	return res.Action
}

func (m *Manager) applyContinue(ctx context.Context, ev *fetch.RequestPausedReply, stage string) {
	if stage == "response" {
		m.client.Fetch.ContinueResponse(ctx, &fetch.ContinueResponseArgs{RequestID: ev.RequestID})
		ilog.L().Debug("continue_response")
	} else {
		m.client.Fetch.ContinueRequest(ctx, &fetch.ContinueRequestArgs{RequestID: ev.RequestID})
		ilog.L().Debug("continue_request")
	}
}

func (m *Manager) applyFail(ctx context.Context, ev *fetch.RequestPausedReply, f *model.Fail) {
	m.client.Fetch.FailRequest(ctx, &fetch.FailRequestArgs{RequestID: ev.RequestID, ErrorReason: network.ErrorReasonFailed})
	m.events <- model.Event{Type: "failed"}
}

func (m *Manager) applyRespond(ctx context.Context, ev *fetch.RequestPausedReply, r *model.Respond) {
	args := &fetch.FulfillRequestArgs{RequestID: ev.RequestID, ResponseCode: r.Status}
	if len(r.Headers) > 0 {
		args.ResponseHeaders = toHeaderEntries(r.Headers)
	}
	if len(r.Body) > 0 {
		args.Body = r.Body
	}
	m.client.Fetch.FulfillRequest(ctx, args)
	m.events <- model.Event{Type: "fulfilled"}
}

func (m *Manager) applyRewrite(ctx context.Context, ev *fetch.RequestPausedReply, rw *model.Rewrite) {
	var url, method *string
	if rw.URL != nil {
		url = rw.URL
	}
	if rw.Method != nil {
		method = rw.Method
	}
	var hdrs []fetch.HeaderEntry
	if rw.Headers != nil {
		for k, v := range rw.Headers {
			if v != nil {
				hdrs = append(hdrs, fetch.HeaderEntry{Name: k, Value: *v})
			}
		}
	}
	m.client.Fetch.ContinueRequest(ctx, &fetch.ContinueRequestArgs{RequestID: ev.RequestID, URL: url, Method: method, Headers: hdrs})
	m.events <- model.Event{Type: "mutated"}
}

func toHeaderEntries(h map[string]string) []fetch.HeaderEntry {
	out := make([]fetch.HeaderEntry, 0, len(h))
	for k, v := range h {
		out = append(out, fetch.HeaderEntry{Name: k, Value: v})
	}
	return out
}

func (m *Manager) applyPause(ctx context.Context, ev *fetch.RequestPausedReply, p *model.Pause, stage string) {
	id := string(ev.RequestID)
	ch := make(chan model.Rewrite, 1)
	m.approvals[id] = ch
	if m.pending != nil {
		m.pending <- struct{ ID string }{ID: id}
	}
	t := time.NewTimer(time.Duration(p.TimeoutMS) * time.Millisecond)
	select {
	case mut := <-ch:
		_ = mut
		m.applyContinue(ctx, ev, stage)
	case <-t.C:
		switch p.DefaultAction.Type {
		case "fulfill":
			m.applyRespond(ctx, ev, &model.Respond{Status: p.DefaultAction.Status})
		case "fail":
			m.applyFail(ctx, ev, &model.Fail{Reason: p.DefaultAction.Reason})
		case "continue_mutated":
			m.applyContinue(ctx, ev, stage)
		default:
			m.applyContinue(ctx, ev, stage)
		}
	}
	delete(m.approvals, id)
}

func (m *Manager) SetRules(rs model.RuleSet) { m.engine = rules.New(rs) }

func (m *Manager) UpdateRules(rs model.RuleSet) {
	if m.engine == nil {
		m.engine = rules.New(rs)
	} else {
		m.engine.Update(rs)
	}
}

func (m *Manager) Approve(itemID string, mutations model.Rewrite) {
	if ch, ok := m.approvals[itemID]; ok {
		ch <- mutations
	}
}
