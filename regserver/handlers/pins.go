package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/qri-io/apiutil"
	"github.com/qri-io/registry/pinset"
)

// NewPinsHandler creates a profiles handler function that operates
// on a *registry.Profiles
func NewPinsHandler(ps pinset.Pinset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var status pinset.PinStatus

		req, err := parsePinReq(r)
		if err != nil {
			apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
			return
		}

		switch r.Method {
		case "POST":
			statusChan, err := ps.Pin(req)
			if err != nil {
				apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
				return
			}
			status = <-statusChan
			apiutil.WriteResponse(w, status)
			return
		case "DELETE":
			if err = ps.Unpin(req); err != nil {
				apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
			}
		case "GET":
			p := apiutil.PageFromRequest(r)
			pins, err := ps.Pins(p.Limit(), p.Offset())
			if err != nil {
				apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
				return
			}
			apiutil.WriteResponse(w, pins)
		default:
			apiutil.NotFoundHandler(w, r)
		}
	}
}

// NewPinStatusHandler creates a handler for getting the pin status of a hash
func NewPinStatusHandler(ps pinset.Pinset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := parsePinReq(r)
		if err != nil {
			apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
			return
		}

		status, err := ps.Status(req)
		if err != nil {
			if err.Error() == "not found" {
				apiutil.WriteErrResponse(w, http.StatusNotFound, err)
			} else {
				apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
			}
			return
		}

		apiutil.WriteResponse(w, status)
	}
}

func parsePinReq(r *http.Request) (req *pinset.PinRequest, err error) {
	req = &pinset.PinRequest{}

	switch r.Header.Get("Content-Type") {
	case "application/json":
		err = json.NewDecoder(r.Body).Decode(req)
	default:
		req.Path = r.FormValue("path")
	}

	return
}
