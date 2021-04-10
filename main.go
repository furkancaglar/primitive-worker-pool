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
	NUMBER_OF_JOBS       = NUMBER_OF_GOROUTINES
)

const (
	RANDOM_INTEGERS_MAX = 10000
	RANDOM_INTEGERS_MIN = 1
)

func run() {

	rand.Seed(time.Now().UnixNano())

	// create a buffered channel with the number of executers
	// so executers won't wait unnecessarily
	chnInteger := make(chan int32, NUMBER_OF_GOROUTINES)
	wg := &sync.WaitGroup{}
	wg.Add(NUMBER_OF_GOROUTINES)

	initGoRoutines(wg, chnInteger)

	fmt.Printf("%v goroutines are ready for executing jobs..\nCurrent total goroutine number=%v\n", NUMBER_OF_GOROUTINES, runtime.NumGoroutine())

	writeNumbers(chnInteger)

	fmt.Printf("\n\nEND of jobs: all jobs are channeled into executer gorouines..\nCurrent total goroutine number=%v\n\n", runtime.NumGoroutine())

	wg.Wait()
}

func initGoRoutines(wg *sync.WaitGroup, chnInteger chan int32) {
	for i := 0; i < NUMBER_OF_GOROUTINES; i++ {
		go execute(wg, chnInteger)
	}
}

func writeNumbers(chnInteger chan<- int32) {
	for i := 0; i < NUMBER_OF_JOBS; i++ {
		rnd := rand.Int31n(RANDOM_INTEGERS_MAX-RANDOM_INTEGERS_MIN) + RANDOM_INTEGERS_MIN
		chnInteger <- rnd
	}

	close(chnInteger)
}

func execute(wg *sync.WaitGroup, chnInteger <-chan int32) {
	for i := range chnInteger {
		// the actually job is just waiting :D
		time.Sleep(time.Millisecond * time.Duration(i))
	}

	fmt.Printf("Current total goroutine number=%v\n", runtime.NumGoroutine())
	wg.Done()
}
