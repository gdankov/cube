package stager

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/runtimeschema/cc_messages"
	"github.com/julienschmidt/httprouter"
	"github.com/julz/cube"
)

func New(stager cube.Stager, backend cube.Backend, logger lager.Logger) http.Handler {
	handler := httprouter.New()

	stagingHandler := NewStagingHandler(stager, backend, logger)

	handler.PUT("/v1/staging/:staging_guid", stagingHandler.Stage)
	handler.DELETE("/v1/staging/:staging_guid", stagingHandler.StopStaging)
	handler.POST("/v1/staging/:staging_guid/completed", stagingHandler.StagingComplete)

	//stagingCompletedHandler := NewStagingCompletionHandler(logger, ccClient, backends, clock)

	return handler
}

type StagingHandler struct {
	stager  cube.Stager
	backend cube.Backend
	logger  lager.Logger
}

func NewStagingHandler(stager cube.Stager, backend cube.Backend, logger lager.Logger) *StagingHandler {
	logger = logger.Session("staging-handler")

	return &StagingHandler{
		stager:  stager,
		backend: backend,
		logger:  logger,
	}
}

func (handler *StagingHandler) Stage(resp http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	stagingGuid := ps.ByName("staging_guid")
	logger := handler.logger.Session("staging-request", lager.Data{"staging-guid": stagingGuid})

	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("read-body-failed", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	var stagingRequest cc_messages.StagingRequestFromCC
	err = json.Unmarshal(requestBody, &stagingRequest)
	if err != nil {
		logger.Error("unmarshal-request-failed", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	envVars := []string{}
	for _, envVar := range stagingRequest.Environment {
		envVars = append(envVars, envVar.Name)
	}

	logger.Info("environment", lager.Data{"keys": envVars})

	stagingTask, err := handler.backend.CreateStagingTask(stagingGuid, stagingRequest)
	if err != nil {
		logger.Error("building-receipe-failed", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = handler.stager.Run(stagingTask)
	if err != nil {
		logger.Error("stage-app-failed", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusAccepted)
}

func (handler *StagingHandler) StopStaging(resp http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//TODO
}

func (handler *StagingHandler) StagingComplete(resp http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//TODO
}

//Wrap httprouter.Hanlde for testing
func TestHandler(f httprouter.Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, nil)
	}
}
