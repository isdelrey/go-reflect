package mesh

import "testing"

func BenchmarkBroadcastWith2Peers(b *testing.B) {
	m1 := New()
	_ = New()

	type Message map[string]interface{}

	for i := 0; i < b.N; i++ {
		m1.Broadcast("test", &Message{
			"value": "test",
		})
	}
}

func BenchmarkBroadcastWith5Peers(b *testing.B) {
	m1 := New()
	_ = New()
	_ = New()
	_ = New()
	_ = New()

	type Message map[string]interface{}

	for i := 0; i < b.N; i++ {
		m1.Broadcast("test", &Message{
			"value": "test",
		})
	}
}
