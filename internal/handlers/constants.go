package handlers

import (
	"github.com/Ayano2000/push/internal/pkg/types"
)

// Context values
const (
	muxContextKey      types.MuxContextKey      = "router"
	urlParamContextKey types.UrlParamContextKey = "parameters"
)

// Error messages
const createWebhookErrorMessage = "Failed to create webhook"
const getWebhookContentErrorMessage = "Failed to fetch webhook content"
const getWebhooksErrorMessage = "Failed to fetch webhooks"
const invalidJQFilterErrorMessage = "JQ filter is invalid"
const jqTransformErrorMessage = "Failed to process JQ filter on request body"
const minioUploadErrorMessage = "Failed to upload request body to minio"
const requestBodyDecodingErrorMessage = "Failed to decode the request body"
