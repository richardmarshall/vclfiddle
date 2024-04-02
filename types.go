package fiddle

import "encoding/json"

func PrettyPrint(v any) string {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return ""
	}
	return string(b)
}

type Config struct {
	Authorizations map[string]string `json:"authorizations"`
}

type Resp struct {
	Fiddle     Fiddle                `json:"fiddle"`
	Valid      bool                  `json:"valid"`
	LintStatus map[string]LintStatus `json:"lintStatus"`
}

type LintStatus struct {
	NoContent bool
	Lints     []Lint
}

func (ls *LintStatus) UnmarshalJSON(data []byte) error {
	if string(data) == "true" {
		ls.NoContent = true
		return nil
	}
	if err := json.Unmarshal(data, &ls.Lints); err != nil {
		return err
	}
	return nil
}

func (ls LintStatus) MarshalJSON() ([]byte, error) {
	if ls.NoContent {
		return json.Marshal(ls.NoContent)
	}
	return json.Marshal(ls.Lints)
}

type Lint struct {
	Level    string `json:"level"`
	Str      string `json:"str"`
	Line     int    `json:"line"`
	StartPos int    `json:"startPos"`
	EndPos   int    `json:"endPos"`
	Message  string `json:"message"`
}

type ExecuteResp struct {
	SessionID  string `json:"sessionId"`
	StreamHost string `json:"streamHost"`
}

type Fiddle struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Type          string    `json:"type"`
	Description   string    `json:"description"`
	Origins       []string  `json:"origins"`
	Source        Source    `json:"src"`
	SourceVersion int       `json:"srcVersion"`
	Requests      []Request `json:"requests"`
	IsLocked      bool      `json:"isLocked"`
	LastEditor    string    `json:"lastEditor"`
}

type Request struct {
	Method          string         `json:"method"`
	Path            string         `json:"path"`
	Headers         string         `json:"headers"`
	Body            string         `json:"body"`
	Data            map[string]any `json:"data"`
	EnableCluster   bool           `json:"enableCluster"`
	EnableShield    bool           `json:"enableShield"`
	UseFreshCache   bool           `json:"useFreshCache"`
	ConnType        string         `json:"connType"`
	SourceIp        string         `json:"sourceIp"`
	FollowRedirects bool           `json:"followRedirects"`
	Tests           string         `json:"tests"`
	Delay           int            `json:"delay"`
}

type Source struct {
	Init string `json:"init"`
	Recv string `json:"recv"`
}

type ExecResults struct {
	ClientFetches map[string]Fetch `json:"clientFetches"`
	Events        []Event          `json:"events"`
	ExecHost      string           `json:"execHost"`
	ExecVersion   int              `json:"execVersion"`
	ID            string           `json:"id"`
	Insights      []Insight        `json:"insights"`
	OriginFetches map[string]Fetch `json:"originFetches"`
	PackageFile   string           `json:"packageFile"`
	RequestCount  int              `json:"requestCount"`
	StartTime     int              `json:"startTime"`
}

type Fetch struct {
	ElapsedTime int       `json:"elapsedTime"`
	FetchId     string    `json:"fetchID"`
	Insights    []Insight `json:"insights"`
	OriginIdx   int       `json:"originIdx"`
	RemoteAddr  string    `json:"remoteAddr"`
	RemoteHost  string    `json:"remoteHost"`
	RemotePort  int       `json:"remotePort"`
	Req         string    `json:"req"`
	ReqID       string    `json:"reqID"`
	Resp        string    `json:"resp"`
	TraceID     string    `json:"traceID"`
}

type Insight struct {
	Data  map[string]any `json:"data"`
	Key   string         `json:"key"`
	Level string         `json:"level"`
	Pos   string         `json:"pos"`
}

type Event struct {
	Attributes     Attributes `json:"attribs"`
	FnName         string     `json:"fnName"`
	IsAsync        bool       `json:"isAsync"`
	IsComplete     bool       `json:"isComplete"`
	Logs           []Log      `json:"logs"`
	PrevAttribs    Attributes `json:"prevAttribs"`
	RanBoilerplate bool       `json:"ranBoilerplate"`
	ReqID          string     `json:"reqID"`
	SeqIdx         int        `json:"seqIdx"`
	Server         Server     `json:"server"`
	Time           int        `json:"time"`
	TraceID        string     `json:"traceID"`
	Type           string     `json:"type"`
}

type Log struct {
	Attribs Attributes `json:"attribs"`
	Content string     `json:"content"`
	FnName  string     `json:"fnName"`
	ReqID   string     `json:"reqID"`
	SeqIdx  int        `json:"seqIdx"`
	Server  Server     `json:"server"`
	Time    int        `json:"time"`
	TraceID string     `json:"traceID"`
	Type    string     `json:"type"`
}

type Server struct {
	NodeId string `json:"nodeID"`
	Pop    string `json:"pop"`
}

type Attributes struct {
	AlwaysMiss bool   `json:"alwaysMiss"`
	Backend    string `json:"backend"`
	IgnoreBusy bool   `json:"ignoreBusy"`
	IsESI      bool   `json:"isESI"`
	IsH2       bool   `json:"isH2"`
	IsPush     bool   `json:"isPush"`
	Method     string `json:"method"`
	Restarts   int    `json:"restarts"`
	Return     string `json:"return"`
	URL        string `json:"url"`
}
