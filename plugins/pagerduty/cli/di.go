package pagerduty

import (
	"sync"

	"github.com/sendinblue/dpe-insights/core/config"
	"github.com/sendinblue/dpe-insights/core/databases/mysql"
	"github.com/sendinblue/dpe-insights/plugins/pagerduty/models/repositories"
	"github.com/sendinblue/dpe-insights/plugins/pagerduty/services"
)

type diContainer struct {
	pagerdutyClient             *services.Client
	pagerdutyIncidentRepository *repositories.IncidentRepository
}

var (
	container *diContainer
	once      sync.Once
)

func (di *diContainer) GetPagerDutyClient() *services.Client {
	return di.pagerdutyClient
}

func (di *diContainer) GetPagerDutyIncidentRepository() *repositories.IncidentRepository {
	return di.pagerdutyIncidentRepository
}

func newDIContainer() *diContainer {
	conf := config.NewConfig()
	db, _ := mysql.NewDB(conf)

	dic := &diContainer{
		pagerdutyClient:             services.New(conf.PluginPagerDutyAccessToken),
		pagerdutyIncidentRepository: &repositories.IncidentRepository{Db: *db},
	}

	return dic
}

func getDIContainer() *diContainer {
	once.Do(func() {
		container = newDIContainer()
	})
	return container
}
