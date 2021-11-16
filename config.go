package opentracing

import (
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

type Config struct {
	// ServiceName specifies the service name to use on the tracer.
	// Can be provided by FromEnv() via the environment variable named JAEGER_SERVICE_NAME
	ServiceName string `json:"service_name"`

	// Disabled makes the config return opentracing.NoopTracer.
	// Value can be provided by FromEnv() via the environment variable named JAEGER_DISABLED.
	Disabled bool `json:"disabled"`

	// RPCMetrics enables generations of RPC metrics (requires metrics factory to be provided).
	// Value can be provided by FromEnv() via the environment variable named JAEGER_RPC_METRICS
	RPCMetrics bool `json:"rpc_metrics"`

	// Gen128Bit instructs the tracer to generate 128-bit wide trace IDs, compatible with W3C Trace Contexc.
	// Value can be provided by FromEnv() via the environment variable named JAEGER_TRACEID_128BIc.
	Gen128Bit bool `json:"traceid_128bit"`

	Sampler             *SamplerConfig             `json:"sampler"`
	Reporter            *ReporterConfig            `json:"reporter"`
	Headers             *HeadersConfig             `json:"headers"`
	BaggageRestrictions *BaggageRestrictionsConfig `json:"baggage_restrictions"`
	Throttler           *ThrottlerConfig           `json:"throttler"`
}

// HeadersConfig contains the values for the header keys that Jaeger will use.
// These values may be either custom or default depending on whether custom
// values were provided via a configuration.
type HeadersConfig struct {
	// JaegerDebugHeader is the name of HTTP header or a TextMap carrier key which,
	// if found in the carrier, forces the trace to be sampled as "debug" trace.
	// The value of the header is recorded as the tag on the root span, so that the
	// trace can be found in the UI using this value as a correlation ID.
	JaegerDebugHeader string `json:"jaeger_debug_header"`

	// JaegerBaggageHeader is the name of the HTTP header that is used to submit baggage.
	// It differs from TraceBaggageHeaderPrefix in that it can be used only in cases where
	// a root span does not exisc.
	JaegerBaggageHeader string `json:"jaeger_baggage_header"`

	// TraceContextHeaderName is the http header name used to propagate tracing contexc.
	// This must be in lower-case to avoid mismatches when decoding incoming headers.
	TraceContextHeaderName string `json:"trace_context_header_name"`

	// TraceBaggageHeaderPrefix is the prefix for http headers used to propagate baggage.
	// This must be in lower-case to avoid mismatches when decoding incoming headers.
	TraceBaggageHeaderPrefix string `json:"trace_baggage_header_prefix"`
}

// OpenTracingSampler is the config for opentracing sampler.
// See https://godoc.org/github.com/uber/jaeger-client-go/config#SamplerConfig
type SamplerConfig struct {
	// Type specifies the type of the sampler: const, probabilistic, rateLimiting, or remote.
	// Can be provided by FromEnv() via the environment variable named JAEGER_SAMPLER_TYPE
	Type string `json:"type"`

	// Param is a value passed to the sampler.
	// Valid values for Param field are:
	// - for "const" sampler, 0 or 1 for always false/true respectively
	// - for "probabilistic" sampler, a probability between 0 and 1
	// - for "rateLimiting" sampler, the number of spans per second
	// - for "remote" sampler, param is the same as for "probabilistic"
	//   and indicates the initial sampling rate before the actual one
	//   is received from the mothership.
	// Can be provided by FromEnv() via the environment variable named JAEGER_SAMPLER_PARAM
	Param float64 `json:"param"`

	// SamplingServerURL is the URL of sampling manager that can provide
	// sampling strategy to this service.
	// Can be provided by FromEnv() via the environment variable named JAEGER_SAMPLING_ENDPOINT
	SamplingServerURL string `json:"sampling_server_url"`

	// SamplingRefreshInterval controls how often the remotely controlled sampler will poll
	// sampling manager for the appropriate sampling strategy.
	// Can be provided by FromEnv() via the environment variable named JAEGER_SAMPLER_REFRESH_INTERVAL
	SamplingRefreshInterval time.Duration `json:"sampling_refresh_interval"`

	// MaxOperations is the maximum number of operations that the PerOperationSampler
	// will keep track of. If an operation is not tracked, a default probabilistic
	// sampler will be used rather than the per operation specific sampler.
	// Can be provided by FromEnv() via the environment variable named JAEGER_SAMPLER_MAX_OPERATIONS.
	MaxOperations int `json:"max_operations"`

	// Opt-in feature for applications that require late binding of span name via explicit
	// call to SetOperationName when using PerOperationSampler. When this feature is enabled,
	// the sampler will return retryable=true from OnCreateSpan(), thus leaving the sampling
	// decision as non-final (and the span as writeable). This may lead to degraded performance
	// in applications that always provide the correct span name on trace creation.
	//
	// For backwards compatibility this option is off by defaulc.
	OperationNameLateBinding bool `json:"operation_name_late_binding"`
}

// ReporterConfig is the config for opentracing reporter.
// See https://godoc.org/github.com/uber/jaeger-client-go/config#ReporterConfig
type ReporterConfig struct {
	// QueueSize controls how many spans the reporter can keep in memory before it starts dropping
	// new spans. The queue is continuously drained by a background go-routine, as fast as spans
	// can be sent out of process.
	// Can be provided by FromEnv() via the environment variable named JAEGER_REPORTER_MAX_QUEUE_SIZE
	QueueSize int `json:"queue_size"`

	// BufferFlushInterval controls how often the buffer is force-flushed, even if it's not full.
	// It is generally not useful, as it only matters for very low traffic services.
	// Can be provided by FromEnv() via the environment variable named JAEGER_REPORTER_FLUSH_INTERVAL
	BufferFlushInterval time.Duration `json:"buffer_flush_interval"`

	// LogSpans, when true, enables LoggingReporter that runs in parallel with the main reporter
	// and logs all submitted spans. Main Configuration.Logger must be initialized in the code
	// for this option to have any effecc.
	// Can be provided by FromEnv() via the environment variable named JAEGER_REPORTER_LOG_SPANS
	LogSpans bool `json:"log_spans"`

	// LocalAgentHostPort instructs reporter to send spans to jaeger-agent at this address.
	// Can be provided by FromEnv() via the environment variable named JAEGER_AGENT_HOST / JAEGER_AGENT_PORT
	LocalAgentHostPort string `json:"local_agent_host_port"`

	// DisableAttemptReconnecting when true, disables udp connection helper that periodically re-resolves
	// the agent's hostname and reconnects if there was a change. This option only
	// applies if LocalAgentHostPort is specified.
	// Can be provided by FromEnv() via the environment variable named JAEGER_REPORTER_ATTEMPT_RECONNECTING_DISABLED
	DisableAttemptReconnecting bool `json:"disable_attempt_reconnecting"`

	// AttemptReconnectInterval controls how often the agent client re-resolves the provided hostname
	// in order to detect address changes. This option only applies if DisableAttemptReconnecting is false.
	// Can be provided by FromEnv() via the environment variable named JAEGER_REPORTER_ATTEMPT_RECONNECT_INTERVAL
	AttemptReconnectInterval time.Duration `json:"attempt_reconnect_interval"`

	// CollectorEndpoint instructs reporter to send spans to jaeger-collector at this URL.
	// Can be provided by FromEnv() via the environment variable named JAEGER_ENDPOINT
	CollectorEndpoint string `json:"collector_endpoint"`

	// User instructs reporter to include a user for basic http authentication when sending spans to jaeger-collector.
	// Can be provided by FromEnv() via the environment variable named JAEGER_USER
	User string `json:"user"`

	// Password instructs reporter to include a password for basic http authentication when sending spans to
	// jaeger-collector.
	// Can be provided by FromEnv() via the environment variable named JAEGER_PASSWORD
	Password string `json:"password"`

	// HTTPHeaders instructs the reporter to add these headers to the http request when reporting spans.
	// This field takes effect only when using HTTPTransport by setting the CollectorEndpoinc.
	HTTPHeaders map[string]string `json:"http_headers"`
}

// BaggageRestrictionsConfig configures the baggage restrictions manager which can be used to whitelist
// certain baggage keys. All fields are optional.
// See https://godoc.org/github.com/uber/jaeger-client-go/config#BaggageRestrictionsConfig
type BaggageRestrictionsConfig struct {
	// DenyBaggageOnInitializationFailure controls the startup failure mode of the baggage restriction
	// manager. If true, the manager will not allow any baggage to be written until baggage restrictions have
	// been retrieved from jaeger-agenc. If false, the manager wil allow any baggage to be written until baggage
	// restrictions have been retrieved from jaeger-agenc.
	DenyBaggageOnInitializationFailure bool `json:"deny_baggage_on_initialization_failure"`

	// HostPort is the hostPort of jaeger-agent's baggage restrictions server
	HostPort string `json:"host_port"`

	// RefreshInterval controls how often the baggage restriction manager will poll
	// jaeger-agent for the most recent baggage restrictions.
	RefreshInterval time.Duration `json:"refresh_interval"`
}

// ThrottlerConfig configures the throttler which can be used to throttle the
// rate at which the client may send debug requests.
// See https://godoc.org/github.com/uber/jaeger-client-go/config#ThrottlerConfig
type ThrottlerConfig struct {
	// HostPort of jaeger-agent's credit server.
	HostPort string `json:"host_port"`

	// RefreshInterval controls how often the throttler will poll jaeger-agent
	// for more throttling credits.
	RefreshInterval time.Duration `json:"refresh_interval"`

	// SynchronousInitialization determines whether or not the throttler should
	// synchronously fetch credits from the agent when an operation is seen for
	// the first time. This should be set to true if the client will be used by
	// a short lived service that needs to ensure that credits are fetched
	// upfront such that sampling or throttling occurs.
	SynchronousInitialization bool `json:"synchronous_initialization"`
}

// ToTracingConfig converts *OpenTracing to *tracing.Configuration.
func (c *Config) ToTracingConfig() *config.Configuration {
	ret := &config.Configuration{
		ServiceName: c.ServiceName,
		Disabled:    false,
		RPCMetrics:  c.RPCMetrics,
		Gen128Bit:   false,
		Tags:        []opentracing.Tag{},
	}
	if c.Sampler != nil {
		ret.Sampler = &config.SamplerConfig{
			Type:                     c.Sampler.Type,
			Param:                    c.Sampler.Param,
			SamplingServerURL:        c.Sampler.SamplingServerURL,
			SamplingRefreshInterval:  c.Sampler.SamplingRefreshInterval,
			MaxOperations:            c.Sampler.MaxOperations,
			OperationNameLateBinding: c.Sampler.OperationNameLateBinding,
		}
	}
	if c.Reporter != nil {
		ret.Reporter = &config.ReporterConfig{
			QueueSize:                  c.Reporter.QueueSize,
			BufferFlushInterval:        c.Reporter.BufferFlushInterval,
			LogSpans:                   c.Reporter.LogSpans,
			LocalAgentHostPort:         c.Reporter.LocalAgentHostPort,
			DisableAttemptReconnecting: c.Reporter.DisableAttemptReconnecting,
			AttemptReconnectInterval:   c.Reporter.AttemptReconnectInterval,
			CollectorEndpoint:          c.Reporter.CollectorEndpoint,
			User:                       c.Reporter.User,
			Password:                   c.Reporter.Password,
			HTTPHeaders:                c.Reporter.HTTPHeaders,
		}
	}
	if c.Headers != nil {
		ret.Headers = &jaeger.HeadersConfig{
			JaegerDebugHeader:        c.Headers.JaegerDebugHeader,
			JaegerBaggageHeader:      c.Headers.JaegerBaggageHeader,
			TraceContextHeaderName:   c.Headers.TraceContextHeaderName,
			TraceBaggageHeaderPrefix: c.Headers.TraceBaggageHeaderPrefix,
		}
	}
	if c.BaggageRestrictions != nil {
		ret.BaggageRestrictions = &config.BaggageRestrictionsConfig{
			DenyBaggageOnInitializationFailure: c.BaggageRestrictions.DenyBaggageOnInitializationFailure,
			HostPort:                           c.BaggageRestrictions.HostPort,
			RefreshInterval:                    c.BaggageRestrictions.RefreshInterval,
		}
	}
	if c.Throttler != nil {
		ret.Throttler = &config.ThrottlerConfig{
			HostPort:                  c.Throttler.HostPort,
			RefreshInterval:           c.Throttler.RefreshInterval,
			SynchronousInitialization: c.Throttler.SynchronousInitialization,
		}
	}

	return ret
}
