package config

type Config struct {
	Log    *Log    `yaml:"log"`
	Data   *Data   `yaml:"data"`
	Server *Server `yaml:"server"`
}

// -----------------------------------------
type Data struct {
	Pgsql *Database `yaml:"Pgsql"`
	Redis *Redis    `yaml:"Redis"`
}

type Database struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	DB           string `yaml:"db"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxLifetime  int    `yaml:"max_life_time"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// -----------------------------------------
type Server struct {
	Http *Http `yaml:"http"`
	Mqtt *Mqtt `yaml:"mqtt"`
	Tcp  *Tcp  `yaml:"tcp"`
	Udp  *Udp  `yaml:"udp"`
}

type Http struct {
	Port int `yaml:"port"`
}

type Mqtt struct {
	Broker   string `yaml:"broker"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Tcp struct {
	Port int `yaml:"port"`
}

type Udp struct {
	Port int `yaml:"port"`
}

// -----------------------------------------
type Log struct {
	Level string `yaml:"level"`
}

// -----------------------------------------
