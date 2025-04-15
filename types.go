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

type Lints map[string]LintStatus

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

type ExecuteResponse struct {
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

const (
	FiddleTypeVCL        = "vcl"
	FiddleTypeJavaScript = "javascript"
	FiddleTypeGo         = "go"
	FittleTypeRust       = "rust"
)

type Request struct {
	Method          string         `json:"method"`
	Path            string         `json:"path"`
	Headers         string         `json:"headers"`
	Body            string         `json:"body"`
	Data            map[string]any `json:"data"` // Uncertain of function
	EnableCluster   bool           `json:"enableCluster"`
	EnableShield    bool           `json:"enableShield"`
	UseFreshCache   bool           `json:"useFreshCache"`
	ConnType        string         `json:"connType"`
	SourceIP        string         `json:"sourceIP"`
	FollowRedirects bool           `json:"followRedirects"`
	Tests           string         `json:"tests"`
	Delay           int            `json:"delay"`
}

func NewRequest() Request {
	return Request{
		Path:          "/",
		Method:        MethodGet,
		EnableCluster: true,
		ConnType:      ConnectionTypeHTTP2,
		SourceIP:      SourceIpClient,
	}
}

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodOptions = "OPTIONS"
	MethodPurge   = "PURGE"
	MethodInvalid = "INVALID"

	SourceIpClient        = "client"
	SourceIpServer        = "server"
	SourceIpChina         = "cn"
	SourceIpGermany       = "de"
	SourceIpJapan         = "jp"
	SourceIpRussia        = "ru"
	SourceIpBrazil        = "br"
	SourceIpSouthAfrica   = "za"
	SourceIpUnitedKingdon = "uk"
	SourceIpUnitedStates  = "us"

	ConnectionTypeHTTP2      = "h2"
	ConnectionTypeHTTP1TLS   = "h1"
	ConnectionTypeHTTP1Plain = "http"
)

type Source struct {
	VCLSource
	ComputeSource
}

type VCLSource struct {
	Init    string `json:"init,omitempty"`
	Recv    string `json:"recv,omitempty"`
	Fetch   string `json:"fetch,omitempty"`
	Deliver string `json:"deliver,omitempty"`
	Hash    string `json:"hash,omitempty"`
	Hit     string `json:"hit,omitempty"`
	Miss    string `json:"miss,omitempty"`
	Pass    string `json:"pass,omitempty"`
	Log     string `json:"log,omitempty"`
}

type ComputeSource struct {
	Deps     string `json:"deps,omitempty"`
	Main     string `json:"main,omitempty"`
	Manifest string `json:"manifest,omitempty"`
}

type StreamEvent struct {
	Type string `json:"-"`
	UpdateResultEvent
	WaitForSyncEvent
}

type UpdateResultEvent struct {
	ClientFetches map[string]ClientFetch `json:"clientFetches"`
	Events        []Event                `json:"events"`
	ExecHost      string                 `json:"execHost"`
	ExecVersion   int                    `json:"execVersion"`
	ID            string                 `json:"id"`
	Insights      []Insight              `json:"insights"`
	OriginFetches map[string]OriginFetch `json:"originFetches"`
	PackageFile   string                 `json:"packageFile"`
	RequestCount  int                    `json:"requestCount"`
	StartTime     int                    `json:"startTime"`
}

type ClientFetch struct {
	Delay             int          `json:"delay"`
	IsRedirected      bool         `json:"isRedirected"`
	Req               string       `json:"req"`
	Complete          bool         `json:"complete"`
	Resp              string       `json:"resp"`
	RespType          string       `json:"respType"`
	IsText            bool         `json:"isText"`
	IsImage           bool         `json:"isImage"`
	PreviewType       string       `json:"previewType"`
	BodyPreview       string       `json:"bodyPreview"`
	BodyBytesReceived int          `json:"bodyBytesReceived"`
	BodyChunkCount    int          `json:"bodyChunkCount"`
	Tests             []TestResult `json:"tests"`
	Insights          []Insight    `json:"insights"`
}

type TestResult struct {
	Label      string `json:"label"`
	TestExpr   string `json:"testExpr"`
	AsyncDelay int    `json:"asyncDelay"`
	Pass       bool   `json:"pass"`
	Expected   any    `json:"expected"`
	Actual     any    `json:"actual"`
}

type OriginFetch struct {
	ElapsedTime int       `json:"elapsedTime"`
	FetchID     string    `json:"fetchID"`
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
	ReqID          string     `json:"reqID"`
	TraceID        string     `json:"traceID"`
	Time           int        `json:"time"`
	SeqIdx         int        `json:"seqIdx"`
	Attribs        Attributes `json:"attribs"`
	Type           string     `json:"type"`
	FnName         string     `json:"fnName,omitempty"`
	IsAsync        *bool      `json:"isAsync,omitempty"`
	IsComplete     *bool      `json:"isComplete,omitempty"`
	Logs           []Log      `json:"logs,omitempty"`
	PrevAttribs    Attributes `json:"prevAttribs,omitempty"`
	RanBoilerplate *bool      `json:"ranBoilerplate,omitempty"`
	Server         Server     `json:"server,omitempty"`
	FetchID        string     `json:"fetchID,omitempty"`
	URL            string     `json:"url,omitempty"`
	Method         string     `json:"method,omitempty"`
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
	AlwaysMiss  bool   `json:"alwaysMiss,omitempty"`
	Backend     string `json:"backend,omitempty"`
	IgnoreBusy  bool   `json:"ignoreBusy,omitempty"`
	IsESI       bool   `json:"isESI,omitempty"`
	IsH2        bool   `json:"isH2,omitempty"`
	IsPush      bool   `json:"isPush,omitempty"`
	Method      string `json:"method,omitempty"`
	Restarts    int    `json:"restarts,omitempty"`
	Return      string `json:"return,omitempty"`
	URL         string `json:"url,omitempty"`
	Hits        int    `json:"hits,omitempty"`
	Digest      string `json:"digest,omitempty"`
	Age         int    `json:"age,omitempty"`
	SKeys       string `json:"skeys,omitempty"`
	TTL         int    `json:"ttl,omitempty"`
	IsPCI       bool   `json:"isPCI,omitempty"`
	LastUse     int    `json:"lastuse,omitempty"`
	StaleExists bool   `json:"staleExists,omitempty"`
	Status      int    `json:"status,omitempty"`
	State       string `json:"state,omitempty"`
	StatusText  string `json:"statusText,omitempty"`
}

type Purge struct {
	Key string `json:"key"`
}

type WaitForSyncEvent struct {
	ServiceID        string `json:"serviceID"`
	ServiceName      string `json:"serviceName"`
	SavedVersion     int    `json:"savedVersion"`
	PublishedVersion int    `json:"publishedVersion"`
	ExecutedVersion  int    `json:"executedVersion"`
	PublishedFiddle  string `json:"publishedFiddle"`
	ExecutedFiddle   string `json:"executedFiddle"`
	ReqHost          string `json:"reqHost"`
	Status           int    `json:"status"`
	RetryCount       int    `json:"retryCount"`
}
