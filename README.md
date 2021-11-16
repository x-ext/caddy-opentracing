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

```caddy
{
	auto_https off
	http_port 80
	https_port 443
}

:80 {
	route /* {
		opentracing {
			service_name hello
			reporter {
				local_agent_host_port localhost:6831
                queue_size 1
			}
		}
		reverse_proxy https://baidu.com
	}
	# curl localhost:80/index.html
}

```
