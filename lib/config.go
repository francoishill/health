package health

import (
	"encoding/json"
	"flag"
	"os"
	"time"
)

const (
	// DefaultMaxStats is the default number of deltas to keep per host.
	DefaultMaxStats = 128
	// DefaultPingTimeout is the connection timeout.
	DefaultPingTimeout = Duration(1000 * time.Millisecond)
	//DefaultPollInterval is the default time between pings.
	DefaultPollInterval = Duration(2000 * time.Millisecond)
)

// NewConfigFromFlags returns a new config object by parsing flags.
func NewConfigFromFlags() (*Config, error) {
	c := &Config{
		PingTimeout:  DefaultPingTimeout,
		MaxStats:     DefaultMaxStats,
		PollInterval: DefaultPollInterval,
	}
	return c, c.ParseFlags()
}

// Config is the healthcheck configuration.
type Config struct {
	PingTimeout      Duration `json:"ping_timeout"`
	MaxStats         int      `json:"max_stats"`
	PollInterval     Duration `json:"interval"`
	Hosts            []string `json:"hosts"`
	ShowNotification bool     `json:"show_notification"`
	Verbose          bool     `json:"verbose"`
}

// LoadFromPath loads a config from a path.
func (c *Config) LoadFromPath(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	return decoder.Decode(c)
}

// ParseFlags parses commandline flags into a config object.
func (c *Config) ParseFlags() error {
	var hosts HostsFlag
	flag.Var(&hosts, "host", "Host(s) to ping.")

	pollInterval := flag.Int("interval", 30000, "Server polling interval in milliseconds")
	showNotifications := flag.Bool("notification", true, "Show OS X Notification on `down`")
	configFilePath := flag.String("config", "", "Load configuration from a file.")

	//parse the arguments
	flag.Parse()

	if configFilePath != nil && len(*configFilePath) != 0 {
		return c.LoadFromPath(*configFilePath)
	}
	if pollInterval != nil {
		c.PollInterval = Duration(time.Duration(*pollInterval) * time.Millisecond)
	}
	if len(hosts) != 0 {
		c.Hosts = append(c.Hosts, hosts...)
	}

	if showNotifications != nil {
		c.ShowNotification = *showNotifications
	}

	return nil
}

// HostNameLength returns the length of the longest host name in the config.
func (c *Config) HostNameLength() int {
	longestHostName := 0
	for x := 0; x < len(c.Hosts); x++ {
		l := len(c.Hosts[x])
		if l > longestHostName {
			longestHostName = l
		}
	}
	return longestHostName
}
