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
		Name:        "-cfg",
		Description: "JSON config file for migration tool",
	})

	app.addFlag(&flag{
		Name:        "-log",
		Description: "Logger level (debug/info/warn/error)",
	})

	app.addFlag(&flag{
		Name:        "-dsn",
		Description: "Database connection URL",
	})

	app.addFlag(&flag{
		Name:        "-dir",
		Description: "Path to directory with migration files",
	})

	app.addFlag(&flag{
		Name:        "-type",
		Description: "Migration type (sql/go)",
	})
}
