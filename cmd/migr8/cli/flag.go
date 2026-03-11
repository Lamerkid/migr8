package cli

type flag struct {
	Name        string
	Description string
}

func (a *App) addFlag(flag *flag) {
	a.flags[flag.Name] = flag
}

// RegisterFlags adds flags to the CLI app.
func RegisterFlags(app *App) {
	app.addFlag(&flag{
		Name:        "-config",
		Description: "Select config file for migr8",
	})

	app.addFlag(&flag{
		Name:        "-dsn",
		Description: "Database connection URL",
	})

	app.addFlag(&flag{
		Name:        "-path",
		Description: "Path to migration files",
	})
}
