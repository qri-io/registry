package registry

import (
	"sort"
	"sync"
)

// Datasets is the interface for working with a set of *Dataset's
// Register, Deregister, Load, Len, Range, and SortedRange should be
// considered safe to hook up to public http endpoints, whereas
// Delete & Store should only be exposed in administrative contexts
type Datasets interface {
	// Register adds a dataset to the set if it's valid
	Register(d *Dataset) error
	// Deregister removes a dataset from the registry if it exists
	Deregister(d *Dataset) error

	// Len returns the number of records in the set
	Len() int
	// Load fetches a dataset from the list by key
	Load(key string) (value *Dataset, ok bool)
	// Range calls an iteration fuction on each element in the map until
	// the end of the list is reached or iter returns true
	Range(iter func(key string, p *Dataset) (brk bool))
	// SortedRange is like range but with deterministic key ordering
	SortedRange(iter func(key string, p *Dataset) (brk bool))

	// Store adds an entry, bypassing the register process
	// store is only exported for administrative use cases.
	// most of the time callers should use Register instead
	Store(key string, value *Dataset)
	// Delete removes a record stored at key
	// Delete is only exported for administrative use cases.
	// most of the time callers should use Deregister instead
	Delete(key string)
}

// MemDatasets is a map of datasets data safe for concurrent use
// heavily inspired by sync.Map
type MemDatasets struct {
	sync.RWMutex
	tips     map[string]string
	internal map[string]*Dataset
}

// NewMemDatasets allocates a new *MemDatasets map
func NewMemDatasets() *MemDatasets {
	return &MemDatasets{
		tips:     make(map[string]string),
		internal: make(map[string]*Dataset),
	}
}

// Len returns the number of records in the map
func (ds *MemDatasets) Len() int {
	return len(ds.internal)
}

// Load fetches a dataset from the list by key
func (ds *MemDatasets) Load(key string) (value *Dataset, ok bool) {
	ds.RLock()
	defer ds.RUnlock()
	value, ok = ds.internal[key]
	return
}

// Range calls an iteration fuction on each element in the map until
// the end of the list is reached or iter returns true
func (ds *MemDatasets) Range(iter func(key string, p *Dataset) (brk bool)) {
	ds.RLock()
	defer ds.RUnlock()
	for key, p := range ds.internal {
		if iter(key, p) {
			break
		}
	}
}

// SortedRange is like range but with deterministic key ordering
func (ds *MemDatasets) SortedRange(iter func(key string, p *Dataset) (brk bool)) {
	ds.RLock()
	defer ds.RUnlock()
	keys := make([]string, len(ds.internal))
	i := 0
	for key := range ds.internal {
		keys[i] = key
		i++
	}
	sort.StringSlice(keys).Sort()
	for _, key := range keys {
		if iter(key, ds.internal[key]) {
			break
		}
	}
}

// Delete removes a record from MemDatasets at key
func (ds *MemDatasets) Delete(key string) {
	ds.Lock()
	delete(ds.internal, key)
	ds.Unlock()
}

// Store adds an entry
func (ds *MemDatasets) Store(key string, value *Dataset) {
	ds.Lock()
	ds.internal[key] = value
	ds.Unlock()
}

// Register adds a dataset to the set if it's valid
func (ds *MemDatasets) Register(d *Dataset) error {
	if err := d.Validate(); err != nil {
		return err
	}
	if err := d.Verify(); err != nil {
		return err
	}

	prev := ""
	dkey := d.Key()
	ds.Range(func(key string, d *Dataset) bool {
		if key == dkey {
			prev = key
			return true
		}
		return false
	})

	if prev != "" {
		ds.Delete(prev)
	}

	ds.Store(dkey, d)
	return nil
}

// Deregister removes a datasets from the registry if it exists
func (ds *MemDatasets) Deregister(d *Dataset) error {
	if err := d.Validate(); err != nil {
		return err
	}
	if err := d.Verify(); err != nil {
		return err
	}

	ds.Delete(d.Key())
	return nil
}
