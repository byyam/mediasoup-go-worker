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

type dtlsTransport struct {
	dtlsConn      *dtls.Conn
	dtlsConnState dtls.State
	config        *dtls.Config
	state         mediasoupdata.DtlsState
	fingerPrints  []mediasoupdata.DtlsFingerprint
	role          mediasoupdata.DtlsRole
	tlsCerts      []tls.Certificate
	logger        utils.Logger
}

type dtlsTransportParam struct {
	role        mediasoupdata.DtlsRole
	connTimeout time.Duration
}

func newDtlsTransport(param dtlsTransportParam) (*dtlsTransport, error) {
	if param.connTimeout <= 0 {
		return nil, common.ErrInvalidParam
	}
	d := &dtlsTransport{
		state:  mediasoupdata.DtlsState_New,
		role:   param.role,
		logger: utils.NewLogger("dtls"),
	}

	if err := d.selfSignCerts(); err != nil {
		return nil, err
	}
	d.prepareConfig(param.connTimeout)
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
	d.logger.Debug("x509 length:%d", len(x509cert.Raw))
	actualSHA256, err := fingerprint.Fingerprint(x509cert, crypto.SHA256)
	if err != nil {
		return err
	}
	d.fingerPrints = append(d.fingerPrints, mediasoupdata.DtlsFingerprint{
		Algorithm: "sha-256",
		Value:     actualSHA256,
	})
	d.tlsCerts = append(d.tlsCerts, certificate)
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
		SRTPProtectionProfiles: []dtls.SRTPProtectionProfile{dtls.SRTP_AES128_CM_HMAC_SHA1_80},
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
	srtpConfig := &srtp.Config{
		Profile: srtp.ProtectionProfileAes128CmHmacSha1_80,
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
