package http

import net_http "net/http"

type Controller interface {
	Router() net_http.Handler
}
