package types

type Webhook struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Path            string `json:"path"`
	Method          string `json:"method"`
	JQFilter        string `json:"jq_filter"`
	ForwardTo       string `json:"forward_to"`
	PreservePayload bool   `json:"preserve_payload"`
}

// WebhookRegistrar defines methods for registering webhooks.
type WebhookRegistrar interface {
	RegisterWebhook(webhook Webhook)
}
