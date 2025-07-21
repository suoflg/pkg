package inject

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 定义一些用于测试的组件类型
type Logger struct {
	Initialized bool
}

func (l *Logger) Initialize() error {
	l.Initialized = true
	return nil
}

func (l *Logger) Finalize() {
	l.Initialized = false
}

type Service struct {
	// 注入 Logger
	Log *Logger `inject:""`
}

func (s *Service) Initialize() error {
	if s.Log == nil {
		return errors.New("logger not injected")
	}
	return nil
}

func (s *Service) Finalize() {}

// 测试基本的依赖注入与生命周期管理
func TestContainer_InjectAndLifecycle(t *testing.T) {
	c := NewContainer()

	logger := &Logger{}
	service := &Service{}

	c.RegisterComponent(logger, service)

	err := c.Initialize()
	require.NoError(t, err)
	assert.True(t, logger.Initialized, "Logger should be initialized")

	c.Finalize()
	assert.False(t, logger.Initialized, "Logger should be finalized")
}

// 测试循环依赖检测
type A struct {
	B *B `inject:""`
}
type B struct {
	A *A `inject:""`
}

func TestContainer_CyclicDependency(t *testing.T) {
	c := NewContainer()

	a := &A{}
	b := &B{}

	c.RegisterComponent(a, b)

	err := c.Initialize()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cyclic dependency")
}

// 测试注册非法组件（非指针或非结构体指针）
func TestContainer_InvalidComponentRegistration(t *testing.T) {
	c := NewContainer()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic when registering invalid component")
		}
	}()

	c.RegisterComponent("not a struct pointer")
}

// 测试字段未设置但具有 `inject` 标签时自动实例化
type AutoCreatedDep struct {
	Name string
}

type NeedsAutoDep struct {
	Dep *AutoCreatedDep `inject:""`
}

func TestContainer_AutoCreateDependency(t *testing.T) {
	c := NewContainer()

	needs := &NeedsAutoDep{}
	c.RegisterComponent(needs)

	err := c.Initialize()
	require.NoError(t, err)
	assert.NotNil(t, needs.Dep)
	assert.Equal(t, "*inject.AutoCreatedDep", reflect.TypeOf(needs.Dep).String())
}
