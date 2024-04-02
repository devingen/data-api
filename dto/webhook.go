package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebhookPreRequest struct {
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	QueryParams map[string][]string    `json:"queryParams"`
	Headers     map[string]string      `json:"headers"`
	Body        map[string]interface{} `json:"body"`
}

type WebhookPreResponse struct {
	// QueryEnhance is used to add extra fields to the query for GET list requests.
	QueryEnhance *QueryEnhance `json:"queryEnhance"`

	// EmailTemplateID is used for webhook api to decide the email template from the email template identifier.
	// It gives the flexibility to allow creating email templates with same identifier
	EmailTemplateID *string `json:"emailTemplateId"`
}

type QueryEnhance struct {
	// IDsIn filters the returned items to contain only the given IDs. All items are returned otherwise.
	IDsIn []primitive.ObjectID `json:"idsIn"`
}

type WebhookFinalRequest struct {
	Method         string              `json:"method"`
	Path           string              `json:"path"`
	PathParameters map[string]string   `json:"pathParameters"`
	QueryParams    map[string][]string `json:"queryParams"`
	Headers        map[string]string   `json:"headers"`
	ResponseBody   interface{}         `json:"responseBody"`
}

type WebhookError struct {
	Status   int      `json:"status"`
	Messages []string `json:"messages"`
}
