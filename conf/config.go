package conf

type Config struct {
	Version             string              `json:"Version"`
	Key                 string              `json:"Key"`
	DbSettings          DbSettings          `json:"DbSettings"`
	EmailSenderSettings EmailSenderSettings `json:"EmailSenderSettings"`
	RedisSettings       RedisSettings       `json:"RedisSettings"`
	EsSettings          EsSettings          `json:"EsSettings"`
	HostName            string              `json:"HostName"`
}
type EsSettings struct {
	Address string `json:"Address"`
	Port    string `json:"Port"`
}
type RedisSettings struct {
	Address  string `json:"Address"`
	Password string `json:"Password"`
	Port     string `json:"Port"`
}
type DbSettings struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Hostname string `json:"Hostname"`
	Dbname   string `json:"Dbname"`
}

// EmailSenderSettings Case Sensitive! Need a uppercase suffix
type EmailSenderSettings struct {
	ServerAddress  string `json:"ServerAddress"`
	ServerPort     int    `json:"ServerPort"`
	ServerHost     string `json:"ServerHost"`
	ServerPassword string `json:"ServerPassword"`
}
