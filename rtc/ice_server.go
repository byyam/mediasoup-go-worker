package rtc

import (
	"github.com/pion/logging"

	"github.com/byyam/mediasoup-go-worker/global"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/utils"
	"github.com/pion/ice/v2"
)

type iceServer struct {
	iceLite    bool
	agent      *ice.Agent
	udpMux     *ice.UDPMuxDefault
	state      mediasoupdata.IceState
	localUfrag string
	localPwd   string
	logger     utils.Logger
}

type iceServerParam struct {
	iceLite bool
	udp4    bool
	tcp4    bool
}

func newIceServer(param iceServerParam) (*iceServer, error) {
	i := &iceServer{
		iceLite: true, // todo
		state:   mediasoupdata.IceState_New,
		logger:  utils.NewLogger("ice"),
	}
	var networkTypes []ice.NetworkType
	if param.udp4 {
		networkTypes = append(networkTypes, ice.NetworkTypeUDP4)
	}
	if param.tcp4 {
		networkTypes = append(networkTypes, ice.NetworkTypeTCP4)
	}

	if err := i.connect(networkTypes); err != nil {
		return nil, err
	}

	return i, nil
}

func (i *iceServer) connect(networkTypes []ice.NetworkType) error {
	var err error

	loggerFactory := logging.NewDefaultLoggerFactory()
	i.udpMux = ice.NewUDPMuxDefault(ice.UDPMuxParams{
		Logger:  loggerFactory.NewLogger("ice"),
		UDPConn: global.UdpMuxConn,
	})
	i.logger.Debug("udp mux addr:%+v", i.udpMux.LocalAddr().String())

	if i.agent, err = ice.NewAgent(&ice.AgentConfig{
		Lite:           i.iceLite,
		NetworkTypes:   networkTypes,
		UDPMux:         i.udpMux,
		CandidateTypes: []ice.CandidateType{ice.CandidateTypeHost},
	}); err != nil {
		return err
	}
	if i.localUfrag, i.localPwd, err = i.agent.GetLocalUserCredentials(); err != nil {
		return err
	}
	if err = i.agent.OnCandidate(func(candidate ice.Candidate) {
		i.logger.Debug("OnCandidate:%+v", candidate)
	}); err != nil {
		return err
	}
	if err = i.agent.OnConnectionStateChange(func(state ice.ConnectionState) {
		i.logger.Debug("OnConnectionStateChange:%v", state)
	}); err != nil {
		return err
	}
	if err = i.agent.OnSelectedCandidatePairChange(func(candidate ice.Candidate, candidate2 ice.Candidate) {
		i.logger.Debug("OnSelectedCandidatePairChange,candidate=%+v,candidate2=%+v", candidate, candidate2)
	}); err != nil {
		return err
	}
	if err = i.agent.GatherCandidates(); err != nil {
		return err
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
	candidates, err := i.agent.GetLocalCandidates()
	if err != nil {
		i.logger.Error("GetLocalCandidates failed:%+v", err)
		return
	}
	for _, c := range candidates {
		can := mediasoupdata.IceCandidate{
			Foundation: c.Foundation(),
			Priority:   c.Priority(),
			Ip:         c.Address(),
			Protocol:   mediasoupdata.TransportProtocol(c.NetworkType().NetworkShort()),
			Port:       uint32(c.Port()),
			Type:       c.Type().String(),
			TcpType:    c.TCPType().String(),
		}
		iceCandidates = append(iceCandidates, can)
	}
	i.logger.Debug("GetLocalCandidates:%+v", candidates)
	return
}
