package rtc

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog"

	mediasoupdata2 "github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"

	"github.com/pion/dtls/v2"
	"github.com/pion/dtls/v2/pkg/crypto/fingerprint"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"
	"github.com/pion/srtp/v2"

	"github.com/byyam/mediasoup-go-worker/mserror"
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
	state                  mediasoupdata2.DtlsState
	fingerPrints           []mediasoupdata2.DtlsFingerprint
	role                   mediasoupdata2.DtlsRole
	tlsCerts               []tls.Certificate
	logger                 zerolog.Logger
	connTimeout            time.Duration
	fingerprintAlgorithms  []crypto.Hash
	sRTPProtectionProfiles []dtls.SRTPProtectionProfile
	srtpProtectionProfile  srtp.ProtectionProfile
}

type dtlsTransportParam struct {
	transportId string
	role        mediasoupdata2.DtlsRole
	connTimeout *time.Duration
}

func newDtlsTransport(param dtlsTransportParam) (*dtlsTransport, error) {
	d := &dtlsTransport{
		state:                  mediasoupdata2.DtlsState_New,
		role:                   param.role,
		logger:                 zerowrapper.NewScope(string(mediasoupdata2.WorkerLogTag_DTLS), param.transportId),
		fingerprintAlgorithms:  defaultFingerprintAlgorithms,
		sRTPProtectionProfiles: defaultSRTPProtectionProfiles,
	}
	d.fingerPrints = make([]mediasoupdata2.DtlsFingerprint, len(d.fingerprintAlgorithms))
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

func (d *dtlsTransport) GetDtlsParameters() mediasoupdata2.DtlsParameters {
	return mediasoupdata2.DtlsParameters{
		Role:         d.role,
		Fingerprints: d.fingerPrints,
	}
}

func (d *dtlsTransport) GetState() mediasoupdata2.DtlsState {
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
	d.logger.Debug().Msgf("x509 length:%d", len(x509cert.Raw))
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
		d.fingerPrints[i] = mediasoupdata2.DtlsFingerprint{
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

func (d *dtlsTransport) SetRole(remoteParam *mediasoupdata2.DtlsParameters) (*mediasoupdata2.TransportConnectData, error) {

	switch remoteParam.Role {
	case mediasoupdata2.DtlsRole_Client, mediasoupdata2.DtlsRole_Auto:
		d.role = mediasoupdata2.DtlsRole_Server
	case mediasoupdata2.DtlsRole_Server:
		d.role = mediasoupdata2.DtlsRole_Client
	default:
		return nil, mserror.ErrInvalidParam
	}
	return &mediasoupdata2.TransportConnectData{DtlsLocalRole: d.role}, nil
}

func (d *dtlsTransport) Connect(iceConn net.Conn) error {
	d.state = mediasoupdata2.DtlsState_Connecting
	var err error
	defer func() {
		if err != nil {
			d.state = mediasoupdata2.DtlsState_Failed
			d.logger.Error().Msgf("dtls connecting failed:%v", err)
		}
	}()
	d.logger.Debug().Msgf("dtlsRole=%s,iceConn=%s|%s", d.role, iceConn.LocalAddr(), iceConn.RemoteAddr())
	if d.role == mediasoupdata2.DtlsRole_Client {
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
	d.state = mediasoupdata2.DtlsState_Connected
	d.logger.Info().Msg("DtlsState_Connected")
	return nil
}

func (d *dtlsTransport) GetSRTPConfig() (*srtp.Config, error) {
	srtpProfile, ok := d.dtlsConn.SelectedSRTPProtectionProfile()
	if !ok {
		return nil, mserror.ErrNoSRTPProtectionProfile
	}
	switch srtpProfile {
	case dtls.SRTP_AEAD_AES_128_GCM:
		d.srtpProtectionProfile = srtp.ProtectionProfileAeadAes128Gcm
	case dtls.SRTP_AES128_CM_HMAC_SHA1_80:
		d.srtpProtectionProfile = srtp.ProtectionProfileAes128CmHmacSha1_80
	default:
		return nil, mserror.ErrNoSRTPProtectionProfile
	}
	srtpConfig := &srtp.Config{
		Profile: d.srtpProtectionProfile,
	}
	var isClient bool
	if d.role == mediasoupdata2.DtlsRole_Client {
		isClient = true
	}
	if err := srtpConfig.ExtractSessionKeysFromDTLS(&d.dtlsConnState, isClient); err != nil {
		return nil, fmt.Errorf("errDtlsKeyExtractionFailed: %v", err)
	}
	return srtpConfig, nil
}

func (d *dtlsTransport) Disconnect() {
	d.logger.Info().Msg("dtls disconnect")
}
