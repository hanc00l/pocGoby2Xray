package gobypoc

type Check struct {
	Type      string `json:"type"`
	Variable  string `json:"variable"`
	Operation string `json:"operation"`
	Value     string `json:"value"`
	Bz        string `json:"bz"`
}
type ResponseTest struct {
	Type      string  `json:"type"`
	Operation string  `json:"operation"`
	Checks    []Check `json:"checks"`
}
type Request struct {
	Method         string            `json:"method"`
	Uri            string            `json:"uri"`
	FollowRedirect bool              `json:"follow_redirect"`
	Header         map[string]string `json:"header"`
	DataType       string            `json:"data_type"`
	Data           string            `json:"data"`
	SetVariable    []string          `json:"set_variable"`
}
type ScanStep struct {
	Requests     Request
	ResponseTest ResponseTest
	SetVariable  []string
}

type PocJson struct {
	Name           string   `json:"Name"`
	Level          string   `json:"Level"`
	Tags           []string `json:"Tags"`
	Query          string   `json:"GobyQuery"`
	Description    string   `json:"Description"`
	Product        string   `json:"Product"`
	Homepage       string   `json:"Homepage"`
	Author         string   `json:"Author"`
	Impact         string   `json:"Impact"`
	Recommandation string   `json:"Recommandation"`
	References     []string `json:"References"`
	HasExp         bool     `json:"HasExp"`
	ExpParams      []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
		Show  string `json:"show"`
	} `json:"ExpParams"`
	ScanSteps    []interface{} `json:"ScanSteps"`
	ExploitSteps []interface{} `json:"ExploitSteps"`
	PostTime     string        `json:"PostTime"`
	Version      string        `json:"GobyVersion"`
}
