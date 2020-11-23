package mxsocket

import (
	asLib "maunium.net/go/mautrix/appservice"
)

var as *asLib.AppService

type mxSocket struct {
	SendQueue chan map[string]string
	RecvQueue chan map[string]string
}

func (ms *mxSocket) Init(config string) error { /*
		var err error
		as, err = asLib.Load(config)
		if err != nil {
			return err
		}*/
	return nil
}

func (ms *mxSocket) Start() error { /*
		as.Start()*/
	return nil
}

func New() mxSocket {
	return mxSocket{}
}
