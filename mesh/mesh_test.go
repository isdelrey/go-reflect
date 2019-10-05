package mesh

import "testing"

func BenchmarkBroadcastWith2Peers(b *testing.B) {
	type Message map[string]interface{}

	b.RunParallel(func(pb *testing.PB) {
		m := New()
		for pb.Next() {
			m.Broadcast("test", &Message{
				"value": "test",
			})
		}
	})
}
