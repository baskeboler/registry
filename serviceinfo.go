package main

import (
	"sync"
	"time"
)

const (
	CHECK_FAIL = 0
	CHECK_OK
)

type CheckResult int

type HealthCheck struct {
	LastCheck       time.Time
	LastCheckResult CheckResult
}
type EndpointEntry struct {
	Host   string        `json:"hostname"`
	Port   int           `json:"port"`
	Url    string        `json:"url"`
	Health []HealthCheck `json:"health-checks"`
}
type ServiceRegistry struct {
	ServiceName string
	Entries     []*EndpointEntry
}

type EndpointRegistration struct {
	ServiceName string `json:"serviceName"`
	Url         string `json:"url`
	Hostname    string `json:"hostname"`
	Port        int    `json:"port"`
}

type ServicesInfo struct {
	sync.RWMutex
	Data map[string]*ServiceRegistry
}

func NewServicesInfo() *ServicesInfo {
	return &ServicesInfo{Data: make(map[string]*ServiceRegistry)}
}

func (s *ServicesInfo) Get(key string) ([]EndpointEntry, bool) {
	(*s).RLock()
	defer (*s).RUnlock()
	reg, ok := s.Data[key]
	var ret []EndpointEntry
	if ok {
		for _, v := range reg.Entries {
			var copy EndpointEntry
			copy = *v
			ret = append(ret, copy)
		}
		return ret, true
	}
	return ret, false
}

func (s *ServicesInfo) Set(key string, registration EndpointRegistration) {
	s.Lock()
	defer s.Unlock()
	reg, ok := s.Data[key]
	if !ok {
		reg = &ServiceRegistry{ServiceName: key, Entries: []*EndpointEntry{}}
		s.Data[key] = reg
	}
	entry := &EndpointEntry{Host: registration.Hostname, Port: registration.Port, Url: registration.Url, Health: []HealthCheck{}}
	reg.Entries = append(reg.Entries, entry)
}

func (s *ServicesInfo) Unset(key string, registration EndpointRegistration) {
	s.Lock()
	defer s.Unlock()
	reg, ok := s.Data[key]
	if !ok {
		return
	}
	for i, entry := range reg.Entries {
		matches := entry.matches(&registration)
		if matches {
			reg.Entries = append(reg.Entries[:i], reg.Entries[i+1:]...)
		}
	}
}

func (e *EndpointEntry) matches(reg *EndpointRegistration) bool {
	return reg != nil && e != nil && e.Host == reg.Hostname && e.Port == reg.Port && e.Url == reg.Url
}
