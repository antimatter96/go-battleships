package stuff

import (
	"sync"
)

// StringSet is a set of strings
type stringSet struct {
	sync.Mutex
	mp map[string]bool
}

// Add - adds a string to the set
// Will return false if already present
func (set *stringSet) add(s string) bool {
	set.Lock()
	defer set.Unlock()
	if set.mp == nil {
		set.mp = make(map[string]bool)
	}
	if _, present := set.mp[s]; present {
		return false
	}
	set.mp[s] = true
	return true
}

// Has - checks for string in the set
func (set *stringSet) has(s string) bool {
	set.Lock()
	defer set.Unlock()
	if set.mp == nil {
		return false
	}
	if _, present := set.mp[s]; present {
		return true
	}
	return false
}

// Delete - deletes the string from the set
// Will return true if was present and deleted
func (set *stringSet) delete(s string) bool {
	set.Lock()
	defer set.Unlock()
	if set.mp == nil {
		return false
	}
	if _, present := set.mp[s]; !present {
		return false
	}
	delete(set.mp, s)
	return true
}

// Size - returns number of elements in the set
func (set *stringSet) size() int {
	set.Lock()
	defer set.Unlock()
	if set.mp == nil {
		return 0
	}
	return len(set.mp)
}
