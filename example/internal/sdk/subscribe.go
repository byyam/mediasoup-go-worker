package sdk

import (
	"log"
	"net"

	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
	"github.com/byyam/mediasoup-go-worker/example/internal/wsconn"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/pkg/errors"
)

type SubscribeOpt struct {
	StreamId   uint64
	DtlsClient bool // default: server
}

func (c *Client) Subscribe(opt SubscribeOpt) (net.Conn, *isignal.SubscribeResponse, error) {
	mediasoupFPs, err := c.prepareDtls(opt.DtlsClient)
	if err != nil {
		return nil, nil, err
	}
	rtpCapabilities := mediasoupdata.GetSupportedRtpCapabilities()

	webRtcTransportOffer := isignal.WebRtcTransportOffer{
		ForceTcp: false,
		DtlsParameters: mediasoupdata.DtlsParameters{
			Role:         mediasoupdata.DtlsRole(c.dtlsRole),
			Fingerprints: mediasoupFPs,
		},
	}
	req := isignal.SubscribeRequest{
		StreamId:    opt.StreamId,
		TransportId: "",
		Offer:       webRtcTransportOffer,
		SubscribeOffer: isignal.SubscribeOffer{
			Kind:            mediasoupdata.MediaKind_Video,
			AppData:         nil,
			RtpCapabilities: &rtpCapabilities,
		},
	}
	rsp, err := wsconn.NewWsClient(c.wsOpt).Subscribe(req)
	if err != nil {
		log.Println("subscribe response error:", err)
		return nil, nil, err
	}
	c.webrtcTransportAnswer = &rsp.Answer

	conn, err := c.Conn()
	if err != nil {
		return nil, nil, errors.WithMessage(err, "connect sfu error")
	}
	log.Println("handle subscribe completed")
	return conn, &rsp, nil
}
