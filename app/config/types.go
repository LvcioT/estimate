package config

type (
	Config struct {
		App    Application
		Gin    GinGonic
		Sqlite Sqlite
	}

	Application struct {
		Name     string
		FullName string `toml:"full_name"`
		Version  string
	}

	GinGonic struct {
		Port  int
		Debug bool
	}

	Sqlite struct {
		File string
	}
)
