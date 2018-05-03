package registry

import (
	"sort"
	"sync"
)

// Datasets is a map of datasets data safe for concurrent use
// heavily inspired by sync.Map
type Datasets struct {
	sync.RWMutex
	tips     map[string]string
	internal map[string]*Dataset
}

// NewDatasets allocates a new *Datasets map
func NewDatasets() *Datasets {
	return &Datasets{
		tips:     make(map[string]string),
		internal: make(map[string]*Dataset),
	}
}

// Len returns the number of records in the map
func (ds *Datasets) Len() int {
	return len(ds.internal)
}

// Load fetches a dataset from the list by key
func (ds *Datasets) Load(key string) (value *Dataset, ok bool) {
	ds.RLock()
	result, ok := ds.internal[key]
	ds.RUnlock()
	return result, ok
}

// Range calls an iteration fuction on each element in the map until
// the end of the list is reached or iter returns true
func (ds *Datasets) Range(iter func(key string, p *Dataset) (brk bool)) {
	ds.RLock()
	defer ds.RUnlock()
	for key, p := range ds.internal {
		if iter(key, p) {
			break
		}
	}
}

// SortedRange is like range but with deterministic key ordering
func (ds *Datasets) SortedRange(iter func(key string, p *Dataset) (brk bool)) {
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

// Delete removes a record from Datasets at key
func (ds *Datasets) Delete(key string) {
	ds.Lock()
	delete(ds.internal, key)
	ds.Unlock()
}

// Store adds an entry
func (ds *Datasets) Store(key string, value *Dataset) {
	ds.Lock()
	ds.internal[key] = value
	ds.Unlock()
}

// Register adds a datasets to the list if it's valid and the desired handle isn't taken
func (ds *Datasets) Register(d *Dataset) error {
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
func (ds *Datasets) Deregister(d *Dataset) error {
	if err := d.Validate(); err != nil {
		return err
	}
	if err := d.Verify(); err != nil {
		return err
	}

	ds.Delete(d.Key())
	return nil
}
