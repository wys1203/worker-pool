package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Subscriber struct {
	id    int
	close chan bool
}

func NewSubscriber(id int) *Subscriber {
	return &Subscriber{
		id:    id,
		close: make(chan bool),
	}
}

func (s *Subscriber) Subscribe(jobq chan int) {
	for {
		time.Sleep(300 * time.Millisecond)
		select {
		case <-s.close:
			fmt.Println("Subscriber closed!")
			return
		case jobq <- rand.Intn(100):
		default:
			fmt.Println("Channel full. Discarding value")
		}
	}

}

func (s *Subscriber) Close() {
	s.close <- true
}

type SubscriberPool struct {
	subscribers chan *Subscriber
	close       chan bool
}

func NewSubscriberPool(jobq chan int) *SubscriberPool {
	sp := &SubscriberPool{
		subscribers: make(chan *Subscriber, 3),
		close:       make(chan bool),
	}

	for i := 0; i < 3; i++ {
		suber := NewSubscriber(i)
		go suber.Subscribe(jobq)
		sp.subscribers <- suber
	}
	close(sp.subscribers)

	fmt.Println("sp.subscribers:", len(sp.subscribers), cap(sp.subscribers))

	go sp.start()

	return sp
}

func (sp *SubscriberPool) start() {
	select {
	case <-sp.close:
		fmt.Println("SubscriberPool closed!")
	}
}

func (sp *SubscriberPool) Close() {
	fmt.Println("close:", len(sp.subscribers), cap(sp.subscribers))
	for s := range sp.subscribers {
		s.Close()
	}
	sp.close <- true
}
