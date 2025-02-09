package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"bikraj.movie_microservice.net/pkg/discovery"
)

type ServiceName string
type InstanceID string

type Registry struct {
	sync.RWMutex
	serviceAddrs map[ServiceName]map[InstanceID]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewRegistry() *Registry {
	return &Registry{serviceAddrs: make(map[ServiceName]map[InstanceID]*serviceInstance)}
}

func (r *Registry) Register(ctx context.Context, isntanceID InstanceID, serviceName ServiceName, hostPort string) error {
	r.Lock()

	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName]; !ok {
		r.serviceAddrs[serviceName] = map[InstanceID]*serviceInstance{}
	}
	r.serviceAddrs[serviceName][isntanceID] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

func (r *Registry) DeRegister(ctx context.Context, instanceID InstanceID, serviceName ServiceName) error {
	r.Lock()

	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return nil
	}

	delete(r.serviceAddrs[serviceName], instanceID)
	return nil
}

func (r *Registry) ReportHealthyState(instanceID InstanceID, servicename ServiceName) error {

	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[servicename]; !ok {
		return errors.New("service is not registered yet")
	}
	if _, ok := r.serviceAddrs[servicename][instanceID]; !ok {
		return errors.New("service instance is not registered yet")
	}

	r.serviceAddrs[servicename][instanceID].lastActive = time.Now()
	return nil

}

func (r *Registry) ServiceAddresses(ctx context.Context, serviceName ServiceName) ([]string, error) {
	r.RLock()
	defer r.Unlock()

	if len(r.serviceAddrs[serviceName]) == 0 {
		return nil, discovery.ErrNotFound
	}

	var res []string

	for _, i := range r.serviceAddrs[serviceName] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}
