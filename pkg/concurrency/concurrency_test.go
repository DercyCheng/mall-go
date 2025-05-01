package concurrency

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptimisticLock(t *testing.T) {
	t.Run("SuccessfulUpdate", func(t *testing.T) {
		lock := NewOptimisticLock()
		version := lock.GetVersion()
		
		success := lock.CheckAndUpdate(version)
		
		assert.True(t, success)
		assert.Equal(t, version+1, lock.GetVersion())
	})
	
	t.Run("ConflictDetection", func(t *testing.T) {
		lock := NewOptimisticLock()
		originalVersion := lock.GetVersion()
		
		// 模拟并发修改
		// 先由一个线程修改
		success1 := lock.CheckAndUpdate(originalVersion)
		// 然后用同样的版本号尝试修改
		success2 := lock.CheckAndUpdate(originalVersion)
		
		assert.True(t, success1)
		assert.False(t, success2) // 应当失败，因为版本号已变更
		assert.Equal(t, originalVersion+1, lock.GetVersion())
	})
}

func TestMemoryMutexManager(t *testing.T) {
	manager := NewMemoryMutexManager()
	
	t.Run("GetMutexCreatesNewIfNotExists", func(t *testing.T) {
		mutex1 := manager.GetMutex("key1")
		mutex2 := manager.GetMutex("key2")
		
		assert.NotNil(t, mutex1)
		assert.NotNil(t, mutex2)
		assert.NotEqual(t, mutex1, mutex2)
	})
	
	t.Run("GetMutexReturnsSameInstance", func(t *testing.T) {
		mutex1 := manager.GetMutex("key3")
		mutex2 := manager.GetMutex("key3")
		
		assert.Equal(t, mutex1, mutex2)
	})
	
	t.Run("MutexCorrectlyLocksResource", func(t *testing.T) {
		mutex := manager.GetMutex("resource")
		shared := 0
		wg := sync.WaitGroup{}
		
		// 创建多个并发的goroutines修改shared变量
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				
				mutex.Lock()
				current := shared
				time.Sleep(time.Microsecond) // 引入延迟，增加竞争可能
				shared = current + 1
				mutex.Unlock()
			}()
		}
		
		wg.Wait()
		assert.Equal(t, 100, shared) // 如果正确加锁，结果应该是100
	})
}

func TestConcurrencyLimiter(t *testing.T) {
	t.Run("LimitsParallelOperations", func(t *testing.T) {
		maxConcurrent := int64(5)
		limiter := NewConcurrencyLimiter(maxConcurrent)
		
		ctx := context.Background()
		currentActive := 0
		maxObserved := 0
		var mu sync.Mutex
		var wg sync.WaitGroup
		
		// 启动10个goroutines，每个尝试获取信号量
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				
				err := limiter.Acquire(ctx)
				assert.NoError(t, err)
				
				mu.Lock()
				currentActive++
				if currentActive > maxObserved {
					maxObserved = currentActive
				}
				mu.Unlock()
				
				// 保持活跃一段时间
				time.Sleep(50 * time.Millisecond)
				
				mu.Lock()
				currentActive--
				mu.Unlock()
				
				limiter.Release()
			}()
		}
		
		wg.Wait()
		assert.LessOrEqual(t, maxObserved, int(maxConcurrent))
	})
	
	t.Run("ExecuteWithLimit", func(t *testing.T) {
		limiter := NewConcurrencyLimiter(2)
		ctx := context.Background()
		
		executed := false
		err := limiter.ExecuteWithLimit(ctx, func() error {
			executed = true
			return nil
		})
		
		assert.NoError(t, err)
		assert.True(t, executed)
	})
	
	t.Run("ExecuteWithLimitError", func(t *testing.T) {
		limiter := NewConcurrencyLimiter(1)
		ctx := context.Background()
		
		// 耗尽许可
		err := limiter.Acquire(ctx)
		assert.NoError(t, err)
		
		// 创建带取消的上下文
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel() // 立即取消
		
		// 尝试执行，应该失败
		executed := false
		err = limiter.ExecuteWithLimit(cancelCtx, func() error {
			executed = true
			return nil
		})
		
		assert.Equal(t, ErrTooManyConcurrentRequests, err)
		assert.False(t, executed)
		
		// 释放许可
		limiter.Release()
	})
}

func TestWithTimeout(t *testing.T) {
	t.Run("CompletesWithinTimeout", func(t *testing.T) {
		ctx := context.Background()
		
		err := WithTimeout(ctx, 100*time.Millisecond, func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)
			return nil
		})
		
		assert.NoError(t, err)
	})
	
	t.Run("TimesOut", func(t *testing.T) {
		ctx := context.Background()
		
		err := WithTimeout(ctx, 50*time.Millisecond, func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond)
			return nil
		})
		
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.DeadlineExceeded))
	})
	
	t.Run("PropagatesError", func(t *testing.T) {
		ctx := context.Background()
		expectedErr := errors.New("test error")
		
		err := WithTimeout(ctx, 100*time.Millisecond, func(ctx context.Context) error {
			return expectedErr
		})
		
		assert.Equal(t, expectedErr, err)
	})
	
	t.Run("RespectsContextCancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()
		
		err := WithTimeout(ctx, 200*time.Millisecond, func(ctx context.Context) error {
			timer := time.NewTimer(100 * time.Millisecond)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
				return nil
			}
		})
		
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled))
	})
}