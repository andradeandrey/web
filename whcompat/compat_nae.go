// +build !appengine
// +build !go1.7

package whcompat

import (
	"net/http"

	"golang.org/x/net/context"
)

func new16Context(r *http.Request) context.Context {
	return context.Background()
}
