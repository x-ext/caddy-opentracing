package opentracing

import (
	"io"
	"net/http"
	"net/url"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go/config"
)

func init() {
	caddy.RegisterModule(Opentracing{})
}

const (
	defaultComponentName = "caddy.module.opentracing"
	defaultServiceName   = "caddy"
	responseSizeKey      = "http.response_size"
)

type Opentracing struct {
	Config
	tr     opentracing.Tracer
	opts   Options
	closer io.Closer
}

// Validate implements caddy.Validator.
func (tracing *Opentracing) Validate() (err error) {
	return nil
}

// CaddyModule returns the Caddy module information.
func (tracing Opentracing) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.opentracing",
		New: func() caddy.Module {
			return new(Opentracing)
		},
	}
}

// Implements caddy.Provisioner.
func (tracing *Opentracing) Provision(ctx caddy.Context) (err error) {
	var cfg *config.Configuration
	if cfg, err = tracing.Config.ToTracingConfig().FromEnv(); err != nil {
		return
	}

	if cfg.ServiceName == "" {
		cfg.ServiceName = defaultServiceName
	}

	if tracing.tr, tracing.closer, err = cfg.New(tracing.Config.ServiceName); err != nil {
		return
	}

	tracing.opts = Options{
		opNameFunc: func(r *http.Request) string {
			return r.Method + " " + r.URL.Path
		},
		spanFilter:   func(r *http.Request) bool { return true },
		spanObserver: func(span opentracing.Span, r *http.Request) {},
		urlTagFunc: func(u *url.URL) string {
			return u.String()
		},
	}
	return nil
}

type Options struct {
	opNameFunc    func(r *http.Request) string
	spanFilter    func(r *http.Request) bool
	spanObserver  func(span opentracing.Span, r *http.Request)
	urlTagFunc    func(u *url.URL) string
	componentName string
}

func (tracing Opentracing) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) (err error) {
	tr := tracing.tr
	opts := tracing.opts

	componentName := opts.componentName
	if componentName == "" {
		componentName = defaultComponentName
	}

	if !opts.spanFilter(r) {
		return next.ServeHTTP(w, r)
	}

	ctx, _ := tr.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	sp := tr.StartSpan(opts.opNameFunc(r), ext.RPCServerOption(ctx))
	ext.HTTPMethod.Set(sp, r.Method)
	ext.HTTPUrl.Set(sp, opts.urlTagFunc(r.URL))
	ext.Component.Set(sp, componentName)
	opts.spanObserver(sp, r)

	r = r.WithContext(opentracing.ContextWithSpan(r.Context(), sp))
	mt := &metricsTracker{ResponseWriter: w}

	err = next.ServeHTTP(mt, r)
	if mt.status > 0 {
		ext.HTTPStatusCode.Set(sp, uint16(mt.status))
	}
	if mt.size > 0 {
		sp.SetTag(responseSizeKey, mt.size)
	}
	if mt.status >= http.StatusInternalServerError {
		ext.Error.Set(sp, true)
	}
	sp.Finish()
	return err
}
