package register

import (
	"embed"
	"errors"
	"sync"

	"github.com/CanftIn/gothafoss/lib/im/http"
	"github.com/CanftIn/gothafoss/lib/im/model"
)

// APIRouter api路由者
type APIRouter interface {
	Route(r *http.IMHttp)
}

// var apiRoutes = make([]APIRouter, 0)

// // Add 添加api
// func Add(r APIRouter) {
// 	apiRoutes = append(apiRoutes, r)
// }

// var taskRoutes = make([]TaskRouter, 0)

// // GetRoutes 获取所有路由者
// func GetRoutes() []APIRouter {
// 	return apiRoutes
// }

// // TaskRouter task路由者
// type TaskRouter interface {
// 	RegisterTasks()
// }

// // AddTask 添加任务
// func AddTask(task TaskRouter) {
// 	taskRoutes = append(taskRoutes, task)
// }

// // GetTasks 获取所有任务
// func GetTasks() []TaskRouter {
// 	return taskRoutes
// }

type ModuleFnc func(ctx interface{}) Module

var modules = make([]ModuleFnc, 0)

type IMDatasourceType int

const (
	IMDatasourceTypeNone        IMDatasourceType = iota
	IMDatasourceTypeSubscribers                  = 1
	IMDatasourceTypeChannelInfo                  = 1 << 1
	IMDatasourceTypeBlacklist                    = 1 << 2
	IMDatasourceTypeWhitelist                    = 1 << 3
	IMDatasourceTypeSystemUIDs                   = 1 << 4
)

func (i IMDatasourceType) Has(d IMDatasourceType) bool {
	return i&d == d
}

var (
	ErrDatasourceNotProcess error = errors.New("datasource not process")
)

type IMDatasource struct {
	// 是否存在数据
	HasData func(channelID string, channelType uint8) IMDatasourceType
	// 获取订阅者
	Subscribers func(channelID string, channelType uint8) ([]string, error)
	// 获取频道信息
	ChannelInfo func(channelID string, channelType uint8) (map[string]interface{}, error)
	// 黑名单列表
	Blacklist func(channelID string, channelType uint8) ([]string, error)
	// 白名单列表
	Whitelist func(channelID string, channelType uint8) ([]string, error)
	// 系统账号
	SystemUIDs func() ([]string, error)
}

type BussDataSource struct {
	// 获取频道详情
	ChannelGet func(channelID string, channelType uint8, loginUID string) (*model.ChannelResp, error)
}

// 模块
type Module struct {
	// 模块名称
	Name string
	// api 路由
	SetupAPI func() APIRouter
	// sql目录
	SQLDir *SQLFS
	// swagger文件
	Swagger string
	// im 数据源
	IMDatasource IMDatasource
	// 业务数据源
	BussDataSource BussDataSource
	// 事件
	Start func() error
	Stop  func() error
}

func AddModule(moduleFnc func(ctx interface{}) Module) {
	modules = append(modules, moduleFnc)
}

type SQLFS struct {
	embed.FS
}

func NewSQLFS(fs embed.FS) *SQLFS {
	return &SQLFS{
		FS: fs,
	}
}

var once sync.Once
var moduleList []Module

func GetModules(ctx any) []Module {
	once.Do(func() {
		for _, m := range modules {
			moduleList = append(moduleList, m(ctx))
		}
	})

	return moduleList
}

func GetModuleByName(name string, ctx any) Module {
	for _, m := range moduleList {
		if m.Name == name {
			return m
		}
	}
	return Module{}
}
