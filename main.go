package main

import (
	"fmt"
)

type jobQ struct {
	q     chan int
	close chan bool
}

func NewJobQ() *jobQ {
	j := &jobQ{
		q:     make(chan int, 30),
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

func (j *jobQ) Close() {
	close(j.q)
	j.close <- true
}

func (j *jobQ) Status() {
	fmt.Println(len(j.q), cap(j.q))
}

func main() {
	jobq := NewJobQ()

	sp := NewSubscriberPool(jobq.q)

	c1 := NewConsumer()
	c1.Consume(jobq.q)

	sp.Close()

	jobq.Close()

	fmt.Println("main closed")

}
