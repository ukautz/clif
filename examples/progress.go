// +build ignore

package main

import (
	"fmt"
	"github.com/ukautz/clif"
	"sync"
	"time"
)

func cmdProgress1(c *clif.Command, out clif.Output) {
	printProgress(c, 1, out)
}

func cmdProgress2(c *clif.Command, out clif.Output) {
	amount := c.Option("amount").Int()
	printProgress(c, amount, out)
}

func cmdProgress3(c *clif.Command, out clif.Output) {
	pbs := out.ProgressBars()
	pbs.Style(clif.ProgressBarStyleAscii)
	if width := c.Option("width").Int(); width > 0 {
		pbs.Width(width)
	}
	pb, _ := pbs.Init("default", 200)
	interval := c.Option("interval").Int()
	pbs.Start()
	var wg sync.WaitGroup
	wg.Add(1)
	go func(b clif.ProgressBar) {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			b.Increment()
			<-time.After(time.Millisecond * time.Duration(interval))
		}
	}(pb)
	wg.Wait()
	<-pbs.Finish()
}

func cmdProgress4(c *clif.Command, out clif.Output) {
	pbs := out.ProgressBars()
	if width := c.Option("width").Int(); width > 0 {
		pbs.Width(width)
	}
	interval := c.Option("interval").Int()
	pbs.Start()
	var wg sync.WaitGroup
	amount := c.Option("amount").Int()
	for i := 0; i < amount; i++ {
		wg.Add(1)
		pb, _ := pbs.Init(fmt.Sprintf("bar-%d", i+1), 200)
		go func(b clif.ProgressBar, ii int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				b.Increment()
				<-time.After(time.Millisecond * time.Duration(interval * ii))
			}
		}(pb, i)
	}
	wg.Wait()
	<-pbs.Finish()
}

func printProgress(c *clif.Command, count int, out clif.Output) {
	out.Printf("<headline>Progressing</headline>\n")
	var wg sync.WaitGroup
	pb := out.ProgressBars()
	if width := c.Option("width").Int(); width > 0 {
		pb.Width(width)
	}
	interval := c.Option("interval").Int()
	style := clif.ProgressBarStyleUtf8
	if c.Option("ascii-style").Bool() {
		style = clif.ProgressBarStyleAscii
	}
	style.Count = clif.PROGRESS_BAR_ADDON_APPEND
	style.Estimate = clif.PROGRESS_BAR_ADDON_PREPEND
	pb.Style(style)
	pb.Start()
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("bar%d", i)
		bar, _ := pb.Init(name, 200)
		wg.Add(1)
		go func(b clif.ProgressBar, t int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				b.Increment()
				//<-time.After(time.Millisecond * time.Duration((t+1)*3))
				<-time.After(time.Millisecond * time.Duration(interval*(t+1)))
			}
		}(bar, i)
	}
	wg.Wait()
	pb.Finish()
}

func main() {
	clif.New("My App", "1.0.0", "An example application").
		New("demo1", "Print the progress bar 1", cmdProgress1).
		New("demo2", "Print the progress bar 2", cmdProgress2).
		New("demo3", "Print the progress bar 3", cmdProgress3).
		New("demo4", "Print the progress bar 3", cmdProgress4).
		AddDefaultOptions(clif.NewFlag("ascii-style", "A", "Use ASCII style instead", false)).
		AddDefaultOptions(clif.NewOption("width", "w", "Width of progress bars", fmt.Sprintf("%d", clif.TermWidthCurrent), false, false)).
		AddDefaultOptions(clif.NewOption("amount", "a", "Amount of progress bars", "3", false, false)).
		AddDefaultOptions(clif.NewOption("interval", "i", "Interval in miliseconds", "10", false, false)).
		SetDefaultCommand("demo1").
		Run()
}
