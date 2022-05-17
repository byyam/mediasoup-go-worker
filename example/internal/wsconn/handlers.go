package wsconn

import (
	"encoding/json"
	"errors"
	"github.com/byyam/mediasoup-go-worker/example/internal/isignal"
	"log"
)

func (c *WsClient) UnPublish(streamId uint64) error {
	req := isignal.UnPublishRequest{StreamId: streamId}
	rsp, err := c.Request(isignal.MethodUnPublish, req)
	if err != nil {
		log.Println("request unpublish error", err)
		return err
	}
	log.Printf("rsp:%+v", rsp)
	return nil
}

func (c *WsClient) Publish(req isignal.PublishRequest) (isignal.PublishResponse, error) {
	rsp := isignal.PublishResponse{}
	rspData, err := c.Request(isignal.MethodPublish, req)
	if err != nil {
		log.Println("request publish error", err)
		return rsp, err
	}
	log.Printf("rsp:%+v", rspData)
	if !rspData.OK {
		return rsp, errors.New("rsp not ok")
	}
	if err := json.Unmarshal(rspData.Data, &rsp); err != nil {
		return rsp, err
	}
	log.Printf("get rsp success %+v", rsp)
	return rsp, nil
}

func (c *WsClient) UnSubscribe(streamId uint64, subId string) error {
	req := isignal.UnSubscribeRequest{StreamId: streamId, SubscribeId: subId}
	rsp, err := c.Request(isignal.MethodUnSubscribe, req)
	if err != nil {
		log.Println("request unsubscribe error", err)
		return err
	}
	log.Printf("rsp:%+v", rsp)
	return nil
}

func (c *WsClient) Subscribe(req isignal.SubscribeRequest) (isignal.SubscribeResponse, error) {
	rsp := isignal.SubscribeResponse{}
	rspData, err := c.Request(isignal.MethodSubscribe, req)
	if err != nil {
		log.Println("request subscribe error", err)
		return rsp, err
	}
	log.Printf("rsp:%+v", rspData)
	if !rspData.OK {
		return rsp, errors.New("rsp not ok")
	}
	if err := json.Unmarshal(rspData.Data, &rsp); err != nil {
		return rsp, err
	}
	log.Printf("get rsp success %+v", rsp)
	return rsp, nil
}
