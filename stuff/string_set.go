package stuff

import (
	"errors"
	"sync"
)

// StringSet is a set of strings
type StringSet struct {
	sync.Mutex
	mp map[string]bool
}

// ErrorAlreadyPresent is returned when
var (
	ErrorAlreadyPresent = errors.New("Element already present")
)

// Add - adds a string to the set
// Will return an error if already present
func (set *StringSet) Add(s string) bool {
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
func (set *StringSet) Has(s string) bool {
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
func (set *StringSet) Delete(s string) bool {
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

func (set *StringSet) Size() int {
	set.Lock()
	defer set.Unlock()
	if set.mp == nil {
		return 0
	}
	return len(set.mp)
}
