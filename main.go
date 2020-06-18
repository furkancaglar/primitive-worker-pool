package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	log.Printf("App Starting..\nCurrent total goroutine number=%v\n", runtime.NumGoroutine())

	run()

	// because garbage collector might not run for the last goroutine
	// (in this case the last goroutine will not be removed)
	// so logging will not work correctly
	runtime.GC()

	log.Printf("Sutting Down..\nCurrent total goroutine number=%v\n", runtime.NumGoroutine())

}

// application configs
const (
	NUMBER_OF_GOROUTINES = 7
	NUMBER_OF_JOBS       = 10
)

const (
	RANDOM_INTEGERS_MAX = 10000
	RANDOM_INTEGERS_MIN = 1
)

type (
	// send all of the goroutines to let them know there is no more
	// job coming so they shut themselves down.
	// that means each goroutine will read from this channel exactly once
	// so a basic for loop will be enough to wait for all goroutines to
	// finish their last job.
	signalEndOfJobs struct{}
)

var (
	chnInteger              chan int32
	chnSignalEndOfJobs      chan signalEndOfJobs

	// FOR LOGGING
	goRoutineCounter = NUMBER_OF_GOROUTINES
)

func run() {

	rand.Seed(time.Now().UnixNano())

	chnInteger = make(chan int32)
	chnSignalEndOfJobs = make(chan signalEndOfJobs)

	initGoRoutines()

	log.Printf("%v goroutine ready for executing jobs..\nCurrent total goroutine number=%v\n", NUMBER_OF_GOROUTINES, runtime.NumGoroutine())

	writeNumbers()

	log.Printf("END of jobs: all jobs are channeled into executer gorouines..\nCurrent total goroutine number=%v\n", runtime.NumGoroutine())

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
			log.Printf("number %v goroutine is shutting down..\nCurrent total goroutine number=%v", goRoutineCounter, runtime.NumGoroutine())
			goRoutineCounter--
			return
		}
	}
}
