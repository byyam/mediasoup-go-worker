package room

import (
	"sync"
)

var once sync.Once

var roomsCollection *RoomsCollection

type RoomsCollection struct {
	sync.RWMutex
	rooms map[string]*Room
}

func init() {
	once.Do(func() {
		roomsCollection = &RoomsCollection{
			rooms: make(map[string]*Room),
		}
	})
}

type Room struct {
}

func GetOrCreateRoom(roomId string) *Room {
	roomsCollection.RLock()
	defer roomsCollection.RUnlock()

	r, ok := roomsCollection.rooms[roomId]
	if ok {
		return r
	}

	return &Room{}
}
