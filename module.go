package appbot

import "github.com/nlopes/slack"

var modules = make(map[string]*Module)

// Module : 모듈은
type Module struct {
	name string

	rtm *RTMReceiver
}

func GetModule(name string) *Module {
	if m, ok := modules[name]; ok {
		return m
	}

	m := &Module{
		name: name,
		rtm: newRTMReceiver(),
	}
	modules[name] = m

	return m
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) RTM() *RTMReceiver {
	return m.rtm
}

func (m *Module) run(api *slack.Client) {
	m.rtm.run(api)
}