package rtc

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog"

	FBS__Transport "github.com/byyam/mediasoup-go-worker/fbs/FBS/Transport"
	FBS__WebRtcTransport "github.com/byyam/mediasoup-go-worker/fbs/FBS/WebRtcTransport"
	"github.com/byyam/mediasoup-go-worker/pkg/iceutil"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/muxpkg"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/pion/ice/v3"
	"github.com/pion/stun/v2"
	"github.com/pion/transport/v3/packetio"

	"github.com/byyam/mediasoup-go-worker/conf"
	"github.com/byyam/mediasoup-go-worker/internal/global"
	"github.com/byyam/mediasoup-go-worker/monitor"
)

const (
	maxBufferSize = 1000 * 1000 // 1MB
	// keepaliveInterval used to keep candidates alive
	defaultKeepaliveInterval = 5 * time.Second
	// defaultDisconnectedTimeout is the default time till an Agent transitions disconnected
	defaultDisconnectedTimeout = 30 * time.Second
)

type iceServer struct {
	iceLite    bool
	state      FBS__WebRtcTransport.IceState
	localUfrag string
	localPwd   string
	logger     zerolog.Logger
	udpMux     *ice.UDPMuxDefault
	udpConn    net.PacketConn
	buffer     *packetio.Buffer
	// timestamp
	lastStunTimestamp   time.Time
	lastPkgTimestamp    time.Time
	disconnectedTimeout time.Duration
	keepaliveInterval   time.Duration
	// remote info
	iceConn  *iceConn
	connDone chan struct{}
	// close
	closedChan chan struct{}
	closeOnce  sync.Once
	// handler
	onPacketReceivedHandler func(data []byte)
}

type iceServerParam struct {
	transportId         string
	iceLite             bool
	tcp4                bool
	OnPacketReceived    func(data []byte)
	DisconnectedTimeout *time.Duration
	KeepaliveInterval   *time.Duration
}

func newIceServer(param iceServerParam) (*iceServer, error) {
	ufrag, _ := iceutil.GenerateUFrag()
	pwd, _ := iceutil.GeneratePwd()
	d := &iceServer{
		iceLite:          param.iceLite,                    // todo: support full ICE
		state:            FBS__WebRtcTransport.IceStateNEW, // todo: completed
		logger:           zerowrapper.NewScope(string(mediasoupdata.WorkerLogTag_ICE), param.transportId),
		localUfrag:       ufrag,
		localPwd:         pwd,
		udpMux:           global.ICEMuxConn,
		connDone:         make(chan struct{}),
		closedChan:       make(chan struct{}),
		buffer:           packetio.NewBuffer(),
		lastPkgTimestamp: time.Now(), // init stun TS to now.
	}
	if param.DisconnectedTimeout == nil {
		d.disconnectedTimeout = defaultDisconnectedTimeout
	} else {
		d.disconnectedTimeout = *param.DisconnectedTimeout
	}
	if param.KeepaliveInterval == nil {
		d.keepaliveInterval = defaultKeepaliveInterval
	} else {
		d.keepaliveInterval = *param.KeepaliveInterval
	}
	d.onPacketReceivedHandler = param.OnPacketReceived
	d.buffer.SetLimitSize(maxBufferSize)
	networkTypes := []ice.NetworkType{ice.NetworkTypeUDP4} // udp is default
	if param.tcp4 {
		networkTypes = append(networkTypes, ice.NetworkTypeTCP4)
	}
	d.logger.Debug().Msgf("ice server start, ufrag=%s", d.localUfrag)

	go func() {
		if err := d.connect(networkTypes); err != nil {
			d.logger.Error().Err(err).Msg("read ice connection failed")
			return
		}
	}()

	go d.connectivityChecks()

	return d, nil
}

func (d *iceServer) connectivityChecks() {
	checkFn := func() {
		if !d.isConnected() {
			return
		}
		if time.Since(d.lastPkgTimestamp) > d.disconnectedTimeout && time.Since(d.lastStunTimestamp) > d.disconnectedTimeout {
			d.logger.Warn().Dur("disconnectedTimeout", d.disconnectedTimeout).Str("localUfrag", d.localUfrag).Msg("ice inactive")
			d.Disconnect()
		}
	}
	for {
		t := time.NewTimer(defaultKeepaliveInterval)
		select {
		case <-t.C:
			checkFn()
		case <-d.closedChan:
			d.logger.Warn().Msg("stop ice connectivityChecks")
			return
		}
	}
}

func (d *iceServer) connect(networkTypes []ice.NetworkType) error {
	var err error
	d.udpConn, err = d.udpMux.GetConn(d.localUfrag, global.ICEMuxAddr)
	if err != nil {
		return err
	}
	d.logger.Debug().Str("localUfrag", d.localUfrag).Msg("get pkg connection from udp mux")
	buf := make([]byte, conf.Settings.ReceiveMTU)
	for {
		n, srcAddr, err := d.udpConn.ReadFrom(buf)
		if err != nil {
			return err
		}
		d.logger.Debug().Msgf("read mux n=%d, addr=%s, err=%v", n, srcAddr.String(), err)
		if err := d.handleInboundMsg(buf[:n], n, srcAddr); err != nil {
			return err
		}
	}
}

func (d *iceServer) handleInboundMsg(buffer []byte, n int, srcAddr net.Addr) error {
	if stun.IsMessage(buffer) {
		monitor.IceCount(monitor.DirectionTypeRecv, monitor.PacketStun)
		m := &stun.Message{
			Raw: make([]byte, len(buffer)),
		}
		// Explicitly copy raw buffer so Message can own the memory.
		copy(m.Raw, buffer)
		if err := m.Decode(); err != nil {
			d.logger.Error().Err(err).Str("LocalAddr", d.udpMux.LocalAddr().String()).
				Str("srcAddr", srcAddr.String()).
				Msg("Failed to handle decode ICE")
			return err
		}
		if err := d.handleInbound(m, srcAddr); err != nil {
			d.logger.Error().Err(err).Msg("Failed to handleInbound")
			return err
		}
		return nil
	}
	if muxpkg.MatchDTLS(buffer) {
		monitor.IceCount(monitor.DirectionTypeRecv, monitor.PacketDtls)
		if _, err := d.buffer.Write(buffer); err != nil {
			d.logger.Warn().Err(err).Msg("Failed to write buffer")
		}
		return nil
	}
	d.onPacketReceivedHandler(buffer[:n])
	d.lastPkgTimestamp = time.Now()
	return nil
}

func (d *iceServer) handleInbound(m *stun.Message, remote net.Addr) error {
	var err error
	if m == nil {
		return errors.New("m stun nil")
	}
	if m.Type.Method != stun.MethodBinding ||
		!(m.Type.Class == stun.ClassSuccessResponse ||
			m.Type.Class == stun.ClassRequest ||
			m.Type.Class == stun.ClassIndication) {
		return fmt.Errorf("unhandled STUN from %s to %s class(%s) method(%s)", remote, d.udpMux.LocalAddr(), m.Type.Class, m.Type.Method)
	}
	if m.Contains(stun.AttrICEControlled) {
		return fmt.Errorf("inbound isControlled && a.isControlling == false")
	}

	if m.Type.Class == stun.ClassRequest {
		if err = iceutil.AssertInboundUsername(m, d.localUfrag+":"+""); err != nil {
			return fmt.Errorf("discard message from (%s), %v", remote, err)
		} else if err = iceutil.AssertInboundMessageIntegrity(m, []byte(d.localPwd)); err != nil {
			return fmt.Errorf("discard message from (%s), %v", remote, err)
		}
		d.logger.Debug().Str("remote", remote.String()).Str("local", d.udpMux.LocalAddr().String()).Msg("inbound STUN (Request)")
		if err := d.handleBindingRequest(m, remote); err != nil {
			return err
		}
		d.lastStunTimestamp = time.Now()
	}
	return nil
}

func (d *iceServer) handleBindingRequest(m *stun.Message, remote net.Addr) error {
	if m.Contains(stun.AttrUseCandidate) {
		// todo
		d.logger.Info().Msg("get AttrUseCandidate")
	}
	return d.sendBindingSuccess(m, remote)
}

func (d *iceServer) sendBindingSuccess(m *stun.Message, remote net.Addr) error {
	ip, port, _, ok := iceutil.ParseAddr(d.udpMux.LocalAddr())
	if !ok {
		return fmt.Errorf("error parsing addr: %s", d.udpMux.LocalAddr())
	}
	if out, err := stun.Build(m, stun.BindingSuccess,
		&stun.XORMappedAddress{
			IP:   ip,
			Port: port,
		},
		stun.NewShortTermIntegrity(d.localPwd),
		stun.Fingerprint,
	); err != nil {
		return fmt.Errorf("failed to handle inbound ICE from: %s to: %s error: %s", d.udpMux.LocalAddr(), remote, err)
	} else {
		if d.iceConn == nil { // todo
			d.iceConn = newIceConn(remote, d)
			d.logger.Debug().Str("remote", remote.String()).Msg("new ice connection")
			close(d.connDone)
		}
		_, err = d.iceConn.Write(out.Raw)
		if err != nil {
			return fmt.Errorf("failed to send STUN message: %s", err)
		}
	}
	return nil
}

func (d *iceServer) GetIceParameters() *FBS__WebRtcTransport.IceParametersT {
	return &FBS__WebRtcTransport.IceParametersT{
		UsernameFragment: d.localUfrag,
		Password:         d.localPwd,
		IceLite:          d.iceLite,
	}
}

func (d *iceServer) GetSelectedTuple() *FBS__Transport.TupleT {
	tuple := &FBS__Transport.TupleT{}
	// ice conn may be nil
	if d.iceConn == nil {
		return tuple
	}
	localAddr := d.iceConn.LocalAddr()
	localUdpAddr, ok := localAddr.(*net.UDPAddr)
	if ok {
		tuple.LocalIp = localUdpAddr.IP.String()
		tuple.LocalPort = uint16(localUdpAddr.Port)
	}
	remoteAddr := d.iceConn.RemoteAddr()
	remoteUdpAddr, ok := remoteAddr.(*net.UDPAddr)
	if ok {
		tuple.RemoteIp = remoteUdpAddr.IP.String()
		tuple.RemotePort = uint16(remoteUdpAddr.Port)
	}
	tuple.Protocol = FBS__Transport.ProtocolUDP
	return tuple
}

func (d *iceServer) GetState() FBS__WebRtcTransport.IceState {
	return d.state
}

func (d *iceServer) GetRole() FBS__WebRtcTransport.IceRole {
	return FBS__WebRtcTransport.IceRoleCONTROLLED
}

func (d *iceServer) GetLocalCandidates() (iceCandidates []*FBS__WebRtcTransport.IceCandidateT) {
	candidate := &FBS__WebRtcTransport.IceCandidateT{
		Foundation: "udpcandidate",
		Priority:   0,
		Ip:         conf.Settings.RtcListenIp,
		Protocol:   FBS__Transport.ProtocolUDP,
		Port:       global.ICEMuxPort,
		Type:       FBS__WebRtcTransport.IceCandidateTypeHOST,
		TcpType:    nil,
	}
	iceCandidates = append(iceCandidates, candidate)

	return
}

func (d *iceServer) GetConn() (*iceConn, error) {
	if d.connDone != nil {
		<-d.connDone
		d.state = FBS__WebRtcTransport.IceStateCONNECTED
		d.logger.Info().Msg("ice connected")
	}
	return d.iceConn, nil
}

func (d *iceServer) Disconnect() {
	d.closeOnce.Do(func() {
		close(d.closedChan)
	})
	if d.iceConn != nil {
		if err := d.iceConn.Close(); err != nil {
			d.logger.Error().Err(err).Msg("disconnect ice failed")
		}
	}
	if d.buffer != nil {
		if err := d.buffer.Close(); err != nil {
			d.logger.Error().Err(err).Msg("close ice buffer failed")
		}
	}
	if d.udpConn != nil {
		if err := d.udpConn.Close(); err != nil {
			d.logger.Error().Err(err).Msg("close udp conn failed")
		}
	}
	d.udpMux.RemoveConnByUfrag(d.localUfrag)
	d.state = FBS__WebRtcTransport.IceStateDISCONNECTED
	d.logger.Info().Msg("ice disconnect")
}

func (d *iceServer) isConnected() bool {
	select {
	case <-d.connDone:
		return true
	default:
		return false
	}
}

func (d *iceServer) isClosed() bool {
	select {
	case <-d.closedChan:
		return true
	default:
		return false
	}
}

func (d *iceServer) CloseChannel() <-chan struct{} {
	return d.closedChan
}
