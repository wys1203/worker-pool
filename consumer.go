package main

import "fmt"

type Consumer struct {
	id int
}

func NewConsumer() *Consumer {
	return &Consumer{
		id: 1,
	}
}

func (c *Consumer) Consume(jobq chan int) {
	fmt.Println("Consumer:", <-jobq)
}
