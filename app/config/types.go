package config

type (
	Config struct {
		Gin    GinConfig
		Sqlite SqliteConfig
	}

	Application struct {
		Name    string
		Version string
	}

	GinConfig struct {
		Port  int
		Debug bool
	}

	SqliteConfig struct {
		File string
	}
)
