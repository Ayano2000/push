package types

type Bucket struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Path            string `json:"path"`
	Method          string `json:"method"`
	JQFilter        string `json:"jq_filter"`
	ForwardTo       string `json:"forward_to"`
	PreservePayload bool   `json:"preserve_payload"`
}
