package clif

import (
	"bytes"
	"fmt"
	"github.com/gosuri/uilive"
	"sync"
	"time"
)

type (
	ProgressBarPool interface {

		// Finish sets all progress to max and stops
		Finish() chan bool

		// Has return whether a bar with name exists or not
		Has(bar string) bool

		// Init adds and initializes a new, named progress bar of given size and returns it
		Init(name string, size int) (ProgressBar, error)

		// Start initializes rendering of all registered bars
		Start()

		// Style sets style of registered progress bars
		Style(style *ProgressBarStyle) error

		// With set's output width of registered progress bars
		Width(width int) error
	}

	progressBarPool struct {
		bars    map[string]ProgressBar
		names   []string
		mux     *sync.Mutex
		style   *ProgressBarStyle
		width   int
		refresh time.Duration
		writer  *uilive.Writer
		started bool
		finishc chan bool
	}
)

const (
	PROGRESS_BAR_DEFAULT_REFRESH = time.Millisecond * 50
)

func NewProgressBarPool(style ...*ProgressBarStyle) ProgressBarPool {
	if len(style) == 0 {
		style = []*ProgressBarStyle{ProgressBarStyleUtf8}
	}
	return &progressBarPool{
		bars:    make(map[string]ProgressBar),
		names:   make([]string, 0),
		mux:     new(sync.Mutex),
		refresh: PROGRESS_BAR_DEFAULT_REFRESH,
		style:   style[0],
		width:   TermWidthCurrent,
		writer:  uilive.New(),
		finishc: make(chan bool),
	}
}

func (this *progressBarPool) Increase(bar string, amount int) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	pb, exist := this.bars[bar]
	if !exist {
		return fmt.Errorf("Progress bar with name \"%s\" does not exist", bar)
	}
	pb.Increase(amount)
	return nil
}

func (this *progressBarPool) Init(name string, size int) (ProgressBar, error) {
	this.mux.Lock()
	defer this.mux.Unlock()
	if _, exist := this.bars[name]; exist {
		return nil, fmt.Errorf("Progress bar with name \"%s\" does already exist", name)
	}
	Dbg("Create new bar: %s (size: %d)", name, size)
	this.names = append(this.names, name)
	this.bars[name] = NewProgressBar(size)
	this.bars[name].SetStyle(this.style)
	this.bars[name].SetRenderWidth(this.width)
	return this.bars[name], nil
}

func (this *progressBarPool) Has(bar string) bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	if _, exist := this.bars[bar]; exist {
		return true
	}
	return false
}

func (this *progressBarPool) Increment(bar string, bars ...string) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	names := append([]string{bar}, bars...)
	for _, name := range names {
		if !this.Has(name) {
			return fmt.Errorf("Progress bar with name \"%s\" does not exist", name)
		}
	}
	for _, name := range names {
		this.bars[name].Increment()
	}
	return nil
}

func (this *progressBarPool) Set(bar string, pos int) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	pb, exist := this.bars[bar]
	if !exist {
		return fmt.Errorf("Progress bar with name \"%s\" does not exist", bar)
	}
	return pb.Set(pos)
}

func (this *progressBarPool) Start() {
	this.mux.Lock()
	defer this.mux.Unlock()
	if !this.started {
		this.started = true
		this.writer.Start()

		go func() {
			defer func() {
				this.mux.Lock()
				defer this.mux.Unlock()
				close(this.finishc)
				this.finishc = make(chan bool)
				this.started = false
				for _, bar := range this.bars {
					bar.Reset()
				}
			}()

			tick := time.NewTicker(this.refresh)
			defer tick.Stop()

			for {
				select {

				// finished
				case <- this.finishc:
					this.render()
					<-time.After(time.Millisecond * 5) // to assure rendering is written fully..
					this.writer.Stop()
					return

				// tick:
				case <-tick.C:
					//this.writer.Flush()
					this.render()
				}
			}
		}()
	}
}

func (this *progressBarPool) render() {
	buf := bytes.NewBuffer(nil)
	for _, name := range this.names {
		buf.WriteString(this.bars[name].Render() + "\n")
	}
	this.writer.Write(buf.Bytes())
}

func (this *progressBarPool) Finish() chan bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	for _, name := range this.names {
		this.bars[name].Finish()
	}
	this.finishc <- true
	return this.finishc
}

func (this *progressBarPool) Style(style *ProgressBarStyle) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	if this.started {
		return fmt.Errorf("Cannot set style after start")
	}
	this.style = style
	for _, bar := range this.bars {
		bar.SetStyle(style)
	}
	return nil
}

func (this *progressBarPool) Width(width int) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	if this.started {
		return fmt.Errorf("Cannot set width after start")
	}
	this.width = width
	for _, bar := range this.bars {
		bar.SetRenderWidth(width)
	}
	return nil
}
