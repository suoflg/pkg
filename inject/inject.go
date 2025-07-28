// Package inject 提供了一个轻量级的依赖注入容器，支持组件的注册、自动依赖注入、初始化和资源释放。
//
// 使用方式：
//  1. 创建容器：container := inject.NewContainer()
//  2. 注册组件：container.RegisterComponent(&MyService{}, &MyRepository{})
//  3. 初始化组件：err := container.Initialize()
//  4. 释放资源（可选）：container.Finalize()
//
// 支持功能：
//   - 基于结构体字段的自动注入（通过 struct tag `inject` 标记依赖字段）
//   - 检测并避免循环依赖
//   - 生命周期管理：通过实现 Initialize 和 Finalize 接口分别处理初始化和清理逻辑
//
// 要求所有被注入的字段和注册的组件都必须是指向结构体的指针。
package inject

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Initializer 是一个组件初始化接口。
// 如果组件实现了该接口，容器在依赖注入完成后会调用其 Initialize 方法。
// Initialize 方法返回 error，用于报告初始化失败。
type Initializer interface {
	Initialize() error
}

// Finalizer 是一个组件销毁接口。
// 如果组件实现了该接口，容器在销毁阶段会调用其 Finalize 方法，按初始化的逆序执行。
// 可用于资源释放、关闭连接等操作。
type Finalizer interface {
	Finalize()
}

// Container 是依赖注入容器的公共接口。
// 它支持组件的注册、依赖注入的初始化、以及资源的销毁。
type Container interface {
	Initializer
	Finalizer
	RegisterComponent(component ...any)
}

type container struct {
	components  []any
	visiting    map[reflect.Type]bool
	visited     map[reflect.Type]bool
	typeMap     map[reflect.Type]any
	initialized []any
}

func (c *container) isValidDependencyType(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

// RegisterComponent 注册一个或多个组件到容器中。
// 每个组件必须是指向结构体的指针，否则会 panic。
// 注册后，组件会被存储并在初始化阶段用于依赖注入。
func (c *container) RegisterComponent(component ...any) {
	for i := 0; i < len(component); i++ {
		t := reflect.TypeOf(component[i])
		if !c.isValidDependencyType(t) {
			panic("component must be a pointer to a struct")
		}
		c.typeMap[t.Elem()] = component[i]
		c.components = append(c.components, component[i])
	}
}

func (c *container) initComponent(component any) error {
	t := reflect.TypeOf(component).Elem()

	if c.visited[t] {
		return nil
	}

	if c.visiting[t] {
		return fmt.Errorf("cyclic dependency detected: %v", t)
	}
	c.visiting[t] = true

	v := reflect.ValueOf(component).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if _, ok := field.Tag.Lookup("inject"); !ok {
			continue
		}

		fieldType := field.Type
		if !c.isValidDependencyType(fieldType) {
			return fmt.Errorf("field %s (type: %v) must be a pointer to a struct", field.Name, fieldType)
		}

		dep, ok := c.typeMap[fieldType.Elem()]
		if !ok {
			dep = reflect.New(fieldType.Elem()).Interface()
			c.typeMap[fieldType.Elem()] = dep
		}

		fieldValue := v.Field(i)
		if fieldValue.CanSet() {
			fieldValue.Set(reflect.ValueOf(dep))
		} else {
			ptr := unsafe.Pointer(v.Field(i).UnsafeAddr())
			reflect.NewAt(fieldType, ptr).Elem().Set(reflect.ValueOf(dep))
		}

		if err := c.initComponent(dep); err != nil {
			return err
		}
	}

	if init, ok := component.(Initializer); ok {
		if err := init.Initialize(); err != nil {
			return fmt.Errorf("%s initialization failed: %v", t, err)
		}
	}

	c.initialized = append(c.initialized, component)
	c.visited[t] = true
	delete(c.visiting, t)

	return nil
}

// Initialize 执行依赖注入并初始化所有已注册的组件。
// 它会递归处理组件的字段，注入标记为 `inject` 的依赖项，并调用组件的 Initialize 方法（如果实现了 Initializer 接口）。
// 如果检测到循环依赖，会返回错误。
func (c *container) Initialize() error {
	for _, component := range c.components {
		if err := c.initComponent(component); err != nil {
			return err
		}
	}
	return nil
}

// Finalize 按照初始化的逆序调用所有实现了 Finalizer 接口的组件的 Finalize 方法。
// 通常用于资源清理和关闭操作。
func (c *container) Finalize() {
	for i := len(c.initialized) - 1; i >= 0; i-- {
		if f, ok := c.initialized[i].(Finalizer); ok {
			f.Finalize()
		}
	}
}

// NewContainer 创建一个新的依赖注入容器实例，并返回 Container 接口。
// 该容器支持注册组件、自动注入依赖、初始化和销毁组件。
func NewContainer() Container {
	return &container{
		visiting: make(map[reflect.Type]bool),
		visited:  make(map[reflect.Type]bool),
		typeMap:  make(map[reflect.Type]any),
	}
}
