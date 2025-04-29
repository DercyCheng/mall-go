package concurrency

import (
	"context"
	"errors"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

var (
	// ErrVersionConflict 版本冲突错误
	ErrVersionConflict = errors.New("version conflict detected, please retry")

	// ErrLockTimeout 获取锁超时错误
	ErrLockTimeout = errors.New("timeout while waiting for lock")

	// ErrTooManyConcurrentRequests 并发请求过多错误
	ErrTooManyConcurrentRequests = errors.New("too many concurrent requests")
)

// OptimisticLock 乐观锁实现
type OptimisticLock struct {
	// 版本号，用于检测并发修改
	version int64
}

// NewOptimisticLock 创建新的乐观锁
func NewOptimisticLock() *OptimisticLock {
	return &OptimisticLock{
		version: 1,
	}
}

// GetVersion 获取当前版本号
func (l *OptimisticLock) GetVersion() int64 {
	return l.version
}

// CheckAndUpdate 检查并更新版本号
// 如果expectedVersion与当前版本匹配，则增加版本号并返回true
// 如果不匹配，则返回false，表示发生了并发修改
func (l *OptimisticLock) CheckAndUpdate(expectedVersion int64) bool {
	if l.version == expectedVersion {
		l.version++
		return true
	}
	return false
}

// DistributedLock 分布式锁接口
type DistributedLock interface {
	// Lock 获取锁
	Lock(ctx context.Context, key string, ttl time.Duration) (bool, error)

	// Unlock 释放锁
	Unlock(ctx context.Context, key string) error
}

// MemoryMutexManager 内存锁管理器
// 注意：这个实现仅适用于单实例场景，分布式系统应使用Redis或etcd等外部存储
type MemoryMutexManager struct {
	locks map[string]*sync.Mutex
	mutex sync.Mutex
}

// NewMemoryMutexManager 创建内存锁管理器
func NewMemoryMutexManager() *MemoryMutexManager {
	return &MemoryMutexManager{
		locks: make(map[string]*sync.Mutex),
	}
}

// GetMutex 获取指定键的互斥锁
func (m *MemoryMutexManager) GetMutex(key string) *sync.Mutex {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.locks[key]; !exists {
		m.locks[key] = &sync.Mutex{}
	}

	return m.locks[key]
}

// ConcurrencyLimiter 并发限制器
type ConcurrencyLimiter struct {
	sem *semaphore.Weighted
}

// NewConcurrencyLimiter 创建并发限制器
func NewConcurrencyLimiter(maxConcurrent int64) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{
		sem: semaphore.NewWeighted(maxConcurrent),
	}
}

// Acquire 获取许可
func (l *ConcurrencyLimiter) Acquire(ctx context.Context) error {
	return l.sem.Acquire(ctx, 1)
}

// Release 释放许可
func (l *ConcurrencyLimiter) Release() {
	l.sem.Release(1)
}

// ExecuteWithLimit 在并发限制下执行函数
func (l *ConcurrencyLimiter) ExecuteWithLimit(ctx context.Context, fn func() error) error {
	if err := l.sem.Acquire(ctx, 1); err != nil {
		return ErrTooManyConcurrentRequests
	}
	defer l.sem.Release(1)

	return fn()
}

// WithTimeout 使用超时执行函数
func WithTimeout(ctx context.Context, timeout time.Duration, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return fn(ctx)
}
