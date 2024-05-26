package types

import "time"

// Application represents an application being monitored.
type Application struct {
	ApplicationID int       // Unique identifier for the application
	Name          string    // Name of the application
	Description   string    // Description of the application
	CreatedAt     time.Time // Time when the application was created
	UpdatedAt     time.Time // Time when the application was last updated
}

// Endpoint represents a specific URL of an application to be monitored.
type Endpoint struct {
	EndpointID         int       // Unique identifier for the endpoint
	ApplicationID      int       // Identifier for the application this endpoint belongs to
	URL                string    // URL to monitor
	MonitoringInterval int       // Frequency in seconds at which the endpoint should be checked
	IsActive           bool      // Whether monitoring is active for this endpoint
	CreatedAt          time.Time // Time when the endpoint was created
	UpdatedAt          time.Time // Time when the endpoint was last updated
}

// UptimeLog represents a log entry for monitoring an endpoint.
type UptimeLog struct {
	EndpointID   int       // Identifier for the endpoint this log belongs to
	StatusCode   int       // HTTP status code received from the endpoint
	ResponseTime int       // Time taken to get a response from the endpoint
	IsUp         bool      // Whether the endpoint was up or down
	Timestamp    time.Time // Time when the check was performed
}

// NotificationChannel represents a channel through which notifications are sent.
type NotificationChannel struct {
	ChannelID     int       // Unique identifier for the channel
	ApplicationID int       // Identifier for the application this channel belongs to
	Type          string    // Type of the notification channel (e.g., Slack, Email)
	Details       string    // Details specific to the notification type, possibly in JSON or another structured format
	IsActive      bool      // Whether this channel is active for sending notifications
	CreatedAt     time.Time // Time when the channel was created
	UpdatedAt     time.Time // Time when the channel was last updated
}

// Alert represents an alert that has been sent out.
type Alert struct {
	AlertID    int       // Unique identifier for the alert
	ChannelID  int       // Identifier for the notification channel used
	EndpointID int       // Identifier for the endpoint related to this alert
	Message    string    // Alert message that was sent
	SentAt     time.Time // Time when the alert was sent
}
