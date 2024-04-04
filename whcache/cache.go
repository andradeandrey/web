package whcache // import "github.com/andradeandrey/web/whcache"

import (
	"net/http"

	"github.com/andradeandrey/web"
	"github.com/andradeandrey/web/whcompat"
	"github.com/andradeandrey/web/whroute"
	"golang.org/x/net/context"
)

var (
	cacheKey = webhelp.GenSym()
)

type reqCache map[interface{}]interface{}

// Register installs a cache in the handler chain.
func Register(h http.Handler) http.Handler {
	return whroute.HandlerFunc(h, func(w http.ResponseWriter, r *http.Request) {
		ctx := whcompat.Context(r)
		if _, ok := ctx.Value(cacheKey).(reqCache); ok {
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, whcompat.WithContext(r,
			context.WithValue(ctx, cacheKey, reqCache{})))
	})
}

// Set stores the key/val pair in the context specific cache, if possible.
func Set(ctx context.Context, key, val interface{}) {
	cache, ok := ctx.Value(cacheKey).(reqCache)
	if !ok {
		return
	}
	cache[key] = val
}

// Remove removes any values stored with key from the context specific cache,
// if possible.
func Remove(ctx context.Context, key interface{}) {
	cache, ok := ctx.Value(cacheKey).(reqCache)
	if !ok {
		return
	}
	delete(cache, key)
}

// Get returns previously stored key/value pairs from the context specific
// cache if one is registered and the value is found, and returns nil
// otherwise.
func Get(ctx context.Context, key interface{}) interface{} {
	cache, ok := ctx.Value(cacheKey).(reqCache)
	if !ok {
		return nil
	}
	val, _ := cache[key]
	return val
}
