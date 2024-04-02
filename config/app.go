package config

// ApiConfig defines the environment variable configuration for the whole app
type ApiConfig struct {
	// Port is the port of the HTTP server.
	Port string `envconfig:"port" default:"1003"`

	// LogLevel defines the log level.
	LogLevel string `envconfig:"log_level" default:"info"`

	// Mongo is the configuration of the MongoDB server.
	Mongo Mongo `envconfig:"mongo"`

	// Webhook is the configuration of the webhook server.
	Webhook Webhook `envconfig:"webhook"`

	// EnableCreate used to enable/disable create function
	EnableCreate bool `envconfig:"enable_create" default:"false"`

	// EnableUpdate used to enable/disable update function
	EnableUpdate bool `envconfig:"enable_update" default:"false"`

	// EnableDelete used to enable/disable delete function
	EnableDelete bool `envconfig:"enable_delete" default:"false"`
}

// Mongo defines the environment variable configuration for MongoDB
type Mongo struct {
	// URI is the MongoDB server URI.
	URI string `envconfig:"uri" default:"mongodb://localhost:27017"`
}

// Webhook defines the environment variable configuration for webhook
type Webhook struct {

	// URL is called before/after every request.
	URL string `envconfig:"url" default:""`

	// URL is called before/after every request.
	Headers string `envconfig:"headers" default:""`
}
