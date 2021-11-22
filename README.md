# caddy-opentracing

---

Enable requests served by caddy for distributed tracing via [The OpenTracing Project](http://opentracing.io).

## Dependencies

- The [Go OpenTracing Library](https://github.com/opentracing/opentracing-go)
  [Jaeger](https://github.com/jaegertracing/cpp-client),

## Getting Started

First, write a configuration for the tracer used. Below's an example of what
a Jaeger configuration might look like:

Caddyfile

```shell
{
	auto_https off
	http_port 80
	https_port 443
}

:80 {
	route /* {
		opentracing {
			# Can be provided by FromEnv() via the environment variable named JAEGER_SERVICE_NAME
			service_name hello #default caddy
			# Value can be provided by FromEnv() via the environment variable named JAEGER_DISABLED.
			disable
			# Value can be provided by FromEnv() via the environment variable named JAEGER_RPC_METRICS
			rpc_metrics
			# Gen128Bit instructs the tracer to generate 128-bit wide trace IDs, compatible with W3C Trace Context.
			traceid_128bit
			# See https://pkg.go.dev/github.com/uber/jaeger-client-go/config#ReporterConfig
			reporter {
				local_agent_host_port localhost:6831
				queue_size 1
			}
			# See https://pkg.go.dev/github.com/uber/jaeger-client-go/config#SamplerConfig
			sampler {
				type const
			}
			# ...
		}
		reverse_proxy https://baidu.com
	}
	# curl localhost:80/abc/pub/example/imap-console-client.png
}

```
