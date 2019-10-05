package mesh

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

func GetChannel(name string) chan map[string]interface{} {
	channel := make(chan map[string]interface{})

	Subscriptions.Add(name, func(message map[string]interface{}) {
		channel <- message
	})

	return channel
}
