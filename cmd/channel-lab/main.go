package main

import "fmt"

// EnrollmentEvent is a message passed between enrollment producer and notifier.
type EnrollmentEvent struct {
	CourseID int64
	Email    string
}

func producer(out chan<- EnrollmentEvent) {
	out <- EnrollmentEvent{CourseID: 101, Email: "an@tayjava.com"}
	out <- EnrollmentEvent{CourseID: 101, Email: "binh@tayjava.com"}
	out <- EnrollmentEvent{CourseID: 102, Email: "chi@tayjava.com"}
	close(out) // only the send owner closes
}

func notifier(in <-chan EnrollmentEvent, done chan<- struct{}) {
	for ev := range in {
		fmt.Printf("notify course=%d email=%s\n", ev.CourseID, ev.Email)
	}
	done <- struct{}{}
}

func main() {
	events := make(chan EnrollmentEvent) // unbuffered rendezvous
	done := make(chan struct{})

	go producer(events)
	go notifier(events, done)

	<-done // wait without Sleep
	fmt.Println("all done")
}
