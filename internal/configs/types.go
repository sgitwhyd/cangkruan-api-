package configs

type (
	Config struct {
		Port      string `mapstructure:"PORT"`
		SecretJWT string `mapstructure:"SECRET_JWT"`
		Database  string `mapstructure:"DATABASE_URL"`
	}
)