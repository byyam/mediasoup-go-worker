package mediasoup_go_worker

import (
	"encoding/json"
	"sync"

	"github.com/byyam/mediasoup-go-worker/internal/utils"
	"github.com/byyam/mediasoup-go-worker/mediasoupdata"
	"github.com/byyam/mediasoup-go-worker/mserror"
	"github.com/byyam/mediasoup-go-worker/rtc"
	"github.com/byyam/mediasoup-go-worker/workerchannel"
)

type workerBase struct {
	pid       int
	logger    utils.Logger
	routerMap sync.Map
}

func (w *workerBase) GetPid() int {
	return w.pid
}

func (w *workerBase) OnChannelRequest(request workerchannel.RequestData) (response workerchannel.ResponseData) {

	w.logger.Debug("method=%s,internal=%+v", request.Method, request.Internal)

	switch request.Method {
	case mediasoupdata.MethodWorkerCreateRouter:
		router := rtc.NewRouter(request.Internal.RouterId)
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
		r, ok := w.routerMap.Load(request.Internal.RouterId)
		if !ok {
			response.Err = mserror.ErrRouterNotFound
			return
		}
		router := r.(*rtc.Router)
		router.HandleRequest(request, &response)
	}
	w.logger.Debug("method:%s, response:%s", request.Method, response)
	return
}

func (w *workerBase) Stop() {
	w.routerMap.Range(func(key, value interface{}) bool {
		router := value.(*rtc.Router)
		router.Close()
		return true
	})
	w.logger.Warn("pid:%d is killed", w.pid)
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
	w.logger.Debug("dumpData:%+v", dumpData)
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
