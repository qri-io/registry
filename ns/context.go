package ns

import (
	"context"
	"net/http"
	"strings"
)

// QriCtxKey defines a distinct type for
// keys for context values should always use custom
// types to avoid collissions.
// see comment on context.WithValue for more info
type QriCtxKey string

// RefCtxKey is the key for adding a dataset reference
// to a context.Context
const RefCtxKey QriCtxKey = "datasetRef"

// RefFromReq examines the path element of a request URL
// to
func RefFromReq(r *http.Request) (Ref, error) {
	if r.URL.String() == "" || r.URL.Path == "" {
		return Ref{}, nil
	}
	return RefFromHTTPPath(r.URL.Path)
}

// RefFromHTTPPath parses a path and returns a datasetRef
func RefFromHTTPPath(path string) (Ref, error) {
	refstr := HTTPPathToQriPath(path)
	return ParseRef(refstr)
}

// RefFromCtx extracts a Dataset reference from a given
// context if one is set, returning nil otherwise
func RefFromCtx(ctx context.Context) Ref {
	iface := ctx.Value(RefCtxKey)
	if ref, ok := iface.(Ref); ok {
		return ref
	}
	return Ref{}
}

// HTTPPathToQriPath converts a http path to a
// qri path
func HTTPPathToQriPath(path string) string {
	paramIndex := strings.Index(path, "?")
	if paramIndex != -1 {
		path = path[:paramIndex]
	}
	path = strings.Replace(path, "/at", "@", 1)
	path = strings.TrimPrefix(path, "/")
	return path
}
