package config

import (
	"github.com/ausro/hyrinx-launcher/internal/hyrinx"
)

// Defines the Profiles available to be used
//   - profile_list: The list of all current Profiles
type Profiles struct {
	Profile_list []Profile
}

// Defines a single Profile containing current application information
//   - name: The name of the Profile
//   - applications: All applications added to this Profile
type Profile struct {
	Name         string
	Applications map[int]hyrinx.Application
}

func createProfileList() *Profiles {
	return &Profiles{
		Profile_list: make([]Profile, 1),
	}
}

func createProfile(name string) *Profile {
	newProfile := Profile{
		Name:         name,
		Applications: make(map[int]hyrinx.Application),
	}

	return &newProfile
}

func (pl *Profiles) addProfile(pf Profile) {
	pl.Profile_list = append(pl.Profile_list, pf)
}

func (p *Profile) Update(name string, apps map[int]hyrinx.Application) {
	p.Name = name
	p.Applications = apps
	WriteConfig(CONF)
}

func (pl *Profiles) GetCurrentProfile() *Profile {
	return &pl.Profile_list[0]
}

// Export all profiles
func (pl *Profiles) export() {
	//TODO
}
