package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/datatogether/api/apiutil"
	"github.com/qri-io/registry"
)

// NewPinsHandler creates a profiles handler function that operates
// on a *registry.Profiles
func NewPinsHandler(pinset registry.Pinset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			req = &registry.PinRequest{}
			err error
		)

		switch r.Header.Get("Content-Type") {
		case "application/json":
			if err = json.NewDecoder(r.Body).Decode(req); err != nil {
				apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
				return
			}
		default:
			req.Path = r.FormValue("path")
		}

		switch r.Method {
		case "POST":
			if _, err = pinset.Pin(req); err != nil {
				apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
				return
			}
		case "DELETE":
			if err = pinset.Unpin(req); err != nil {
				apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
				return
			}
		case "GET":
			if req.Pinned, err = pinset.Pinned(req.Path); err != nil {
				apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
				return
			}
		}
		apiutil.WriteResponse(w, req)
	}
}
