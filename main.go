package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Printf("\nApp Starting..\nCurrent total goroutine number=%v\n\n\n", runtime.NumGoroutine())

	run()

	// because garbage collector might not run for the last goroutine
	// (in this case the last goroutine will not be removed)
	// so logging will not work correctly
	runtime.GC()

	fmt.Printf("\n\n\nSutting Down..\nCurrent total goroutine number=%v\n", runtime.NumGoroutine())

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

var (
	chnInteger chan int32
	wg         *sync.WaitGroup

	// FOR LOGGING ONLY
	goRoutineCounter = NUMBER_OF_GOROUTINES
)

func run() {

	rand.Seed(time.Now().UnixNano())

	chnInteger = make(chan int32)
	wg = new(sync.WaitGroup)
	wg.Add(NUMBER_OF_GOROUTINES)

	initGoRoutines()

	fmt.Printf("%v goroutine ready for executing jobs..\nCurrent total goroutine number=%v\n", NUMBER_OF_GOROUTINES, runtime.NumGoroutine())

	writeNumbers()

	fmt.Printf("\n\nEND of jobs: all jobs are channeled into executer gorouines..\nCurrent total goroutine number=%v\n\n", runtime.NumGoroutine())

	wg.Wait()
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

	close(chnInteger)
}

func execute() {
	for {
		select {
		case i, ok := <-chnInteger:
			if !ok {
				fmt.Printf("Current total goroutine number=%v\nnumber %v goroutine is shutting down..\n", runtime.NumGoroutine(), goRoutineCounter)
				goRoutineCounter--
				wg.Done()
				return
			}
			time.Sleep(time.Millisecond * time.Duration(i))
		}
	}
}
