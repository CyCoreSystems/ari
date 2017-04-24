package ari

// DialplanCEP describes a location in the dialplan (context,extension,priority)
type DialplanCEP struct {
	Context  string `json:"context"`
	Exten    string `json:"exten"`
	Priority int64  `json:"priority"`
}
