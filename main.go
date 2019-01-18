package main

import (
	"fmt"
	"math/rand"
	"time"
)

type jobQ struct {
	q     chan int
	close chan bool
}

func NewJobQ() *jobQ {
	j := &jobQ{
		q:     make(chan int, 3),
		close: make(chan bool),
	}
	go j.start()
	return j
}

func (j *jobQ) start() {
	select {
	case <-j.close:
		fmt.Println("jobq q channel closed!")
	}
}

func (j *jobQ) Close(s *subscriber) {
	s.close <- true
	j.close <- true
	close(j.q)
}

func (j *jobQ) Status() {
	fmt.Println(len(j.q), cap(j.q))
}

type subscriber struct {
	id    int
	close chan bool
}

func NewSubscriber() *subscriber {
	return &subscriber{
		id:    1,
		close: make(chan bool),
	}
}

func (s *subscriber) subscribe(jobq chan int) {
	for {
		select {
		case <-s.close:
			fmt.Println("subscriber closed!")
			return
		case jobq <- rand.Intn(100):
		default:
			fmt.Println("Channel full. Discarding value")
		}
	}

}

func (s *subscriber) Close() {
	s.close <- true
}

type consumer struct {
	id int
}

func NewConsumer() *consumer {
	return &consumer{
		id: 1,
	}
}

func (c *consumer) consume(jobq chan int) {
	fmt.Println("consumer:", <-jobq)
}

func main() {
	jobq := NewJobQ()

	s1 := NewSubscriber()
	go s1.subscribe(jobq.q)

	c1 := NewConsumer()
	c1.consume(jobq.q)

	jobq.Status()
	jobq.Close(s1)

	fmt.Println("main closed")

	jobq.Status()
	time.Sleep(5 * time.Second)

}
