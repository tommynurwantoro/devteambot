package container

import (
	"fmt"

	"devteambot/internal/pkg/logger"

	"github.com/facebookgo/inject"
)

type Graph interface {
	Provide(objects ...*inject.Object) error
	Populate() error
}

type Service interface {
	Startup() error
	Shutdown() error
}

type Container interface {
	Ready() error
	GetServiceOrNil(id string) interface{}
	RegisterService(id string, svc interface{})
	Shutdown()
}

type container struct {
	graph    Graph
	order    []string
	ready    bool
	services map[string]interface{}
}

func New() Container {
	return &container{
		graph:    new(inject.Graph),
		order:    make([]string, 0),
		services: make(map[string]interface{}, 0),
		ready:    false,
	}
}

// Ready starts up the service graph and returns error if it's not ready
func (c *container) Ready() error {
	if c.ready {
		return nil
	}
	if err := c.graph.Populate(); err != nil {
		return err
	}
	for _, key := range c.order {
		obj := c.services[key]
		if s, ok := obj.(Service); ok {
			logger.Info(fmt.Sprintf("[starting up] %s", key))
			if err := s.Startup(); err != nil {
				return err
			}
		}
	}
	c.ready = true
	return nil
}

func (c *container) RegisterService(id string, svc interface{}) {
	err := c.graph.Provide(&inject.Object{Name: id, Value: svc, Complete: false})
	if err != nil {
		logger.Panic(fmt.Sprintln("Error provide", id), err)
	}
	c.order = append(c.order, id)
	c.services[id] = svc
}

func (c *container) GetServiceOrNil(id string) interface{} {
	svc, ok := c.services[id]
	if !ok {
		panic(fmt.Errorf("service %s nil", id))
	}
	return svc
}

func (c *container) Shutdown() {
	for _, key := range c.order {
		if service, ok := c.services[key]; ok {
			if s, ok := service.(Service); ok {
				logger.Info(fmt.Sprint("[shutting down] ", key))
				if err := s.Shutdown(); err != nil {
					logger.Error(fmt.Sprintf("ERROR: [shutting down] %s", key), err)
				}
			}
		}
	}
	c.ready = false
}
