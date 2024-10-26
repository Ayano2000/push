package types

type Bucket struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	JQFilter    string `json:"jq_filter,omitempty"`
	ForwardTo   string `json:"forward_to,omitempty"`
}
