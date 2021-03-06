package driver

import (
	"fmt"

	"github.com/spaceuptech/galaxy/model"
	"github.com/spaceuptech/galaxy/utils/auth"
)

// New creates a new instance of the driver module
func New(auth *auth.Module, c *Config) (Driver, error) {

	switch c.DriverType {
	case model.TypeIstio:
		// Generate the config file
		var istioConfig *istio.Config
		if c.IsInCluster {
			istioConfig = istio.GenerateInClusterConfig()
		} else {
			istioConfig = istio.GenerateOutsideClusterConfig(c.ConfigFilePath)
		}
		istioConfig.SetProxyPort(c.ProxyPort)

		return istio.NewIstioDriver(auth, istioConfig)
	default:
		return nil, fmt.Errorf("invalid driver type (%s) provided", c.DriverType)
	}
}

// Config describes the configuration required by the driver module
type Config struct {
	DriverType     model.DriverType
	ConfigFilePath string
	IsInCluster    bool
	ProxyPort      uint32
}

// Driver is the interface of the modules which interact with the deployment targets
type Driver interface {
	CreateProject(project *model.Project) error
	ApplyService(service *model.Service) error
	AdjustScale(service *model.Service, activeReqs int32) error
	WaitForService(service *model.Service) error
	Type() model.DriverType
}
