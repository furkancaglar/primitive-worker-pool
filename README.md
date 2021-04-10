### Example project

This is a project I used to show some features of Go lang to a friend. So nothing serious here :)

The app initiates a number of workers that can work in parallel and then gives them jobs (just an int value that makes
them sleep as many milliseconds. After saying this I think I should have called the project as `sleepy workers` :D )
to execute. When the initiator goroutine is out of jobs to give, it closes the channel
so each worker then can receive signal to stop excepting any job. Along the process we log the number of goroutines to
see what is going on with the goroutine number.

Type `go run main.go` command in the root project directory to run the application.
