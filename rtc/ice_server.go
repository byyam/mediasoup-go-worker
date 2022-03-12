package rtc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/pion/transport/packetio"

	"github.com/byyam/mediasoup-go-worker/conf"

	"github.com/pion/stun"

	"github.com/byyam/mediasoup-go-worker/global"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/utils"

	"github.com/pion/ice/v2"
)

const (
	maxBufferSize = 1000 * 1000 // 1MB
)

type iceServer struct {
	iceLite    bool
	state      mediasoupdata.IceState
	localUfrag string
	localPwd   string
	logger     utils.Logger
	udpMux     *ice.UDPMuxDefault
	udpConn    net.PacketConn
	buffer     *packetio.Buffer
	// remote info
	iceConn  *iceConn
	connDone chan struct{}
	// handler
	onPacketReceivedHandler atomic.Value
}

type iceServerParam struct {
	iceLite          bool
	tcp4             bool
	OnPacketReceived func(data []byte, len int)
}

func newIceServer(param iceServerParam) (*iceServer, error) {
	ufrag, _ := utils.GenerateUFrag()
	pwd, _ := utils.GeneratePwd()
	i := &iceServer{
		iceLite:    param.iceLite, // todo: support full ICE
		state:      mediasoupdata.IceState_New,
		logger:     utils.NewLogger("ice"),
		localUfrag: ufrag,
		localPwd:   pwd,
		udpMux:     global.UdpMuxConn,
		connDone:   make(chan struct{}),
		buffer:     packetio.NewBuffer(),
	}
	i.onPacketReceivedHandler.Store(param.OnPacketReceived)
	i.buffer.SetLimitSize(maxBufferSize)
	networkTypes := []ice.NetworkType{ice.NetworkTypeUDP4} // udp is default
	if param.tcp4 {
		networkTypes = append(networkTypes, ice.NetworkTypeTCP4)
	}
	i.logger.Debug("ice server start, ufrag=%s", i.localUfrag)

	go func() {
		if err := i.connect(networkTypes); err != nil {
			i.logger.Error("ice connecting failed:%v", err)
			return
		}
	}()

	return i, nil
}

func (i *iceServer) connect(networkTypes []ice.NetworkType) error {
	var err error
	i.udpConn, err = i.udpMux.GetConn(i.localUfrag, false)
	if err != nil {
		return err
	}
	i.logger.Debug("get pkg connection from udp mux:%s", i.localUfrag)
	buf := make([]byte, global.ReceiveMTU)
	for {
		n, srcAddr, err := i.udpConn.ReadFrom(buf)
		if err != nil {
			return err
		}
		i.logger.Debug("read mux n=%d, addr=%s, err=%v", n, srcAddr.String(), err)
		if err := i.handleInboundMsg(buf[:n], n, srcAddr); err != nil {
			return err
		}
	}
}

func (i *iceServer) handleInboundMsg(buffer []byte, n int, srcAddr net.Addr) error {
	if stun.IsMessage(buffer) {
		m := &stun.Message{
			Raw: make([]byte, len(buffer)),
		}
		// Explicitly copy raw buffer so Message can own the memory.
		copy(m.Raw, buffer)
		if err := m.Decode(); err != nil {
			i.logger.Error("Failed to handle decode ICE from %s to %s: %v", i.udpMux.LocalAddr(), srcAddr, err)
			return err
		}
		if err := i.handleInbound(m, srcAddr); err != nil {
			i.logger.Error("Failed to handleInbound: %v", err)
			return err
		}
		return nil
	}
	if utils.MatchDTLS(buffer) {
		if _, err := i.buffer.Write(buffer); err != nil {
			i.logger.Warn("Failed to write buffer: %v", err)
		}
		return nil
	}
	if handler, ok := i.onPacketReceivedHandler.Load().(func(data []byte, len int)); ok && handler != nil {
		handler(buffer[:n], n)
	}
	return nil
}

func (i *iceServer) handleInbound(m *stun.Message, remote net.Addr) error {
	var err error
	if m == nil {
		return errors.New("m stun nil")
	}
	if m.Type.Method != stun.MethodBinding ||
		!(m.Type.Class == stun.ClassSuccessResponse ||
			m.Type.Class == stun.ClassRequest ||
			m.Type.Class == stun.ClassIndication) {
		return fmt.Errorf("unhandled STUN from %s to %s class(%s) method(%s)", remote, i.udpMux.LocalAddr(), m.Type.Class, m.Type.Method)
	}
	if m.Contains(stun.AttrICEControlled) {
		return fmt.Errorf("inbound isControlled && a.isControlling == false")
	}

	if m.Type.Class == stun.ClassRequest {
		if err = utils.AssertInboundUsername(m, i.localUfrag+":"+""); err != nil {
			return fmt.Errorf("discard message from (%s), %v", remote, err)
		} else if err = utils.AssertInboundMessageIntegrity(m, []byte(i.localPwd)); err != nil {
			return fmt.Errorf("discard message from (%s), %v", remote, err)
		}
		log.Printf("inbound STUN (Request) from %s to %s", remote.String(), i.udpMux.LocalAddr())
		if err := i.handleBindingRequest(m, remote); err != nil {
			return err
		}
	}
	return nil
}

func (i *iceServer) handleBindingRequest(m *stun.Message, remote net.Addr) error {
	if m.Contains(stun.AttrUseCandidate) {
		// todo
		log.Printf("get AttrUseCandidate")
	}
	return i.sendBindingSuccess(m, remote)
}

func (i *iceServer) sendBindingSuccess(m *stun.Message, remote net.Addr) error {
	ip, port, _, ok := utils.ParseAddr(i.udpMux.LocalAddr())
	if !ok {
		return fmt.Errorf("error parsing addr: %s", i.udpMux.LocalAddr())
	}
	if out, err := stun.Build(m, stun.BindingSuccess,
		&stun.XORMappedAddress{
			IP:   ip,
			Port: port,
		},
		stun.NewShortTermIntegrity(i.localPwd),
		stun.Fingerprint,
	); err != nil {
		return fmt.Errorf("failed to handle inbound ICE from: %s to: %s error: %s", i.udpMux.LocalAddr(), remote, err)
	} else {
		if i.iceConn == nil { // todo
			i.iceConn = newIceConn(remote, i)
			i.logger.Debug("new ice connection,remote addr=%s", remote.String())
			close(i.connDone)
		}
		_, err = i.iceConn.Write(out.Raw)
		if err != nil {
			return fmt.Errorf("failed to send STUN message: %s", err)
		}
	}
	return nil
}

func (i *iceServer) GetIceParameters() mediasoupdata.IceParameters {
	return mediasoupdata.IceParameters{
		UsernameFragment: i.localUfrag,
		Password:         i.localPwd,
		IceLite:          i.iceLite,
	}
}

func (i *iceServer) GetSelectedTuple() mediasoupdata.TransportTuple {
	return mediasoupdata.TransportTuple{}
}

func (i *iceServer) GetState() mediasoupdata.IceState {
	return i.state
}

func (i *iceServer) GetRole() string {
	return "controlled"
}

func (i *iceServer) GetLocalCandidates() (iceCandidates []mediasoupdata.IceCandidate) {
	candidate := mediasoupdata.IceCandidate{
		Foundation: "udpcandidate",
		Priority:   0,
		Ip:         conf.Settings.RtcListenIp,
		Protocol:   "udp",
		Port:       uint32(conf.Settings.RtcStaticPort),
		Type:       "host",
		TcpType:    "",
	}
	iceCandidates = append(iceCandidates, candidate)

	return
}

func (i *iceServer) GetConn() (*iceConn, error) {
	if i.connDone != nil {
		<-i.connDone
		i.logger.Debug("ice connected")
	}
	return i.iceConn, nil
}
