package hyrinx

// Defines an application added by a user.
// 	- name: 	Name of the application
// 	- path: 	Absolute path to the executable
// 	- icon: 	Display icon
//	- options:  Launch args
type Application struct {
	Name    string
	Path    string
	Icon    string
	Options []string
}

func CreateApplication(name string, path string, icon string) *Application {
	newApp := Application{
		Name:    name,
		Path:    path,
		Icon:    icon,
		Options: make([]string, 0),
	}

	return &newApp
}

func (app *Application) EditDetails(name string, path string, icon string) {
	app.Name = name
	if app.Path != path {
		app.Path = path
	}

	if icon != "" {
		app.Icon = icon
	}
}

func (app *Application) EditOptions(opts []string) {
	app.Options = opts
}
