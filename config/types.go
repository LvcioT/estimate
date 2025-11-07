package config

type (
	AppConfig struct {
		Gin    GinConfig
		Sqlite SqliteConfig
	}

	GinConfig struct {
		Port  int
		Debug bool
	}

	SqliteConfig struct {
		File string
	}
)
