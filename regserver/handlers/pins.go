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
		var status registry.PinStatus

		req, err := parsePinReq(r)
		if err != nil {
			apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
			return
		}

		switch r.Method {
		case "POST":
			if statusChan, err := pinset.Pin(req); err != nil {
				apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
				return
			} else {
				status = <-statusChan
			}
			apiutil.WriteResponse(w, status)
		case "DELETE":
			if err = pinset.Unpin(req); err != nil {
				apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
			}
		case "GET":
			p := apiutil.PageFromRequest(r)
			pins, err := pinset.Pins(p.Limit(), p.Offset())
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
func NewPinStatusHandler(pinset registry.Pinset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := parsePinReq(r)
		if err != nil {
			apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
			return
		}

		status, err := pinset.Status(req)
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

func parsePinReq(r *http.Request) (req *registry.PinRequest, err error) {
	req = &registry.PinRequest{}

	switch r.Header.Get("Content-Type") {
	case "application/json":
		err = json.NewDecoder(r.Body).Decode(req)
	default:
		req.Path = r.FormValue("path")
	}

	return
}
