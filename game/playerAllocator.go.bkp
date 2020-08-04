package game

import "sync"

// Allocator as as
type Allocator struct {
	sync.Mutex
	list []string
	mp   map[string]chan string
}

// NewAllocator returns a pointer to a new allocator
func NewAllocator() *Allocator {
	return &Allocator{mp: make(map[string]chan string)}
}

// Find is used to
func (all *Allocator) Find(name string, output chan string) {
	all.Lock()
	defer all.Unlock()
	if len(all.list) != 0 {
		otherUser := all.list[0]
		all.list = all.list[1:]
		otherUsersChannel := all.mp[otherUser]
		all.mp[otherUser] = nil
		otherUsersChannel <- name
		output <- otherUser
	} else {
		all.list = append(all.list, name)
		all.mp[name] = output
	}
}

func (all *Allocator) IDontNeedAnyMore(name string) {
	all.Lock()
	defer all.Unlock()

	if all.mp[name] == nil {
	} else {
		close(all.mp[name])
	}

}
