package model

type SessionID string
type TargetID string
type RuleID string

type SessionConfig struct {
	DevToolsURL       string `json:"devToolsURL"`
	Concurrency       int    `json:"concurrency"`
	BodySizeThreshold int64  `json:"bodySizeThreshold"`
	PendingCapacity   int    `json:"pendingCapacity"`
	ProcessTimeoutMS  int    `json:"processTimeoutMS"`
}

type RuleSet struct {
	Version string `json:"version"`
	Rules   []Rule `json:"rules"`
}

type Rule struct {
	ID       RuleID `json:"id"`
	Priority int    `json:"priority"`
	Mode     string `json:"mode"`
	Match    Match  `json:"match"`
	Action   Action `json:"action"`
}

type Match struct {
	AllOf  []Condition `json:"allOf"`
	AnyOf  []Condition `json:"anyOf"`
	NoneOf []Condition `json:"noneOf"`
}

type Condition struct {
	Type    string   `json:"type"`
	Mode    string   `json:"mode"`
	Pattern string   `json:"pattern"`
	Values  []string `json:"values"`
	Key     string   `json:"key"`
	Op      string   `json:"op"`
	Value   string   `json:"value"`
	Pointer string   `json:"pointer"`
}

type Action struct {
	Rewrite  *Rewrite `json:"rewrite"`
	Respond  *Respond `json:"respond"`
	Fail     *Fail    `json:"fail"`
	DelayMS  int      `json:"delayMS"`
	DropRate float64  `json:"dropRate"`
	Pause    *Pause   `json:"pause"`
}

type Rewrite struct {
	URL     *string            `json:"url"`
	Method  *string            `json:"method"`
	Headers map[string]*string `json:"headers"`
	Query   map[string]*string `json:"query"`
	Cookies map[string]*string `json:"cookies"`
	Body    *BodyPatch         `json:"body"`
}

type BodyPatch struct {
	Type string `json:"type"`
	Ops  []any  `json:"ops"`
}

type Respond struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
	Base64  bool              `json:"base64"`
}

type Fail struct {
	Reason string `json:"reason"`
}

type Pause struct {
	Stage         string `json:"stage"`
	TimeoutMS     int    `json:"timeoutMS"`
	DefaultAction struct {
		Type   string `json:"type"`
		Status int    `json:"status"`
		Reason string `json:"reason"`
	} `json:"defaultAction"`
}

type EngineStats struct {
	Total   int64            `json:"total"`
	Matched int64            `json:"matched"`
	ByRule  map[RuleID]int64 `json:"byRule"`
}

type Event struct {
	Type    string    `json:"type"`
	Session SessionID `json:"session"`
	Target  TargetID  `json:"target"`
	Rule    *RuleID   `json:"rule"`
	Error   error     `json:"error"`
}
