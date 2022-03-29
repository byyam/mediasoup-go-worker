package rtc

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"time"

	"github.com/pion/dtls/v2"
	"github.com/pion/dtls/v2/pkg/crypto/fingerprint"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"
	"github.com/pion/srtp/v2"

	"github.com/byyam/mediasoup-go-worker/common"
	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
)

const (
	defaultDtlsConnectTimeout = 30 * time.Second
)

var (
	defaultSRTPProtectionProfiles = []dtls.SRTPProtectionProfile{
		dtls.SRTP_AEAD_AES_128_GCM,
		dtls.SRTP_AES128_CM_HMAC_SHA1_80,
	}

	defaultFingerprintAlgorithms = []crypto.Hash{
		crypto.SHA1,
		crypto.SHA224,
		crypto.SHA256,
		crypto.SHA384,
		crypto.SHA512,
	}
)

type dtlsTransport struct {
	dtlsConn               *dtls.Conn
	dtlsConnState          dtls.State
	config                 *dtls.Config
	state                  mediasoupdata.DtlsState
	fingerPrints           []mediasoupdata.DtlsFingerprint
	role                   mediasoupdata.DtlsRole
	tlsCerts               []tls.Certificate
	logger                 utils.Logger
	connTimeout            time.Duration
	fingerprintAlgorithms  []crypto.Hash
	sRTPProtectionProfiles []dtls.SRTPProtectionProfile
	srtpProtectionProfile  srtp.ProtectionProfile
}

type dtlsTransportParam struct {
	transportId string
	role        mediasoupdata.DtlsRole
	connTimeout *time.Duration
}

func newDtlsTransport(param dtlsTransportParam) (*dtlsTransport, error) {
	d := &dtlsTransport{
		state:                  mediasoupdata.DtlsState_New,
		role:                   param.role,
		logger:                 utils.NewLogger("dtls", param.transportId),
		fingerprintAlgorithms:  defaultFingerprintAlgorithms,
		sRTPProtectionProfiles: defaultSRTPProtectionProfiles,
	}
	d.fingerPrints = make([]mediasoupdata.DtlsFingerprint, len(d.fingerprintAlgorithms))
	if param.connTimeout == nil {
		d.connTimeout = defaultDtlsConnectTimeout
	} else {
		d.connTimeout = *param.connTimeout
	}
	if err := d.selfSignCerts(); err != nil {
		return nil, err
	}
	d.prepareConfig(d.connTimeout)
	return d, nil
}

func (d *dtlsTransport) GetDtlsParameters() mediasoupdata.DtlsParameters {
	return mediasoupdata.DtlsParameters{
		Role:         d.role,
		Fingerprints: d.fingerPrints,
	}
}

func (d *dtlsTransport) GetState() mediasoupdata.DtlsState {
	return d.state
}

func (d *dtlsTransport) selfSignCerts() error {
	certificate, err := selfsign.GenerateSelfSigned()
	if err != nil {
		return err
	}
	x509cert, err := x509.ParseCertificate(certificate.Certificate[0])
	if err != nil {
		return err
	}
	d.tlsCerts = append(d.tlsCerts, certificate)
	d.logger.Debug("x509 length:%d", len(x509cert.Raw))
	// set fingerprint
	for i, algo := range d.fingerprintAlgorithms {
		name, err := fingerprint.StringFromHash(algo)
		if err != nil {
			return err
		}
		value, err := fingerprint.Fingerprint(x509cert, algo)
		if err != nil {
			return err
		}
		d.fingerPrints[i] = mediasoupdata.DtlsFingerprint{
			Algorithm: name,
			Value:     value,
		}
	}
	return nil
}

func (d *dtlsTransport) prepareConfig(timeout time.Duration) {
	d.config = &dtls.Config{
		Certificates:         d.tlsCerts,
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
		// Create timeout context for accepted connection.
		ConnectContextMaker: func() (context.Context, func()) {
			return context.WithTimeout(context.Background(), timeout)
		},
		SRTPProtectionProfiles: d.sRTPProtectionProfiles,
	}
}

func (d *dtlsTransport) SetRole(remoteParam *mediasoupdata.DtlsParameters) (*mediasoupdata.TransportConnectData, error) {

	switch remoteParam.Role {
	case mediasoupdata.DtlsRole_Client, mediasoupdata.DtlsRole_Auto:
		d.role = mediasoupdata.DtlsRole_Server
	case mediasoupdata.DtlsRole_Server:
		d.role = mediasoupdata.DtlsRole_Client
	default:
		return nil, common.ErrInvalidParam
	}
	return &mediasoupdata.TransportConnectData{DtlsLocalRole: d.role}, nil
}

func (d *dtlsTransport) Connect(iceConn net.Conn) error {
	d.state = mediasoupdata.DtlsState_Connecting
	var err error
	defer func() {
		if err != nil {
			d.state = mediasoupdata.DtlsState_Failed
			d.logger.Error("dtls connecting failed:%v", err)
		}
	}()
	d.logger.Debug("dtlsRole=%s,iceConn=%s|%s", d.role, iceConn.LocalAddr(), iceConn.RemoteAddr())
	if d.role == mediasoupdata.DtlsRole_Client {
		d.config.InsecureSkipVerify = true
		if d.dtlsConn, err = dtls.Client(iceConn, d.config); err != nil {
			return err
		}
	} else {
		if d.dtlsConn, err = dtls.Server(iceConn, d.config); err != nil {
			return err
		}
	}
	d.dtlsConnState = d.dtlsConn.ConnectionState()
	d.state = mediasoupdata.DtlsState_Connected
	d.logger.Info("DtlsState_Connected")
	return nil
}

func (d *dtlsTransport) GetSRTPConfig() (*srtp.Config, error) {
	srtpProfile, ok := d.dtlsConn.SelectedSRTPProtectionProfile()
	if !ok {
		return nil, common.ErrNoSRTPProtectionProfile
	}
	switch srtpProfile {
	case dtls.SRTP_AEAD_AES_128_GCM:
		d.srtpProtectionProfile = srtp.ProtectionProfileAeadAes128Gcm
	case dtls.SRTP_AES128_CM_HMAC_SHA1_80:
		d.srtpProtectionProfile = srtp.ProtectionProfileAes128CmHmacSha1_80
	default:
		return nil, common.ErrNoSRTPProtectionProfile
	}
	srtpConfig := &srtp.Config{
		Profile: d.srtpProtectionProfile,
	}
	var isClient bool
	if d.role == mediasoupdata.DtlsRole_Client {
		isClient = true
	}
	if err := srtpConfig.ExtractSessionKeysFromDTLS(&d.dtlsConnState, isClient); err != nil {
		return nil, fmt.Errorf("errDtlsKeyExtractionFailed: %v", err)
	}
	return srtpConfig, nil
}

func (d *dtlsTransport) Disconnect() {
	d.logger.Info("dtls disconnect")
}
