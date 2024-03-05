package mediasoup_go_worker

import (
	"encoding/json"
	"sync"

	"github.com/rs/zerolog"

	"github.com/byyam/mediasoup-go-worker/mserror"
	"github.com/byyam/mediasoup-go-worker/pkg/mediasoupdata"
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

	w.logger.Info().Str("request", request.String()).Msg("handle channel request start")

	switch request.Method {
	case mediasoupdata.MethodWorkerCreateRouter:
		router := rtc.NewRouter(request.Internal.RouterId)
		if router == nil {
			response.Err = mserror.ErrInvalidParam
			return
		}
		w.routerMap.Store(request.Internal.RouterId, router)
	case mediasoupdata.MethodWorkerClose:
		w.Stop()
	case mediasoupdata.MethodWorkerDump:
		response.Data = w.FillJson()
	case mediasoupdata.MethodWorkerGetResourceUsage:
		response.Data = w.FillJsonResourceUsage()
	case mediasoupdata.MethodWorkerUpdateSettings:
		// todo
	default:
		h, err := workerchannel.GetChannelRequestHandler(request.HandlerId)
		if err != nil {
			response.Err = err
			return
		}
		h(request, &response)

		//r, ok := w.routerMap.Load(request.Internal.RouterId)
		//if !ok {
		//	response.Err = mserror.ErrRouterNotFound
		//	return
		//}
		//router := r.(*rtc.Router)
		//router.HandleRequest(request, &response)
	}
	w.logger.Info().Str("request", request.String()).Msg("handle channel request done")
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

func (w *workerBase) FillJson() json.RawMessage {
	var routerIds []string
	w.routerMap.Range(func(key, value interface{}) bool {
		routerIds = append(routerIds, key.(string))
		return true
	})
	dumpData := mediasoupdata.WorkerDump{
		Pid:       w.pid,
		RouterIds: routerIds,
	}
	data, _ := json.Marshal(&dumpData)
	w.logger.Debug().Msgf("dumpData:%+v", dumpData)
	return data
}

func (w *workerBase) FillJsonResourceUsage() json.RawMessage {
	// todo
	ruData := mediasoupdata.WorkerResourceUsage{
		RU_Utime:    0,
		RU_Stime:    0,
		RU_Maxrss:   0,
		RU_Ixrss:    0,
		RU_Idrss:    0,
		RU_Isrss:    0,
		RU_Minflt:   0,
		RU_Majflt:   0,
		RU_Nswap:    0,
		RU_Inblock:  0,
		RU_Oublock:  0,
		RU_Msgsnd:   0,
		RU_Msgrcv:   0,
		RU_Nsignals: 0,
		RU_Nvcsw:    0,
		RU_Nivcsw:   0,
	}
	data, _ := json.Marshal(&ruData)
	return data
}
