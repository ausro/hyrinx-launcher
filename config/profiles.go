package config

import (
	"github.com/ausro/hyrinx-launcher/internal/hyrinx"
)

// Defines a single Profile containing current application information
//   - name: The name of the Profile
//   - applications: All applications added to this Profile
type Profile struct {
	Name         string
	Applications []hyrinx.Application
}

func CreateProfile(name string) *Profile {
	newProfile := Profile{
		Name:         name,
		Applications: []hyrinx.Application{},
	}

	return &newProfile
}

func (p *Profile) Update(name string, apps []hyrinx.Application) {
	p.Name = name
	p.Applications = apps
	WriteConfig(CONF)
}
