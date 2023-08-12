package server

import (
	"github.com/gofiber/fiber/v2"
	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type httpServerConfigRoute struct {
	Data map[string]interface{}
}

type httpServerConfigRouteItem struct {
	Path     string
	Handlers []fiber.Handler
}

type httpServerConfigGroup struct {
	Path     string
	Handlers []fiber.Handler
	Children []interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	httpServerConfigRoute
//----------------------------------------------------------------------------------------------------------------------

func NewHttpServerConfigRoute() *httpServerConfigRoute {
	instance := new(httpServerConfigRoute)
	instance.Data = make(map[string]interface{})
	return instance
}

func (instance *httpServerConfigRoute) Group(path string, handlers ...fiber.Handler) *httpServerConfigGroup {
	m := instance.Data
	g := &httpServerConfigGroup{
		Path:     path,
		Handlers: handlers,
	}
	g.Children = make([]interface{}, 0)
	m[buildKey("GROUP", path)] = g
	return g
}

func (instance *httpServerConfigRoute) All(path string, handlers ...fiber.Handler) {
	m := instance.Data
	m[buildKey("ALL", path)] = &httpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
}

func (instance *httpServerConfigRoute) Get(path string, handlers ...fiber.Handler) {
	m := instance.Data
	m[buildKey("GET", path)] = &httpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
}

func (instance *httpServerConfigRoute) Post(path string, handlers ...fiber.Handler) {
	m := instance.Data
	m[buildKey(fiber.MethodPost, path)] = &httpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
}

func (instance *httpServerConfigRoute) Delete(path string, handlers ...fiber.Handler) {
	m := instance.Data
	m[buildKey(fiber.MethodDelete, path)] = &httpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
}

func (instance *httpServerConfigRoute) Put(path string, handlers ...fiber.Handler) {
	m := instance.Data
	m[buildKey(fiber.MethodPut, path)] = &httpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
}


//----------------------------------------------------------------------------------------------------------------------
//	httpServerConfigGroup
//----------------------------------------------------------------------------------------------------------------------

func (instance *httpServerConfigGroup) All(path string, handlers ...fiber.Handler) {
	g := NewHttpServerConfigRoute()
	m := g.Data
	m[buildKey("ALL", path)] = &httpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
	instance.Children = append(instance.Children, g)
}

func (instance *httpServerConfigGroup) Get(path string, handlers ...fiber.Handler) {
	g := NewHttpServerConfigRoute()
	m := g.Data
	m[buildKey("GET", path)] = &httpServerConfigRouteItem{
		Path:     path,
		Handlers: handlers,
	}
	instance.Children = append(instance.Children, g)
}


//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func buildKey(method, path string) string {
	return method + "_" + qbc.Coding.MD5(path)
}
