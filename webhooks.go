package paypal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/yinyajiang/paypal/v4/json"
)

// CreateWebhook - Subscribes your webhook listener to events.
// Endpoint: POST /v1/notifications/webhooks
func (c *Client) CreateWebhook(ctx context.Context, createWebhookRequest *CreateWebhookRequest) (*Webhook, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/webhooks"), createWebhookRequest)
	webhook := &Webhook{}
	if err != nil {
		return webhook, err
	}

	err = c.SendWithAuth(req, webhook)
	return webhook, err
}

// GetWebhook - Shows details for a webhook, by ID.
// Endpoint: GET /v1/notifications/webhooks/ID
func (c *Client) GetWebhook(ctx context.Context, webhookID string) (*Webhook, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/notifications/webhooks/", webhookID), nil)
	webhook := &Webhook{}
	if err != nil {
		return webhook, err
	}

	err = c.SendWithAuth(req, webhook)
	return webhook, err
}

// UpdateWebhook - Updates a webhook to replace webhook fields with new values.
// Endpoint: PATCH /v1/notifications/webhooks/ID
func (c *Client) UpdateWebhook(ctx context.Context, webhookID string, fields []WebhookField) (*Webhook, error) {
	req, err := c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("%s/v1/notifications/webhooks/%s", c.APIBase, webhookID), fields)
	webhook := &Webhook{}
	if err != nil {
		return webhook, err
	}

	err = c.SendWithAuth(req, webhook)
	return webhook, err
}

// ListWebhooks - Lists webhooks for an app.
// Endpoint: GET /v1/notifications/webhooks
func (c *Client) ListWebhooks(ctx context.Context, anchorType string) (*ListWebhookResponse, error) {
	if len(anchorType) == 0 {
		anchorType = AncorTypeApplication
	}
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/webhooks"), nil)
	q := req.URL.Query()
	q.Add("anchor_type", anchorType)
	req.URL.RawQuery = q.Encode()
	resp := &ListWebhookResponse{}
	if err != nil {
		return nil, err
	}

	err = c.SendWithAuth(req, resp)
	return resp, err
}

// DeleteWebhook - Deletes a webhook, by ID.
// Endpoint: DELETE /v1/notifications/webhooks/ID
func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) error {
	req, err := c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/v1/notifications/webhooks/%s", c.APIBase, webhookID), nil)
	if err != nil {
		return err
	}

	err = c.SendWithAuth(req, nil)
	return err
}

// VerifyWebhookSignature - Use this to verify the signature of a webhook recieved from paypal.
// Endpoint: POST /v1/notifications/verify-webhook-signature
func (c *Client) VerifyWebhookSignature(ctx context.Context, httpReq *http.Request, webhookID string) (*VerifyWebhookResponse, error) {
	// Read the content
	var bodyBytes []byte
	if httpReq.Body != nil {
		bodyBytes, _ = io.ReadAll(httpReq.Body)
	} else {
		return nil, errors.New("Cannot verify webhook for HTTP Request with empty body.")
	}
	// Restore the io.ReadCloser to its original state
	httpReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return c.VerifyWebhookSignature2(ctx, bodyBytes, httpReq.Header, webhookID)
}

// VerifyWebhookSignature - Use this to verify the signature of a webhook recieved from paypal.
// Endpoint: POST /v1/notifications/verify-webhook-signature
func (c *Client) VerifyWebhookSignature2(ctx context.Context, bodyBytes []byte, header http.Header, webhookID string) (*VerifyWebhookResponse, error) {
	type verifyWebhookSignatureRequest struct {
		AuthAlgo         string          `json:"auth_algo,omitempty"`
		CertURL          string          `json:"cert_url,omitempty"`
		TransmissionID   string          `json:"transmission_id,omitempty"`
		TransmissionSig  string          `json:"transmission_sig,omitempty"`
		TransmissionTime string          `json:"transmission_time,omitempty"`
		WebhookID        string          `json:"webhook_id,omitempty"`
		Event            json.RawMessage `json:"webhook_event,omitempty"`
	}

	// Read the content
	var bodyBytes []byte
	if httpReq.Body != nil {
		bodyBytes, _ = io.ReadAll(httpReq.Body)
	} else {
		return nil, errors.New("cannot verify webhook for HTTP Request with empty body")
	}
	// Restore the io.ReadCloser to its original state
	httpReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	verifyRequest := verifyWebhookSignatureRequest{
		AuthAlgo:         header.Get("PAYPAL-AUTH-ALGO"),
		CertURL:          header.Get("PAYPAL-CERT-URL"),
		TransmissionID:   header.Get("PAYPAL-TRANSMISSION-ID"),
		TransmissionSig:  header.Get("PAYPAL-TRANSMISSION-SIG"),
		TransmissionTime: header.Get("PAYPAL-TRANSMISSION-TIME"),
		WebhookID:        webhookID,
		Event:            json.RawMessage(bodyBytes),
	}

	response := &VerifyWebhookResponse{}

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/verify-webhook-signature"), verifyRequest)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetWebhookEventTypes - Lists all webhook event types.
// Endpoint: GET /v1/notifications/webhooks-event-types
func (c *Client) GetWebhookEventTypes(ctx context.Context) (*WebhookEventTypesResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/webhooks-event-types"), nil)
	q := req.URL.Query()

	req.URL.RawQuery = q.Encode()
	resp := &WebhookEventTypesResponse{}
	if err != nil {
		return nil, err
	}

	err = c.SendWithAuth(req, resp)
	return resp, err
}
