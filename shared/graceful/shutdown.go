package graceful

import (
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"SituationBak/shared/logger"
)

// Closer 可关闭资源接口
type Closer interface {
	Close() error
}

// NamedCloser 带名称的可关闭资源
type NamedCloser struct {
	Name   string
	Closer Closer
}

// Shutdown 优雅关闭管理器
type Shutdown struct {
	mu      sync.Mutex
	closers []NamedCloser
	timeout time.Duration
}

// NewShutdown 创建优雅关闭管理器
func NewShutdown(timeout time.Duration) *Shutdown {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &Shutdown{
		closers: make([]NamedCloser, 0),
		timeout: timeout,
	}
}

// Register 注册需要关闭的资源
func (s *Shutdown) Register(name string, closer Closer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closers = append(s.closers, NamedCloser{Name: name, Closer: closer})
}

// RegisterIO 注册 io.Closer 类型的资源
func (s *Shutdown) RegisterIO(name string, closer io.Closer) {
	s.Register(name, &ioCloserWrapper{closer})
}

// ioCloserWrapper 包装 io.Closer 以实现 Closer 接口
type ioCloserWrapper struct {
	closer io.Closer
}

func (w *ioCloserWrapper) Close() error {
	return w.closer.Close()
}

// Wait 等待关闭信号
// 返回一个 channel，当收到信号时会关闭
func (s *Shutdown) Wait(signals ...os.Signal) <-chan struct{} {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	done := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)

	go func() {
		sig := <-sigChan
		logger.Info("Received shutdown signal", logger.String("signal", sig.String()))
		close(done)
	}()

	return done
}

// WaitWithCallback 等待关闭信号并执行回调
func (s *Shutdown) WaitWithCallback(callback func(), signals ...os.Signal) {
	<-s.Wait(signals...)
	if callback != nil {
		callback()
	}
	s.Close()
}

// Close 关闭所有已注册的资源
// 按照注册的相反顺序关闭（后注册的先关闭）
func (s *Shutdown) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	logger.Info("Starting graceful shutdown", logger.Int("resources", len(s.closers)))

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	// 创建错误收集通道
	errChan := make(chan error, len(s.closers))
	var wg sync.WaitGroup

	// 按相反顺序关闭资源
	for i := len(s.closers) - 1; i >= 0; i-- {
		closer := s.closers[i]
		wg.Add(1)
		go func(nc NamedCloser) {
			defer wg.Done()
			logger.Info("Closing resource", logger.String("name", nc.Name))
			if err := nc.Closer.Close(); err != nil {
				logger.Error("Failed to close resource",
					logger.String("name", nc.Name),
					logger.Err(err),
				)
				errChan <- err
			} else {
				logger.Info("Resource closed successfully", logger.String("name", nc.Name))
			}
		}(closer)
	}

	// 等待所有资源关闭或超时
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("All resources closed successfully")
	case <-ctx.Done():
		logger.Warn("Shutdown timeout, some resources may not be properly closed")
	}

	close(errChan)

	// 收集错误
	var firstErr error
	for err := range errChan {
		if firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

// CloseSequential 按顺序关闭所有资源（阻塞式）
func (s *Shutdown) CloseSequential() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	logger.Info("Starting sequential shutdown", logger.Int("resources", len(s.closers)))

	var firstErr error

	// 按相反顺序关闭资源
	for i := len(s.closers) - 1; i >= 0; i-- {
		closer := s.closers[i]
		logger.Info("Closing resource", logger.String("name", closer.Name))
		if err := closer.Closer.Close(); err != nil {
			logger.Error("Failed to close resource",
				logger.String("name", closer.Name),
				logger.Err(err),
			)
			if firstErr == nil {
				firstErr = err
			}
		} else {
			logger.Info("Resource closed successfully", logger.String("name", closer.Name))
		}
	}

	return firstErr
}

// Global 全局优雅关闭管理器实例
var Global = NewShutdown(30 * time.Second)

// Register 注册资源到全局管理器
func Register(name string, closer Closer) {
	Global.Register(name, closer)
}

// RegisterIO 注册 io.Closer 到全局管理器
func RegisterIO(name string, closer io.Closer) {
	Global.RegisterIO(name, closer)
}

// Wait 等待关闭信号（全局）
func Wait(signals ...os.Signal) <-chan struct{} {
	return Global.Wait(signals...)
}

// Close 关闭所有资源（全局）
func Close() error {
	return Global.Close()
}
