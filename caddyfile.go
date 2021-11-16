package opentracing

import (
	"fmt"
	"strconv"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	httpcaddyfile.RegisterHandlerDirective("opentracing", parseCaddyfile)
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	tracing := new(Opentracing)
	err := tracing.UnmarshalCaddyfile(h.Dispenser)
	return tracing, err
}

// UnmarshalCaddyfile sets up the handler from Caddyfile tokens. Syntax:
// Specifying the formats on the first line will use those formats' defaults.
func (tracing *Opentracing) UnmarshalCaddyfile(d *caddyfile.Dispenser) (err error) {
	var cfg Config
	for d.Next() {
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "service_name":
				if !d.NextArg() {
					return d.ArgErr()
				}
				cfg.ServiceName = d.Val()
			case "disabled":
				cfg.Disabled = true
			case "rpc_metrics":
				cfg.RPCMetrics = true
			case "traceid_128bit":
				cfg.Gen128Bit = true
			case "sampler":
				cfg.Sampler = new(SamplerConfig)
				for nesting := d.Nesting(); d.NextBlock(nesting); {
					switch d.Val() {
					case "type":
						if !d.NextArg() {
							return d.ArgErr()
						}
						cfg.Sampler.Type = d.Val()
					case "param":
						if !d.NextArg() {
							return d.ArgErr()
						}
						if cfg.Sampler.Param, err = strconv.ParseFloat(d.Val(), 64); err != nil {
							return
						}
					case "sampling_server_url":
						if !d.NextArg() {
							return d.ArgErr()
						}
						cfg.Sampler.SamplingServerURL = d.Val()
					case "sampling_refresh_interval":
						if !d.NextArg() {
							return d.ArgErr()
						}
						if cfg.Sampler.SamplingRefreshInterval, err = time.ParseDuration(d.Val()); err != nil {
							return
						}
					case "max_operations":
						if !d.NextArg() {
							return d.ArgErr()
						}
						if cfg.Sampler.MaxOperations, err = strconv.Atoi(d.Val()); err != nil {
							return
						}
					case "operation_name_late_binding":
						cfg.Sampler.OperationNameLateBinding = true
					}
				}
			case "reporter":
				cfg.Reporter = new(ReporterConfig)
				for nesting := d.Nesting(); d.NextBlock(nesting); {
					switch d.Val() {
					case "collector_endpoint":
						if !d.NextArg() {
							return d.ArgErr()
						}
						cfg.Reporter.CollectorEndpoint = d.Val()
					case "user":
						if !d.NextArg() {
							return d.ArgErr()
						}
						cfg.Reporter.User = d.Val()
					case "password":
						if !d.NextArg() {
							return d.ArgErr()
						}
						cfg.Reporter.Password = d.Val()
					case "local_agent_host_port":
						if !d.NextArg() {
							return d.ArgErr()
						}
						cfg.Reporter.LocalAgentHostPort = d.Val()
					case "buffer_flush_interval":
						if !d.NextArg() {
							return d.ArgErr()
						}
						if cfg.Reporter.BufferFlushInterval, err = time.ParseDuration(d.Val()); err != nil {
							return
						}
					case "attempt_reconnect_interval":
						if !d.NextArg() {
							return d.ArgErr()
						}
						if cfg.Reporter.AttemptReconnectInterval, err = time.ParseDuration(d.Val()); err != nil {
							return
						}
					case "queue_size":
						if !d.NextArg() {
							return d.ArgErr()
						}
						if cfg.Reporter.QueueSize, err = strconv.Atoi(d.Val()); err != nil {
							return
						}
					case "log_spans":
						cfg.Reporter.LogSpans = true
					case "disable_attempt_reconnecting":
						cfg.Reporter.DisableAttemptReconnecting = true
					case "http_headers":
					}
				}
			case "headers":
				cfg.Headers = new(HeadersConfig)
				for nesting := d.Nesting(); d.NextBlock(nesting); {
					switch d.Val() {
					case "jaeger_debug_header":
						if !d.NextArg() {
							return d.ArgErr()
						}
						cfg.Headers.JaegerDebugHeader = d.Val()
					case "jaeger_baggage_header":
						if !d.NextArg() {
							return d.ArgErr()
						}
						cfg.Headers.JaegerBaggageHeader = d.Val()
					case "trace_context_header_name":
						if !d.NextArg() {
							return d.ArgErr()
						}
						cfg.Headers.TraceContextHeaderName = d.Val()
					case "trace_baggage_header_prefix":
						if !d.NextArg() {
							return d.ArgErr()
						}
						cfg.Headers.TraceBaggageHeaderPrefix = d.Val()
					}
				}
			case "baggage_restrictions":
				cfg.BaggageRestrictions = new(BaggageRestrictionsConfig)
				for nesting := d.Nesting(); d.NextBlock(nesting); {
					switch d.Val() {
					case "deny_baggage_on_initialization_failure":
						cfg.BaggageRestrictions.DenyBaggageOnInitializationFailure = true
					case "host_port":
						if !d.NextArg() {
							return d.ArgErr()
						}
						cfg.BaggageRestrictions.HostPort = d.Val()
					case "refresh_interval":
						if !d.NextArg() {
							return d.ArgErr()
						}
						if cfg.BaggageRestrictions.RefreshInterval, err = time.ParseDuration(d.Val()); err != nil {
							return
						}
					}
				}
			case "throttler":
				cfg.Throttler = new(ThrottlerConfig)
				for nesting := d.Nesting(); d.NextBlock(nesting); {
					switch d.Val() {
					case "synchronous_initialization":
						cfg.Throttler.SynchronousInitialization = true
					case "host_port":
						if !d.NextArg() {
							return d.ArgErr()
						}
						cfg.Throttler.HostPort = d.Val()
					case "refresh_interval":
						if !d.NextArg() {
							return d.ArgErr()
						}
						if cfg.Throttler.RefreshInterval, err = time.ParseDuration(d.Val()); err != nil {
							return
						}
					}
				}
			}
		}
	}
	tracing.Config = cfg
	fmt.Printf("%+v", cfg.Reporter)
	return nil
}

// Interface guard
var (
	_ caddyfile.Unmarshaler       = (*Opentracing)(nil)
	_ caddyhttp.MiddlewareHandler = (*Opentracing)(nil)
	_ caddy.Validator             = (*Opentracing)(nil)
)
