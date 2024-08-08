package iceutil

import (
	"net"

	"github.com/pion/ice/v3"
	"github.com/pion/stun/v2"
)

func AssertInboundUsername(m *stun.Message, expectedUsername string) error {
	var username stun.Username
	if err := username.GetFrom(m); err != nil {
		return err
	}
	//if string(username) != expectedUsername {
	//	return fmt.Errorf("mismatch stun username: expected(%x) actual(%x)", expectedUsername, string(username))
	//}

	return nil
}

func AssertInboundMessageIntegrity(m *stun.Message, key []byte) error {
	messageIntegrityAttr := stun.MessageIntegrity(key)
	return messageIntegrityAttr.Check(m)
}

func ParseAddr(in net.Addr) (net.IP, int, ice.NetworkType, bool) {
	switch addr := in.(type) {
	case *net.UDPAddr:
		return addr.IP, addr.Port, ice.NetworkTypeUDP4, true
	case *net.TCPAddr:
		return addr.IP, addr.Port, ice.NetworkTypeTCP4, true
	}
	return nil, 0, 0, false
}
