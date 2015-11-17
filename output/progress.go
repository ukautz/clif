package output

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

type (
	ProgressBarStyle struct {
		Progress    string
		Rightmost   string
		None        string
		LeftBorder  string
		RightBorder string
	}

	ProgressBar struct {

		// RenderTimeEstimate enables output of estimated time
		RenderTimeEstimate bool

		// RenderProgressCount enables output of progressed count vs size
		RenderProgressCount bool

		// RenderProgressPercentage enables output of progressed percentage
		RenderProgressPercentage bool

		// RefreshTime is the time
		RefreshTime time.Duration

		// size is total length of progress
		size int

		// position is current location in progress
		position int

		// RenderWidth is the max size of characters
		renderWidth int

		// style is the rendering style (ASCII, UTF-8) of the progress bar.. i.e. the characters used
		style *ProgressBarStyle

		// started holds the time when Start() was called
		started time.Time

		mux *sync.Mutex
	}
)

const (
	PB_DEFAULT_REFRESH_TIME = time.Millisecond * 200
	PB_DEFAULT_RENDER_WIDTH = 80
)

var (
	// PbOutOfBoundError is returned when increasing, adding or setting the position
	// with a value which is below 0 or beyond size of the progress bar
	PbOutOfBoundError = fmt.Errorf("Position is out of bounds")

	// PbStyleAscii is an ASCII encoding based style for rendering the progress bar
	PbStyleAscii = &ProgressBarStyle{
		Progress:    "=",
		Rightmost:   ">",
		None:        "-",
		LeftBorder:  "[",
		RightBorder: "]",
	}

	// PbStyleAscii is an UTF-8 encoding based style for rendering the progress bar
	PbStyleUtf8 = &ProgressBarStyle{
		Progress:    "█",
		Rightmost:   "▓",
		None:        "░",
		LeftBorder:  "▕",
		RightBorder: "▏",
	}

	// PbEraseLine prints on output to erease current line
	PbEraseLine = func(out io.Writer) {
		out.Write([]byte(fmt.Sprintf("\033[2K\r")))
	}
)

func NewProgressBar(size int) *ProgressBar {
	return &ProgressBar{
		size:                     size,
		renderWidth:              PB_DEFAULT_RENDER_WIDTH,
		RenderTimeEstimate:       true,
		RenderProgressPercentage: true,
		RefreshTime:              PB_DEFAULT_REFRESH_TIME,
		mux:                      new(sync.Mutex),
		style:                    PbStyleUtf8,
	}
}

// Add increases progress by given amount
func (this *ProgressBar) Add(amount int) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	if this.position+amount > this.size {
		return PbOutOfBoundError
	}
	this.position += amount
	return nil
}

// Finish ends progress by skipping to the end
func (this *ProgressBar) Finish() {
	this.Set(this.size)
}

// Increase adds a one to progress
func (this *ProgressBar) Increase() error {
	return this.Add(1)
}

// Position returns the current position
func (this *ProgressBar) Position() int {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.position
}

// Set moves progress to given position (must be between 0 and size)
func (this *ProgressBar) Set(position int) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	if position < 0 || position > this.size {
		return PbOutOfBoundError
	}
	this.position = position
	return nil
}

// SetRenderWidth is builder method to set render width (defaults to PB_DEFAULT_RENDER_WIDTH)
func (this *ProgressBar) SetRenderWidth(v int) *ProgressBar {
	this.renderWidth = v
	return this
}

// SetRenderWidth is builder method to set size (i.e. max length) of progress
func (this *ProgressBar) SetSize(v int) *ProgressBar {
	this.size = v
	return this
}

// SetStyle sets the rendering (output) style (i.e. the characters used to print the progress bar)
func (this *ProgressBar) SetStyle(v *ProgressBarStyle) *ProgressBar {
	this.style = v
	return this
}

// Start runs the continuous rendering of the progress bar. The returned channel
// can be used to pre-maturely close the progress bar - and to check if the
// progress bar has finished:
//
//      // close before finished
//      pbc := pb.Start(out)
//      errc := make(chan error)
//      pbc <- errc
//      <-errc // now closed
//
//      // wait for finished
//      pbc := pb.Start(out)
//      <-pbc // now finished
func (this *ProgressBar) Start(out io.Writer) chan chan error {
	this.started = time.Now()
	donec := make(chan chan error)
	go func() {
		defer close(donec)
		out.Write([]byte(this.Render()))
		lastPos := this.Position()
		next := time.After(this.RefreshTime)
		for {
			select {
			case errc := <-donec:
				PbEraseLine(out)
				out.Write([]byte(this.Render()))
				close(errc)
				return
			case <-next:
				pos := this.Position()
				if pos >= this.size {
					this.Finish()
					PbEraseLine(out)
					out.Write([]byte(this.Render()))
					return
				} else if this.RenderTimeEstimate || pos != lastPos {
					PbEraseLine(out)
					out.Write([]byte(this.Render()))
				}
				next = time.After(this.RefreshTime)
				lastPos = pos
			}
		}
	}()
	return donec
}

// Style returns the currently used style
func (this *ProgressBar) Style() *ProgressBarStyle {
	return this.style
}

// Render returns rendered progress bar in current progress position
func (this *ProgressBar) Render() string {
	this.mux.Lock()
	defer this.mux.Unlock()
	out := ""
	size := this.size
	if size == 0 {
		size = 1
	}
	percentage := float32(this.position) * 100 / float32(size)
	infoPrefix := this.buildProgressInfoPrefix(percentage, size)
	infoSuffix := this.buildProgressInfoSuffix(percentage, size)
	infoSize := StringLength(infoPrefix) + StringLength(infoSuffix) +
	StringLength(this.style.LeftBorder) + StringLength(this.style.RightBorder)
	width := this.renderWidth - infoSize
	if width == 0 {
		width = 1
	}
	out += infoPrefix
	out += this.style.LeftBorder
	progress := int(percentage * float32(width) / 100)
	none := width - progress
	if progress > 1 {
		if progress < width {
			out += strings.Repeat(this.style.Progress, progress-1)
		} else {
			out += strings.Repeat(this.style.Progress, progress)
		}
	}
	//fmt.Printf("PROGRESS1: %d -- NONE1: %d", progress, none)
	if progress > 0 && progress < width {
		out += this.style.Rightmost
	}
	diff := width - (none + progress)
	//fmt.Printf(" -- DIFF: %d", diff)
	if diff != 0 {
		none += diff
	}
	//fmt.Printf(" -- NONE2: %d", none)
	/*if progress == 0 {
		none++
	}*/
	//fmt.Printf(" -- WIDTH: %d -- PROGRESS2: %d -- SIZE: %d", width, progress, size)
	if none > 0 {
		out += strings.Repeat(this.style.None, none)
	}

	out += this.style.RightBorder
	out += infoSuffix
	//fmt.Printf(" -- TOTAL: %d\n", len(out))

	return out
}

func (this *ProgressBar) buildProgressInfoSuffix(percentage float32, size int) string {
	out := []string{}
	if this.RenderProgressPercentage {
		if percentage > 99.99 {
			out = append(out, "100.%")
		} else {
			out = append(out, fmt.Sprintf("%4.1f%%", percentage))
		}
	}
	if this.RenderTimeEstimate {
		if this.started.Year() == 1 {
			this.started = time.Now()
		}
		duration := time.Now().Sub(this.started)
		var expected uint64
		if this.position > 0 {
			// dur / pos = total / size
			// total = dur * size / pos
			expected = uint64(float64(duration.Nanoseconds()) * float64(size) / float64(this.position))
		}
		//fmt.Printf("> EXPECTED: %d .. %.1f -- DURATON: %.1f\n", expected, time.Duration(expected).Seconds(), duration.Seconds())
		remaining := time.Duration(expected - uint64(duration.Nanoseconds()))
		out = append(out, this.renderFixedSizeDuration(remaining))
	}

	if len(out) > 0 {
		return fmt.Sprintf(" [%s]", strings.Join(out, " - "))
	}
	return ""
}

func (this *ProgressBar) buildProgressInfoPrefix(percentage float32, size int) string {
	out := []string{}
	if this.RenderProgressCount {
		l := fmt.Sprintf("%d", len(fmt.Sprintf("%d", size)))
		out = append(out, fmt.Sprintf("%"+l+"d/%d", this.position, size))
	}
	if len(out) > 0 {
		return fmt.Sprintf("[%s] ", strings.Join(out, " - "))
	}
	return ""
}

var (
	pbYear = float64(365 * 24)
	pbWeek = float64(7 * 24)
	pbDay  = float64(24)
)

func (this *ProgressBar) renderFixedSizeDuration(dur time.Duration) string {
	h := dur.Hours()
	m := dur.Minutes()
	if h > pbYear*10 {
		y := int(h / pbYear)
		h -= float64(y) * pbYear
		w := h / pbWeek
		return fmt.Sprintf("%02dy%02dw", y, int(w))
	} else if h > pbWeek*10 {
		return fmt.Sprintf("%05dw", int(h/pbWeek))
	} else if h > pbDay*2 {
		d := int(h / pbDay)
		h -= pbDay * float64(d)
		return fmt.Sprintf("%02dd%02dh", d, int(h))
	} else if h > 1 {
		o := int(h)
		i := m - float64(o)*60
		return fmt.Sprintf("%02dh%02dm", o, int(i))
	} else if dur.Seconds() < 0 {
		return "00m00s"
	} else {
		i := int(m)
		s := dur.Seconds() - float64(i)*60
		return fmt.Sprintf("%02dm%02ds", i, int(s))
	}
}
