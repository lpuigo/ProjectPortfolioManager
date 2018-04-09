package eventbus

import "github.com/oskca/gopherjs-vue"

var eventBus *vue.ViewModel

func NewEventBus() {
	eventBus = vue.New(nil, nil)
}

func Get() *vue.ViewModel {
	return eventBus
}
