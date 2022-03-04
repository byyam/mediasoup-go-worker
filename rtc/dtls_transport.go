package rtc

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"strings"

	"github.com/byyam/mediasoup-go-worker/utils"
	"github.com/pion/dtls/v2/pkg/crypto/fingerprint"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"

	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
)

type dtlsTransport struct {
	state        mediasoupdata.DtlsState
	fingerPrints []mediasoupdata.DtlsFingerprint
	role         mediasoupdata.DtlsRole
	tlsCerts     []tls.Certificate
	logger       utils.Logger
}

type dtlsTransportParam struct {
	isClient bool
}

func newDtlsTransport(param dtlsTransportParam) (*dtlsTransport, error) {

	d := &dtlsTransport{
		state:  mediasoupdata.DtlsState_New,
		logger: utils.NewLogger("dtls"),
	}
	switch param.isClient {
	case true:
		d.role = "client"
	case false:
		d.role = "server"
	}
	if err := d.selfSignCerts(); err != nil {
		return nil, err
	}
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
		Value:     strings.ToUpper(actualSHA256),
	})
	return nil
}
