package web_test

import (
	"fmt"
	"net/http"

	"github.com/andradeandrey/web/whcompat"
	"github.com/andradeandrey/web/whlog"
	"github.com/andradeandrey/web/whmux"
)

var (
	pageName = whmux.NewStringArg()
)

func page(w http.ResponseWriter, r *http.Request) {
	name := pageName.Get(whcompat.Context(r))

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Welcome to %s", name)
}

func Example() {
	pageHandler := pageName.Shift(whmux.Exact(http.HandlerFunc(page)))

	whlog.ListenAndServe(":0", whmux.Dir{
		"wiki": pageHandler,
	})
}
