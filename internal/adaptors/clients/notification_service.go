package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"payment-engine/internal/domain"
	"payment-engine/internal/domain/application"

	"github.com/recodextech/container"
)

const notificationSendEndpoint = "/notifications"

// SendNotificationRequest is a push notification to every registered device of a user.
type SendNotificationRequest struct {
	UserID      string            `json:"user_id"`
	Type        string            `json:"type"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Data        map[string]string `json:"data,omitempty"`
}

// NotificationServiceClient sends push notifications through notification-service (Firebase).
type NotificationServiceClient struct {
	httpClient *http.Client
	conf       *Config
}

func (n *NotificationServiceClient) Init(di container.Container) error {
	n.conf = di.GetGlobalConfig(application.ModuleHTTPClient).(*Config)
	n.httpClient = &http.Client{Timeout: time.Duration(n.conf.HTTPTimeoutInSeconds) * time.Second}
	return nil
}

// SendToUser pushes the notification to the recipient's devices. actorUserID and actorAccountID are
// sent as the calling user headers; for system triggered pushes they are the recipient's own ids.
func (n *NotificationServiceClient) SendToUser(ctx context.Context, actorUserID, actorAccountID string, req SendNotificationRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal notification request: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, n.conf.NotificationServiceBaseURL+notificationSendEndpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create notification request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set(domain.ContextKeyUserID.String(), actorUserID)
	httpReq.Header.Set(domain.ContextKeyAccountID.String(), actorAccountID)

	res, err := n.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("notification service request failed: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		raw, _ := io.ReadAll(res.Body)
		return fmt.Errorf("notification service returned status %d: %s", res.StatusCode, strings.TrimSpace(string(raw)))
	}
	return nil
}
