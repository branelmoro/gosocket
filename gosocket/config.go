package gosocket

import "time"

type serverConfig struct {
	httpReadTimeOut time.Duration
	httpMaxUriSize int
	httpMaxHeaderSize int

	wsMaxFrameSize int
	wsMaxMessageSize int

	wsHeaderReadTimeout time.Duration
	wsMinByteRatePerSec int
	wsCloseReadTimeout time.Duration
}

var serverConf = serverConfig{
	httpReadTimeOut:    20,
    httpMaxUriSize: 	256,
	httpMaxHeaderSize:	8192,

	wsMaxFrameSize:		1024,
	wsMaxMessageSize:	4096,

	wsHeaderReadTimeout:1,
	wsMinByteRatePerSec:100,
	wsCloseReadTimeout: 1,
}

func setConfig(config serverConfig) {
	serverConf = config
}

func GetConfig() serverConfig {
	return serverConf
}
