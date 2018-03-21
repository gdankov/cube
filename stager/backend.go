package stager

import (
	"encoding/json"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/runtimeschema/cc_messages"
	"github.com/julz/cube"
	"github.com/julz/cube/opi"
)

//type Config struct {
//TaskDomain     string
//StagerURL      string
//FileServerURL  string
//CCUploaderURL  string
//Lifecycles     map[string]string
//SkipCertVerify bool
//}

type backend struct {
	config cube.BackendConfig
	logger lager.Logger
}

func NewBackend(config cube.BackendConfig, logger lager.Logger) cube.Backend {
	return &backend{
		config: config,
		logger: logger.Session("kubernetes"),
	}
}

func (b backend) CreateStagingTask(stagingGuid string, request cc_messages.StagingRequestFromCC) (opi.Task, error) {
	logger := b.logger.Session("create-staging-task", lager.Data{"app-id": request.AppId, "staging-guid": stagingGuid})
	logger.Info("staging-request")

	var lifecycleData cc_messages.BuildpackStagingData
	err := json.Unmarshal(*request.LifecycleData, &lifecycleData)
	if err != nil {
		return opi.Task{}, err
	}

	stagingTask := opi.Task{
		Image: "diegoteam/recipe:build",
		Env: map[string]string{
			"DOWNLOAD_URL":        lifecycleData.AppBitsDownloadUri,
			"UPLOAD_URL":          lifecycleData.DropletUploadUri,
			"APP_ID":              request.AppId,
			"STAGING_GUID":        stagingGuid,
			"COMPLETION_CALLBACK": request.CompletionCallback,
			"CF_USERNAME":         b.config.CfUsername,
			"CF_PASSWORD":         b.config.CfPassword,
			"API_ADDRESS":         b.config.ApiAddress,
		},
	}

	return stagingTask, nil
}

func (b backend) BuildStagingResponse() (cc_messages.StagingResponseForCC, error) {
	//TODO
	return cc_messages.StagingResponseForCC{}, nil
}
