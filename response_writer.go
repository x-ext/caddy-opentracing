package opentracing

import (
	"net/http"
)

type metricsTracker struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *metricsTracker) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *metricsTracker) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}
