package clif

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type (
	ProgressBar interface {

		// Done returns whether progress bar is done (position == size)
		Done() bool

		// Finish sets progress to max and
		Finish()

		// Increase adds given amount to progress
		Increase(amount int) error

		// Increment increases progress by one
		Increment() error

		// Position returns the progress position
		Position() int

		// Render returns the rendered progress bar in it's current position
		Render() string

		// RenderWidth returns the width of the rendered output (amount of chars (rune))
		RenderWidth() int

		// Reset sets position to zero and unsets all timers
		Reset()

		// Set moves the progress to given position
		Set(position int) error

		// SetRenderWidth sets total width of rendered output to given amount of characters (runes)
		SetRenderWidth(v int) ProgressBar

		// SetSize sets the total size of the progress bar
		SetSize(v int) ProgressBar

		// SetStyle sets the style of the progress bar
		SetStyle(style *ProgressBarStyle) ProgressBar

		// Style sets the
		Style() *ProgressBarStyle
	}

	ProgressBarSimple struct {
		// size is total length of progress
		size int

		// position is current location in progress
		position int

		// RenderWidth is the max size of characters
		renderWidth int

		// style is the rendering style (ASCII, UTF-8) of the progress bar.. i.e. the characters used
		style *ProgressBarStyle

		// started holds the first time
		started time.Time

		// stopped holds the time when Start() was called
		stopped time.Time

		mux *sync.Mutex
	}
)

var (
	// PbOutOfBoundError is returned when increasing, adding or setting the position
	// with a value which is below 0 or beyond size of the progress bar
	PbOutOfBoundError = fmt.Errorf("Position is out of bounds")
)

func NewProgressBar(size int) *ProgressBarSimple {
	if size == 0 {
		if s, err := TermWidth(); err != nil {
			size = s
		} else {
			s = TERM_DEFAULT_WIDTH
		}
	}

	return &ProgressBarSimple{
		size:        size,
		renderWidth: TermWidthCurrent,
		mux:         new(sync.Mutex),
		style:       ProgressBarStyleUtf8,
	}
}

// Done returns bool whether progress bar is done (at 100%)
func (this *ProgressBarSimple) Done() bool {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.done()
}

func (this *ProgressBarSimple) done() bool {
	return this.position == this.size
}

// Finish ends progress by skipping to the end
func (this *ProgressBarSimple) Finish() {
	this.Set(this.size)
}

// Add increases progress by given amount
func (this *ProgressBarSimple) Increase(amount int) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	if this.position+amount > this.size {
		return PbOutOfBoundError
	}
	this.position += amount
	this.setTimes()
	return nil
}

// Increase adds a one to progress
func (this *ProgressBarSimple) Increment() error {
	return this.Increase(1)
}

// Position returns the current position
func (this *ProgressBarSimple) Position() int {
	this.mux.Lock()
	defer this.mux.Unlock()
	return this.position
}

func (this *ProgressBarSimple) Reset() {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.position = 0
	this.started = time.Time{}
	this.stopped = time.Time{}
}

// Render returns rendered progress bar in current progress position
func (this *ProgressBarSimple) Render() string {
	this.mux.Lock()
	defer this.mux.Unlock()
	this.setTimes()
	out := ""
	size := this.size
	if size == 0 {
		size = 1
	}
	percentage := float32(this.position) * 100 / float32(size)
	infoPrefix := this.buildProgressInfo(percentage, size, PROGRESS_BAR_ADDON_PREPEND)
	infoSuffix := this.buildProgressInfo(percentage, size, PROGRESS_BAR_ADDON_APPEND)
	infoSize := StringLength(infoPrefix) + StringLength(infoSuffix) +
	StringLength(string(this.style.LeftBorder)) + StringLength(string(this.style.RightBorder))
	width := this.renderWidth - infoSize
	if width == 0 {
		width = 1
	}
	out += infoPrefix
	out += string(this.style.LeftBorder)
	progress := int(percentage * float32(width) / 100)
	none := width - progress
	if progress > 1 {
		if progress < width {
			out += strings.Repeat(string(this.style.Progress), progress-1)
		} else {
			out += strings.Repeat(string(this.style.Progress), progress)
		}
	}
	if progress > 0 && progress < width {
		out += string(this.style.Rightmost)
	}
	diff := width - (none + progress)
	if diff != 0 {
		none += diff
	}
	if none > 0 {
		out += strings.Repeat(string(this.style.None), none)
	}

	out += string(this.style.RightBorder)
	out += infoSuffix

	return out
}

// Render returns rendered progress bar in current progress position
func (this *ProgressBarSimple) RenderWidth() int {
	return this.renderWidth
}

// Set moves progress to given position (must be between 0 and size)
func (this *ProgressBarSimple) Set(position int) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	if position < 0 || position > this.size {
		return PbOutOfBoundError
	}
	this.position = position
	this.setTimes()
	return nil
}

// SetRenderWidth is builder method to set render width (defaults to PB_DEFAULT_RENDER_WIDTH)
func (this *ProgressBarSimple) SetRenderWidth(v int) ProgressBar {
	this.renderWidth = v
	return this
}

// SetSize is builder method to set size (i.e. max length) of progress
func (this *ProgressBarSimple) SetSize(v int) ProgressBar {
	this.size = v
	return this
}

// SetStyle sets the rendering (output) style (i.e. the characters used to print the progress bar)
func (this *ProgressBarSimple) SetStyle(v *ProgressBarStyle) ProgressBar {
	this.style = v
	return this
}

// Style returns the currently used style
func (this *ProgressBarSimple) Style() *ProgressBarStyle {
	return this.style
}

func (this *ProgressBarSimple) buildProgressInfo(percentage float32, size int, pos progressBarAddon) string {
	// count, elapsed, estimate, percentage
	out := []string{"", "", "", ""}

	if this.style.Count == pos {
		out[0] = this.renderCount(size)
	}
	if this.style.Elapsed == pos {
		out[1] = this.renderElapsed()
	}
	if this.style.Estimate == pos {
		out[2] = this.renderEstimate(size)
	}
	if this.style.Percentage == pos {
		out[3] = this.renderPercentage(percentage)
	}

	if pos == PROGRESS_BAR_ADDON_APPEND {
		return this.style.RenderSuffix(out[0], out[1], out[2], out[3])
	} else {
		return this.style.RenderPrefix(out[0], out[1], out[2], out[3])
	}
}

func (this *ProgressBarSimple) setTimes() {
	if this.started.Year() <= 1 {
		this.started = time.Now()
	}
	if this.size == this.position {
		if this.stopped.Year() <= 1 {
			this.stopped = time.Now()
		}
	} else {
		this.stopped = time.Time{}
	}
}

func (this *ProgressBarSimple) renderElapsed() string {
	var duration time.Duration
	if this.done() {
		duration = this.stopped.Sub(this.started)
	} else {
		duration = time.Now().Sub(this.started)
	}
	return this.style.RenderElapsed(duration, this)
}

func (this *ProgressBarSimple) renderCount(size int) string {
	return this.style.RenderCount(this.position, size, this)
}

func (this *ProgressBarSimple) renderPercentage(percentage float32) string {
	return this.style.RenderPercentage(percentage, this)
}

func (this *ProgressBarSimple) renderEstimate(size int) string {
	if this.done() {
		return this.style.RenderEstimate(time.Duration(0), this)
	}
	duration := time.Now().Sub(this.started)
	var expected uint64
	if this.position > 0 {
		// dur / pos = total / size
		// total = dur * size / pos
		expected = uint64(float64(duration.Nanoseconds()) * float64(size) / float64(this.position))
	}
	remaining := time.Duration(expected - uint64(duration.Nanoseconds()))
	return this.style.RenderEstimate(remaining, this)
}

var (
	pbYear = float64(365 * 24)
	pbWeek = float64(7 * 24)
	pbDay  = float64(24)
)

func (this *ProgressBarSimple) renderFixedSizeDuration(dur time.Duration) string {
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

func RenderFixedSizeDuration(dur time.Duration) string {
	h := dur.Hours()
	m := dur.Minutes()
	s := dur.Seconds()
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
	} else if s > 99 {
		i := int(m)
		o := s - float64(i)*60
		return fmt.Sprintf("%02dm%02ds", i, int(o))
	} else if s < 1 {
		return "00m00s"
	} else {
		ms := (s - float64(int(s))) * 1000
		return fmt.Sprintf("%02ds%03d", int(s), int(ms))
	}
}
