package cube

import (
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
}

//******** STAGING TASK BY JULZ.S
type StagingTask struct {
	Actions               map[string]Action
	CompletionCallbackUrl string
}

type Action interface {
	ActionType() string
}

type DownloadAction struct {
	Artifact string `json:"artifact"`
	From     string `json:"from"`
	To       string `json:"to"`
	User     string `json:"user"`
}

func (a DownloadAction) ActionType() string {
	return "ActionTypeDownload"
}

type UploadAction struct {
	Artifact string `json:"artifact,omitempty"`
	From     string `json:"from"`
	To       string `json:"to"`
	User     string `json:"user"`
}

func (a UploadAction) ActionType() string {
	return "ActionTypeUpload"
}

type StagingAction struct {
	User  string
	From  string
	To    string
	Image string
}

func (a StagingAction) ActionType() string {
	return "StagingActionType"
}

//go:generate counterfeiter . WorkspaceManager
type WorkspaceManager interface {
	Create(string) (Workspace, error)
}

type Workspace struct {
	DropletUploadPath   string
	AppBitsDownloadPath string
}
