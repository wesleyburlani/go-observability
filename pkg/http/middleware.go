package http

import net_http "net/http"

type Middleware interface {
	Handle(next net_http.Handler) net_http.Handler
}
