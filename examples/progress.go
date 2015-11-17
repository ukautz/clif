// +build ignore

package main

import (
	"gopkg.in/ukautz/clif.v1"
	"sync"
	"time"
)

func printProgress(out clif.Output) {
	out.Printf("<headline>Progressing</headline>\n")
	pb := out.Progress(200)
	pb.RenderProgressPercentage = true
	pb.RenderTimeEstimate = true
	pb.RenderProgressCount = true
	var wg sync.WaitGroup
	wg.Add(1)
	pbc := pb.Start(out.Writer())
	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			pb.Add(1)
			<-time.After(time.Millisecond * 20)
		}
		errc := make(chan error)
		pbc <- errc
		if err := <-errc; err != nil {
			out.Printf("\n <error>Failure: %s</error>\n", err)
		} else {
			out.Printf("\n <info>Done</info>\n")
		}
	}()
	wg.Wait()
}

func main() {
	clif.New("My App", "1.0.0", "An example application").
		New("demo", "Print the progress bar", printProgress).
		SetDefaultCommand("demo").
		Run()
}
