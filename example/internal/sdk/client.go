package sdk

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
	"github.com/byyam/mediasoup-go-worker/example/internal/wsconn"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/pion/dtls/v2"
	"github.com/pion/dtls/v2/pkg/crypto/fingerprint"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"
	"github.com/pion/ice/v2"
	"github.com/pion/srtp/v2"
)

var (
	iceKeepAlive = 5 * time.Second
)

type Client struct {
	iceConn               *ice.Conn
	dtlsConnState         dtls.State
	dtlsClient            bool
	dtlsRole              string
	webrtcTransportAnswer *isignal.WebRtcTransportAnswer
	tlsCerts              []tls.Certificate
	wsOpt                 wsconn.WsClientOpt
}

func NewClient(opt wsconn.WsClientOpt) *Client {
	return &Client{wsOpt: opt}
}

func (c *Client) GetSRTPConfig() (*srtp.Config, error) {
	srtpConfig := &srtp.Config{
		Profile: srtp.ProtectionProfileAes128CmHmacSha1_80,
	}
	if err := srtpConfig.ExtractSessionKeysFromDTLS(&c.dtlsConnState, c.dtlsClient); err != nil {
		return nil, fmt.Errorf("errDtlsKeyExtractionFailed: %v", err)
	}
	return srtpConfig, nil
}

func (c *Client) Conn() (net.Conn, error) {
	if len(c.webrtcTransportAnswer.IceCandidates) == 0 {
		return nil, errors.New("ice candidate list empty")
	}
	iceAgent, err := ice.NewAgent(&ice.AgentConfig{
		NetworkTypes:      []ice.NetworkType{ice.NetworkTypeUDP4},
		KeepaliveInterval: &iceKeepAlive,
	})
	if err != nil {
		return nil, err
	}
	candidate, err := c.convertCandidates(c.webrtcTransportAnswer.IceCandidates[0])
	if err != nil {
		log.Println("ice candidate failed:", err)
		return nil, err
	}
	if err := iceAgent.AddRemoteCandidate(candidate); err != nil {
		log.Println("add remote candidate failed:", err)
		return nil, err
	}
	if err := iceAgent.OnConnectionStateChange(func(state ice.ConnectionState) {
		log.Println("ICE connection state change:", state.String())
	}); err != nil {
		log.Println("ICE connection state change failed:", err)
		return nil, err
	}
	if err := iceAgent.OnSelectedCandidatePairChange(func(candidate ice.Candidate, candidate2 ice.Candidate) {
		log.Println("ICE OnSelectedCandidatePairChange:", candidate)
		log.Println("ICE OnSelectedCandidatePairChange2:", candidate2)
	}); err != nil {
		log.Println("ICE OnSelectedCandidatePairChange failed:", err)
	}
	if err := iceAgent.OnCandidate(func(candidate ice.Candidate) {
		log.Println("ICE OnCandidate", candidate)
	}); err != nil {
		log.Println("ICE OnCandidate failed:", err)
	}
	if err := iceAgent.GatherCandidates(); err != nil {
		log.Println("gather candidate failed:", err)
		return nil, err
	}
	log.Println("start dialing...")

	iceConn, err := iceAgent.Dial(context.Background(), c.webrtcTransportAnswer.IceParameters.UsernameFragment, c.webrtcTransportAnswer.IceParameters.Password)
	if err != nil {
		return nil, err
	}
	c.iceConn = iceConn

	config := &dtls.Config{
		Certificates:         c.tlsCerts,
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
		// Create timeout context for accepted connection.
		ConnectContextMaker: func() (context.Context, func()) {
			return context.WithTimeout(context.Background(), 30*time.Second)
		},
		SRTPProtectionProfiles: []dtls.SRTPProtectionProfile{dtls.SRTP_AES128_CM_HMAC_SHA1_80},
	}

	if c.dtlsClient {
		config.InsecureSkipVerify = true
		dtlsConn, err := dtls.Client(iceConn, config)
		if err != nil {
			return nil, err
		}
		log.Println("dtls client completed...")
		c.dtlsConnState = dtlsConn.ConnectionState()
	} else {
		dtlsConn, err := dtls.Server(iceConn, config)
		if err != nil {
			return nil, err
		}
		log.Println("dtls server completed...")
		c.dtlsConnState = dtlsConn.ConnectionState()
	}
	return iceConn, nil
}

func (c *Client) convertCandidates(candidate mediasoupdata.IceCandidate) (ice.Candidate, error) {
	iceCandidate, err := ice.NewCandidateHost(&ice.CandidateHostConfig{
		CandidateID: "",
		Network:     string(candidate.Protocol),
		Address:     candidate.Ip,
		Port:        int(candidate.Port),
		Component:   0,
		Priority:    candidate.Priority,
		Foundation:  "",
		TCPType:     0,
	})
	return iceCandidate, err
}

func (c *Client) prepareDtls(isClient bool) ([]mediasoupdata.DtlsFingerprint, error) {
	mediasoupFPs, err := c.selfSignCerts()
	if err != nil {
		return nil, err
	}
	log.Println("mediasoupFPs:", mediasoupFPs)

	c.dtlsRole = "server"
	if isClient {
		c.dtlsClient = true
		c.dtlsRole = "client"
	}
	return c.selfSignCerts()
}

func (c *Client) selfSignCerts() ([]mediasoupdata.DtlsFingerprint, error) {
	var mediasoupFPs []mediasoupdata.DtlsFingerprint
	c.tlsCerts = []tls.Certificate{} // init
	certificate, err := selfsign.GenerateSelfSigned()
	if err != nil {
		return mediasoupFPs, err
	}
	x509cert, err := x509.ParseCertificate(certificate.Certificate[0])
	if err != nil {
		return mediasoupFPs, err
	}
	log.Println("x509:", len(x509cert.Raw))
	actualSHA256, err := fingerprint.Fingerprint(x509cert, crypto.SHA256)
	if err != nil {
		return mediasoupFPs, err
	}
	mediasoupFPs = append(mediasoupFPs, mediasoupdata.DtlsFingerprint{
		Algorithm: "sha-256",
		Value:     strings.ToUpper(actualSHA256),
	})
	c.tlsCerts = append(c.tlsCerts, certificate)
	return mediasoupFPs, nil
}
