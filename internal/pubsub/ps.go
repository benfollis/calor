package pubsub

import "fmt"

type PubSub struct {
	subscribers map[string][]chanPacket
	capacity int
}

type chanPacket struct {
	channel chan interface{}
	seq int
}

func Initialize(capacity int) *PubSub {
	subs := make(map[string][]chanPacket)
	ps := PubSub{
		subscribers: subs,
		capacity:    capacity,
	}
	return &ps
}

/** This is not thread safe, so set up all the subscribers in the main thread

 */
func (ps *PubSub) Subscribe(topic string) chan interface{} {
	fmt.Println(ps)
	current := ps.subscribers[topic]
	newChan := make(chan interface{}, ps.capacity)
	packet := chanPacket{
		channel: newChan,
		seq: len(current),
	}
	fmt.Println("New subscription", packet)
	newCurrent := append(current, packet)
	ps.subscribers[topic] = newCurrent
	return newChan
}

/** Note
This will block the producer if any consumer hits the capacity depth
*/
func (ps *PubSub) Publish(topic string, message interface{}) {
	subs := ps.subscribers[topic]
	for _, chanPacket := range subs {
		fmt.Println("Sending message to chan seq", chanPacket.seq)
		chanPacket.channel <- message
	}
}
