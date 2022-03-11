package utils

import "sync"

type Hashmap struct {
	sync.RWMutex
	map1 map[interface{}]map[interface{}]interface{}
}

func NewHashMap() *Hashmap {
	return &Hashmap{
		map1: make(map[interface{}]map[interface{}]interface{}),
	}
}

func (h *Hashmap) Store(id1, id2, data interface{}) {
	if data == nil {
		return
	}
	h.Lock()
	defer h.Unlock()

	map2, ok := h.map1[id1]
	if ok {
		map2[id2] = data
	} else {
		map2 := make(map[interface{}]interface{})
		map2[id2] = data
		h.map1[id1] = map2
	}
}

func (h *Hashmap) Load(id1, id2 interface{}) (interface{}, bool) {
	h.RLock()
	defer h.RUnlock()

	map2, ok := h.map1[id1]
	if !ok || map2 == nil {
		return nil, false
	}
	data, ok := map2[id2]
	if !ok {
		return nil, false
	}
	return data, true
}

func (h *Hashmap) Erase(id1 interface{}) {
	h.Lock()
	defer h.Unlock()

	delete(h.map1, id1)
}

func (h *Hashmap) Delete(id1, id2 interface{}) {
	h.Lock()
	defer h.Unlock()

	map2, ok := h.map1[id1]
	if !ok || map2 == nil {
		return
	}
	delete(map2, id2)
}

func (h *Hashmap) Get(id1 interface{}) (interface{}, bool) {
	h.RLock()
	defer h.RUnlock()

	map2, ok := h.map1[id1]
	if !ok || map2 == nil {
		return nil, false
	}
	return map2, ok
}
