package server

import (
	"context"
	"errors"
	log "github.com/hillguo/sanlog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// EnvGraceRestartStr 是否热重启环境变量
const EnvGraceRestartStr = "TRPC_IS_GRACEFUL=1"

// Server trpc server 一个服务进程只有一个server，一个server可以有多个service
type Server struct {
	services map[string]Service // k=serviceName,v=Service
}

// AddService 添加一个service到server里面，service name为配置文件指定的用于名字服务的name
// trpc.NewServer() 会遍历配置文件中定义的service配置项，调用AddService完成serviceName与实现的映射关系
func (s *Server) AddService(serviceName string, service Service) {

	if s.services == nil {
		s.services = make(map[string]Service)
	}
	s.services[serviceName] = service
}

// Service 通过serviceName获取对应的Service
func (s *Server) Service(serviceName string) Service {

	if s.services == nil {
		return nil
	}
	return s.services[serviceName]
}

// Register 把业务实现接口注册到server里面，一般一个server只有一个service，
// 有多个service的情况下请使用 Service("servicename") 指定, 否则默认会把实现注册到server里面的所有service
func (s *Server) Register(serviceDesc interface{}, serviceImpl interface{}) error {

	desc, ok := serviceDesc.(*ServiceDesc)
	if !ok {
		return errors.New("service desc type invalid")
	}

	for _, srv := range s.services {
		err := srv.Register(desc, serviceImpl)
		if err != nil {
			return err
		}
	}

	return nil
}

// Close 通知各个service执行关闭动作,最长等待10s
func (s *Server) Close(ch chan struct{}) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var wg sync.WaitGroup
	for _, service := range s.services {

		wg.Add(1)
		go func(srv Service) {

			defer wg.Done()

			c := make(chan struct{}, 1)
			go srv.Close(c)

			select {
			case <-c:
			case <-ctx.Done():
			}
		}(service)
	}

	wg.Wait()
	if ch != nil {
		ch <- struct{}{}
	}

	return nil
}

// Serve 启动所有服务
func (s *Server) Serve() error {

	if len(s.services) == 0 {
		panic("service empty")
	}

	ch := make(chan os.Signal)
	var failedServices sync.Map
	var err error
	for name, service := range s.services {

		go func(n string, srv Service) {

			e := srv.Serve()
			if e != nil {
				err = e
				failedServices.Store(n, srv)
				time.Sleep(time.Millisecond * 300)
				ch <- syscall.SIGTERM
			}
		}(name, service)
	}

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGSEGV)

	sig := <-ch
	// 热重启单独处理
	if sig == syscall.SIGUSR2 {
		_, err = s.StartNewProcess()
		if err != nil {
			panic(err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var wg sync.WaitGroup
	for name, service := range s.services {

		if _, ok := failedServices.Load(name); ok {
			continue
		}

		wg.Add(1)
		go func(srv Service) {

			defer wg.Done()

			c := make(chan struct{}, 1)
			go srv.Close(c)

			select {
			case <-c:
			case <-ctx.Done():
			}
		}(service)
	}

	wg.Wait()

	if err != nil {
		panic(err)
	}

	return nil
}

// StartNewProcess 启动新进程, 由于trpc使用了reuseport形式，可以直接fork子进程
// 不必再传递net.Listener的文件描述符
func (s *Server) StartNewProcess() (uintptr, error) {

	log.Infof("receive USR2 signal, so restart the process")

	execSpec := &syscall.ProcAttr{
		Env:   append(os.Environ(), EnvGraceRestartStr),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
	}
	fork, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	if err != nil {
		log.Error("failed to forkexec with err: %s", err.Error())
		return 0, err
	}

	return uintptr(fork), nil
}



