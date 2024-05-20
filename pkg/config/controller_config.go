package config

import "github.com/spf13/pflag"

const (
	flagLogLevel                   = "log-level"
	flagLogFile                    = "log-file"
	flagMaxConcurrentReconciles    = "max-concurrent-reconciles"
	defaultLogLevel                = "info"
	defaultLogFile                 = "/var/log/aws-routed-eni/network-policy-agent.log"
	defaultMaxConcurrentReconciles = 3
	flagEnableIPv6                 = "enable-ipv6"
)

// ControllerConfig contains the controller configuration
type ControllerConfig struct {
	// Log level for the controller logs
	LogLevel string
	// Local log file for Network Policy Agent
	LogFile string
	// MaxConcurrentReconciles specifies the max number of reconcile loops
	MaxConcurrentReconciles int
	// Enable IPv6 mode
	EnableIPv6 bool
	// Configurations for the Controller Runtime
	RuntimeConfig RuntimeConfig
}

func (cfg *ControllerConfig) BindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cfg.LogLevel, flagLogLevel, defaultLogLevel,
		"Set the controller log level - info, debug")
	fs.StringVar(&cfg.LogFile, flagLogFile, defaultLogFile, ""+
		"Set the controller log file - if not specified logs are written to stdout")
	fs.IntVar(&cfg.MaxConcurrentReconciles, flagMaxConcurrentReconciles, defaultMaxConcurrentReconciles, ""+
		"Maximum number of concurrent reconcile loops")
	fs.BoolVar(&cfg.EnableIPv6, flagEnableIPv6, false, "If enabled, Network Policy agent will operate in IPv6 mode")

	cfg.RuntimeConfig.BindFlags(fs)
}
