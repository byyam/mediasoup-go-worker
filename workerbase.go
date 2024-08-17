package mediasoup_go_worker

import (
	"sync"

	"github.com/rs/zerolog"

	FBS__Request "github.com/byyam/mediasoup-go-worker/fbs/FBS/Request"
	FBS__Response "github.com/byyam/mediasoup-go-worker/fbs/FBS/Response"
	FBS__Worker "github.com/byyam/mediasoup-go-worker/fbs/FBS/Worker"
	"github.com/byyam/mediasoup-go-worker/mserror"
	"github.com/byyam/mediasoup-go-worker/pkg/zerowrapper"
	"github.com/byyam/mediasoup-go-worker/rtc"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type workerBase struct {
	pid       int
	logger    zerolog.Logger
	routerMap sync.Map
}

func NewWorkerBase(pid int) *workerBase {
	w := &workerBase{
		pid:       pid,
		logger:    zerowrapper.NewScope("worker", pid),
		routerMap: sync.Map{},
	}
	// init channel handlers
	workerchannel.InitChannelHandlers()

	return w
}

func (w *workerBase) GetPid() int {
	return w.pid
}

func (w *workerBase) OnChannelRequest(request workerchannel.RequestData) (response workerchannel.ResponseData) {

	w.logger.Info().Any("body", request.Data).Str("header", request.String()).Msg("[channelMsg]request")

	// set from request
	response.Id = request.Request.Id
	response.MethodType = request.MethodType

	switch request.MethodType {
	case FBS__Request.MethodWORKER_CREATE_ROUTER:
		requestT := request.Request.Body.Value.(*FBS__Worker.CreateRouterRequestT)
		router := rtc.NewRouter(requestT.RouterId)
		if router == nil {
			response.Err = mserror.ErrInvalidParam
			return
		}
		w.routerMap.Store(requestT.RouterId, router)
		response.RspBody = &FBS__Response.BodyT{}
	case FBS__Request.MethodWORKER_CLOSE:
		w.Stop()
	case FBS__Request.MethodWORKER_DUMP:
		// request body is null
		dataDump := w.FillJson()
		// set rsp
		rspBody := &FBS__Response.BodyT{
			Type:  FBS__Response.BodyWorker_DumpResponse,
			Value: dataDump,
		}
		response.RspBody = rspBody
	case FBS__Request.MethodWORKER_GET_RESOURCE_USAGE:
		// request body is null
		dataDump := w.FillJsonResourceUsage()
		// set rsp
		rspBody := &FBS__Response.BodyT{
			Type:  FBS__Response.BodyWorker_ResourceUsageResponse,
			Value: dataDump,
		}
		response.RspBody = rspBody
	case FBS__Request.MethodWORKER_UPDATE_SETTINGS:
	// todo
	case FBS__Request.MethodWORKER_CLOSE_ROUTER:
		requestT := request.Request.Body.Value.(*FBS__Worker.CloseRouterRequestT)
		v, ok := w.routerMap.Load(requestT.RouterId)
		if !ok {
			response.Err = mserror.ErrRouterNotFound
			return
		}
		router := v.(*rtc.Router)
		router.Close()

	default:
		h, err := workerchannel.GetChannelRequestHandler(request.HandlerId)
		if err != nil {
			response.Err = err
			return
		}
		h(request, &response)

	}
	return
}

func (w *workerBase) Stop() {
	w.routerMap.Range(func(key, value interface{}) bool {
		router := value.(*rtc.Router)
		router.Close()
		return true
	})
	w.logger.Warn().Int("pid", w.pid).Msg("worker is killed")
}

func (w *workerBase) FillJson() *FBS__Worker.DumpResponseT {
	var routerIds []string
	w.routerMap.Range(func(key, value interface{}) bool {
		routerIds = append(routerIds, key.(string))
		return true
	})

	channelRequestHandlerIds := &FBS__Worker.ChannelMessageHandlersT{
		ChannelRequestHandlers:      workerchannel.GetChannelRequestHandlerStats(),
		ChannelNotificationHandlers: nil,
	}

	data := &FBS__Worker.DumpResponseT{
		Pid:                    uint32(w.pid),
		WebRtcServerIds:        nil,
		RouterIds:              routerIds,
		ChannelMessageHandlers: channelRequestHandlerIds,
		Liburing:               nil,
	}
	return data
}

func (w *workerBase) FillJsonResourceUsage() *FBS__Worker.ResourceUsageResponseT {

	ruData := &FBS__Worker.ResourceUsageResponseT{
		RuUtime:    0,
		RuStime:    0,
		RuMaxrss:   0,
		RuIxrss:    0,
		RuIdrss:    0,
		RuIsrss:    0,
		RuMinflt:   0,
		RuMajflt:   0,
		RuNswap:    0,
		RuInblock:  0,
		RuOublock:  0,
		RuMsgsnd:   0,
		RuMsgrcv:   0,
		RuNsignals: 0,
		RuNvcsw:    0,
		RuNivcsw:   0,
	}
	return ruData
}
