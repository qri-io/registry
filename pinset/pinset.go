package pinset

import (
	"fmt"
	"sort"

	"github.com/qri-io/registry"
)

// Pinset is the interface for acting as a remote pinning service.
// Pinset implementations are expected to keep a store of PinStatus
// that callers can use to probe the progress of a request
type Pinset interface {
	// Pin can take a while, so Pin returns a channel of PinStatus to
	// deliver updates structs that should all share an ID for the request that
	// the Pinset implementation will provide
	Pin(req *PinRequest) (chan PinStatus, error)
	// Unpin removes a pin
	Unpin(req *PinRequest) error
	// Status gets the current pin state value for a given PinRequest
	Status(req *PinRequest) (PinStatus, error)
	// Pins lists pins within the range defined by limit & offset in
	// lexographical order, smallest to largest
	Pins(limit, offset int) ([]string, error)
	// PinLen returns the number of pins in the set
	PinLen() (int, error)
}

// MemPinset is a completely ficticious implementation of a pinset
// it shouldn't ever be used in real-world scenarios. We use it for mocking
// a pinning service without an actual backing store keeping pins
type MemPinset struct {
	pk       PinStatusStore
	Profiles registry.Profiles
	pins     []string
}

func insertSorted(list []string, elem string) []string {
	index := sort.Search(len(list), func(i int) bool { return list[i] > elem })
	list = append(list, "")
	copy(list[index + 1:], list[index:])
	list[index] = elem
	return list
}

// Pin a dataset
func (m *MemPinset) Pin(req *PinRequest) (chan PinStatus, error) {
	m.pins = insertSorted(m.pins, req.Path)

	pc := make(chan PinStatus)
	go func() {
		prog := PinStatus{
			Path:        req.Path,
			PctComplete: 1.0,
			Pinned:      true,
		}

		m.pk.Set(prog)
		pc <- prog
		close(pc)
	}()
	return pc, nil
}

// Status gives a hydrated, up-to-date progress struct for a given request ID
func (m *MemPinset) Status(req *PinRequest) (PinStatus, error) {
	ps := m.pk.Get(req.Path)
	if ps == nil {
		return PinStatus{}, fmt.Errorf("not found")
	}

	return *ps, nil
}

// Unpin a dataset
func (m *MemPinset) Unpin(req *PinRequest) error {
	for i, p := range m.pins {
		if p == req.Path {
			m.pk.Delete(req.Path)
			m.pins = append(m.pins[:i], m.pins[i+1:]...)
			return nil
		}
	}
	return nil
}

// Pinned gets the pin status of a path
func (m *MemPinset) Pinned(path string) (pinned bool, err error) {
	for _, p := range m.pins {
		if p == path {
			return true, nil
		}
	}
	return false, nil
}

// Pins reads from the list present in the pinset
func (m *MemPinset) Pins(limit, offset int) (pins []string, err error) {
	for i, p := range m.pins {
		if i < offset {
			continue
		}
		pins = append(pins, p)
		if len(pins) == limit {
			break
		}
	}
	return
}

// PinLen returns the number of pins in the pinset
func (m *MemPinset) PinLen() (int, error) {
	return len(m.pins), nil
}
