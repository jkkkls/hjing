package {{lowServiceName}}

// {{serviceName}}Service 服务
type {{serviceName}}Service struct {
}

func (service *{{serviceName}}Service) NodeConn(name string)                  {}
func (service *{{serviceName}}Service) NodeClose(name string)                 {}
func (service *{{serviceName}}Service) OnEvent(eventName string, args ...any) {}

// Exit 退出处理
func (service *{{serviceName}}Service) Exit() {}

// Run 服务启动函数
func (service *{{serviceName}}Service) Run() error {return nil}
