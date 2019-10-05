package sync

type event struct {
	Name    string
	Message interface{}
}

func Broadcast(Name string, Message interface{}) {
	for _, stream := range Streams.list {
		stream.channel <- event{
			Name,
			Message,
		}
	}
}

func Subscribe(name string, handler func(message map[string]interface{})) {
	Subscriptions.Add(name, handler)
}
