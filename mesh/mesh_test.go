package mesh

import "testing"

func BenchmarkBroadcastWith2Peers(b *testing.B) {
	type message.Message map[string]interface{}

	b.RunParallel(func(pb *testing.PB) {
		m := New(Options{
			Key: "TEST",
		})
		for pb.Next() {
			m.Broadcast("test", &message.Message{
				"value": "test",
			})
		}
	})
}
