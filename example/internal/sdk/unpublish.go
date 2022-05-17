package sdk

import (
	"github.com/byyam/mediasoup-go-worker/example/internal/wsconn"
	"log"
)

func (c *Client) UnPublish(streamId uint64) error {
	if err := wsconn.NewWsClient(c.wsOpt).UnPublish(streamId); err != nil {
		log.Println("unpublish error:", err)
		return err
	}
	log.Println("unpublish completed, streamId:", streamId)
	return nil
}
