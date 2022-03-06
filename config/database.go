package config

type Database struct {
	Username           string `config:"USERNAME"`
	Password           string `config:"PASSWORD"`
	Port               string `config:"PORT"`
	Host               string `config:"HOST"`
	DBName             string `config:"DB_NAME"`
	MaxOpenConnections int    `config:"MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections int    `config:"MAX_IDLE_CONNECTIONS"`
	Timezone           string `config:"TIMEZONE"`
}
