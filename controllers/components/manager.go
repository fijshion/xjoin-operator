package components

import "github.com/go-errors/errors"

type Component interface {
	Name() string
	Create() error
	Delete() error
	CheckDeviation() (error, error)
	Exists() (bool, error)
	SetName(string)
	SetVersion(string)
	ListInstalledVersions() ([]string, error)
}

type ComponentManager struct {
	components []Component
	name       string
	version    string
}

func NewComponentManager(name string, version string) ComponentManager {
	return ComponentManager{
		name:    name,
		version: version,
	}
}

func (c *ComponentManager) AddComponent(component Component) {
	component.SetName(c.name)
	component.SetVersion(c.version)
	c.components = append(c.components, component)
}

// CreateAll creates all components. No-op if the components are already created.
func (c *ComponentManager) CreateAll() error {
	for _, component := range c.components {
		componentExists, err := component.Exists()
		if err != nil {
			return errors.Wrap(err, 0)
		}
		if !componentExists {
			err = component.Create()
			if err != nil {
				return errors.Wrap(err, 0)
			}
		}
	}
	return nil
}

// DeleteAll deletes all components. No-op if the components are already deleted.
func (c *ComponentManager) DeleteAll() error {
	for _, component := range c.components {
		componentExists, err := component.Exists()
		if err != nil {
			return err
		}
		if componentExists {
			err = component.Delete()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CheckForDeviations checks each component's stored value against the expected value, returns true if deviation is found
func (c *ComponentManager) CheckForDeviations() (problems []error, err error) {
	for _, component := range c.components {
		problem, err := component.CheckDeviation()
		if err != nil {
			return problems, errors.Wrap(err, 0)
		}
		if problem != nil {
			problems = append(problems, problem)
		}
	}

	return problems, nil
}
