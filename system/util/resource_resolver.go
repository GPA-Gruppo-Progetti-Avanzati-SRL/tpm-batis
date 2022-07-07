package util

import (
	"io/ioutil"
)

type ResourceResolver interface {
	GetResource(resourceName string) ([]byte, error)
}

type ResourceResolverFunc func(resourceName string) ([]byte, error)

func (f ResourceResolverFunc) GetResource(resourceName string) ([]byte, error) {
	return f(resourceName)
}

type fileRegistryCascadeResolver struct {
	registry map[string][]byte
}

func (r *fileRegistryCascadeResolver) GetResource(resourceName string) ([]byte, error) {

	if resource, ok := r.registry[resourceName]; ok {
		return resource, nil
	}
	if content, err := ioutil.ReadFile(resourceName); err != nil {
		return nil, err
	} else {
		return content, nil
	}
}

func NewFileRegistryCascadeResolver(aRegistry map[string][]byte) ResourceResolver {
	return &fileRegistryCascadeResolver{registry: aRegistry}
}
