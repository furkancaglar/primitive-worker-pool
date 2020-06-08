package main

import (
	"math/rand"
	"time"
)

func main() {

	run()

}

// application configs
const (
	NUMBER_OF_GOROUTINES = 5
	NUMBER_OF_JOBS       = 10
)

const (
	RANDOM_INTEGERS_MAX = 10000
	RANDOM_INTEGERS_MIN = 1
)

type (
	// send all of the go routines to let them know there is no more
	// job coming so they shut themselves down.
	// that means each go routine will read from this channel exactly once
	// so a basic for loop will be enough to wait for all go routines to
	// finish their last job.
	signalEndOfJobs struct{}
)

var (
	chnInteger         chan int32
	chnSignalEndOfJobs chan signalEndOfJobs
)

func run() {

	rand.Seed(time.Now().UnixNano())

	chnInteger = make(chan int32)
	chnSignalEndOfJobs = make(chan signalEndOfJobs)

	initGoRoutines()

	writeNumbers()

	waitForLastGoRoutine()
}

func initGoRoutines() {
	for i := 0; i < NUMBER_OF_GOROUTINES; i++ {
		go execute()
	}
}

func writeNumbers() {
	for i := 0; i < NUMBER_OF_JOBS; i++ {
		rnd := rand.Int31n(RANDOM_INTEGERS_MAX-RANDOM_INTEGERS_MIN) + RANDOM_INTEGERS_MIN
		chnInteger <- rnd
	}
}

func waitForLastGoRoutine() {
	for i := 0; i < NUMBER_OF_GOROUTINES; i++ {
		chnSignalEndOfJobs <- struct{}{}
	}
}

func execute() {
	for {
		select {
		case i, ok := <-chnInteger:
			if !ok {
				return
			}
			time.Sleep(time.Millisecond * time.Duration(i))
		case <-chnSignalEndOfJobs:
			return
		}
	}
}
