package cube

import (
	"net/http"

	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/runtimeschema/cc_messages"
	"github.com/julz/cube/opi"
)

type AppInfo struct {
	AppName   string `json:"name"`
	SpaceName string `json:"space_name"`
	AppGuid   string `json:"application_id"`
}

//go:generate counterfeiter . CfClient
type CfClient interface {
	GetDropletByAppGuid(string) ([]byte, error)
	PushDroplet(string, string) error
	GetAppBitsByAppGuid(string) (*http.Response, error)
}

type SyncConfig struct {
	Properties SyncProperties `yaml:"sync"`
}

type SyncProperties struct {
	KubeConfig         string `yaml:"kube_config"`
	RegistryEndpoint   string `yaml:"registry_endpoint"`
	CcApi              string `yaml:"api_endpoint"`
	Backend            string `yaml:"backend"`
	CfUsername         string `yaml:"cf_username"`
	CfPassword         string `yaml:"cf_password"`
	CcUser             string `yaml:"cc_internal_user"`
	CcPassword         string `yaml:"cc_internal_password"`
	SkipSslValidation  bool   `yaml:"skip_ssl_validation"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
}

//go:generate counterfeiter . Stager
type Stager interface {
	Run(task opi.Task) error
}

//go:generate counterfeiter . Backend
type Backend interface {
	CreateStagingTask(string, cc_messages.StagingRequestFromCC) (opi.Task, error)
	BuildStagingResponse(*models.TaskCallbackResponse) (cc_messages.StagingResponseForCC, error)
}

type BackendConfig struct {
	CfUsername        string
	CfPassword        string
	ApiAddress        string
	SkipSslValidation bool
}
