package xraypocv2

import (
	"gopkg.in/yaml.v2"
)

type Rule struct {
	Request    RuleRequest   `yaml:"request"`
	Expression string        `yaml:"expression"`
	Output     yaml.MapSlice `yaml:"output,omitempty"`
}

// 用于帮助yaml解析，保证Rule有序
type RuleMapItem struct {
	Key   string
	Value Rule
}

type VariableMapItem struct {
	Key   string
	Value string
}

// 用于帮助yaml解析，保证Rule有序
//type RuleMapSlice []RuleMapItem

type RuleRequest struct {
	Cache           bool              `yaml:"cache,omitempty"`
	Method          string            `yaml:"method"`
	Path            string            `yaml:"path"`
	Headers         map[string]string `yaml:"headers,omitempty"`
	Body            string            `yaml:"body,omitempty"`
	FollowRedirects bool              `yaml:"follow_redirects,omitempty"`
	Content         string            `yaml:"content,omitempty"`
	ReadTimeout     string            `yaml:"read_timeout,omitempty"`
	ConnectionID    string            `yaml:"connection_id,omitempty"`
}

type Infos struct {
	ID         string `yaml:"id,omitempty"`
	Name       string `yaml:"name,omitempty"`
	Version    string `yaml:"version,omitempty"`
	Type       string `yaml:"type,omitempty"`
	Confidence int    `yaml:"confidence,omitempty"`
}

type HostInfo struct {
	Hostname string `yaml:"hostname"`
}

type Vulnerability struct {
	ID    string `yaml:"id"`
	Match string `yaml:"match"`
}

type FingerPrint struct {
	Infos    []Infos  `yaml:"infos"`
	HostInfo HostInfo `yaml:"host_info"`
}
type Detail struct {
	Author        string        `yaml:"author,omitempty"`
	Links         []string      `yaml:"links,omitempty"`
	FingerPrint   FingerPrint   `yaml:"fingerprint,omitempty"`
	Vulnerability Vulnerability `yaml:"vulnerability,omitempty"`
	Description   string        `yaml:"description,omitempty"`
	Version       string        `yaml:"version,omitempty"`
	Tags          string        `yaml:"tags,omitempty"`
}

type Payloads struct {
	Continue bool          `yaml:"continue,omitempty"`
	Payloads yaml.MapSlice `yaml:"payloads,omitempty"`
}

type Poc struct {
	Name       string        `yaml:"name"`
	Manual     bool          `yaml:"manual"`
	Query      string        `yaml:"query,omitempty"`
	Transport  string        `yaml:"transport"`
	Set        yaml.MapSlice `yaml:"set,omitempty"`
	Payloads   Payloads      `yaml:"payloads,omitempty"`
	Rules      yaml.MapSlice `yaml:"rules"`
	Expression string        `yaml:"expression"`
	Detail     Detail        `yaml:"detail,omitempty"`
}
