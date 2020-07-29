package config

type ServerConfig struct {
	ServerName string
	Services []ServiceConfig
}

type ServiceConfig struct {
	Name string
	NetWork string
	Address string

	InMsgChanSize uint32
	OutMsgChanSize uint32
	ReadTimeout uint32
	WriteTimeout uint32

	Transport string
	Protocol string
}

type RedisConfig struct {

}

type MysqlConfig struct {

}

type LogConfig struct {

}