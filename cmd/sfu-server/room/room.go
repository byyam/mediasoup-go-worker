package room

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
)

var manager *sessionManager

var logger zerolog.Logger

type sessionManager struct {
	rooms map[string]*Room
}

type Room struct {
	Peers map[string]*Peer
}

type Peer struct {
	Joined                     bool // websocket create and join not atomic processing
	ProducingWebRtcTransportId string
	ConsumingWebRtcTransportId string
}

func InitSessionManager() {
	manager = &sessionManager{
		rooms: make(map[string]*Room),
	}
	logger = zerowrapper.NewScope("session")
	logger.Info().Msgf("session manager initialized")
}

func NewSession(roomId, peerId string, onCreateRoom func()) error {
	room := getOrCreateRoom(roomId, onCreateRoom)
	p, ok := room.Peers[peerId]
	if ok {
		return fmt.Errorf("peer %s already exists", peerId)
	}

	p = newPeer(room)
	room.Peers[peerId] = p

	logger.Info().Str("roomId", roomId).Str("peerId", peerId).Msgf("new seesion created")
	return nil
}

func CloseSession(roomId, peerId string, onCloseRoom func()) error {
	r, ok := manager.rooms[roomId]
	if !ok {
		return fmt.Errorf("room %s not found", roomId)
	}

	p, ok := r.Peers[peerId]
	if !ok {
		return fmt.Errorf("peer %s not found", peerId)
	}
	p.Close()
	logger.Info().Msgf("closed peer %s", peerId)

	// check room empty
	delete(r.Peers, peerId)
	if len(r.Peers) == 0 {
		r.Close()
		delete(manager.rooms, roomId)
		logger.Info().Msgf("closed empty room %s", roomId)
	}
	if onCloseRoom != nil {
		onCloseRoom()
	}

	logger.Info().Str("roomId", roomId).Str("peerId", peerId).Msgf("seesion closed")
	return nil
}

func newRoom() *Room {
	return &Room{
		Peers: make(map[string]*Peer),
	}
}

func (r *Room) Close() {

}

func getOrCreateRoom(roomId string, onCreateRoom func()) *Room {
	r, ok := manager.rooms[roomId]
	if ok {
		return r
	}

	r = newRoom()
	manager.rooms[roomId] = r
	if onCreateRoom != nil {
		onCreateRoom()
	}

	return r
}

func newPeer(r *Room) *Peer {
	return &Peer{}
}

func (p *Peer) Close() {}

func GetRoom(roomId string) (*Room, error) {
	r, ok := manager.rooms[roomId]
	if !ok {
		return nil, fmt.Errorf("room %s not found", roomId)
	}

	return r, nil
}

func GetPeer(roomId, peerId string) (*Peer, error) {
	r, err := GetRoom(roomId)
	if err != nil {
		return nil, err
	}

	p, ok := r.Peers[peerId]
	if !ok {
		return nil, fmt.Errorf("peer %s not found", peerId)
	}

	return p, nil
}
