package clif

import (
	"fmt"
	"strings"
	"time"
)

type (

	// progressBarAddon type of progress par addon position
	progressBarAddon uint

	// ProgressBarStyle
	ProgressBarStyle struct {
		Empty            rune
		Progress         rune
		Rightmost        rune
		None             rune
		LeftBorder       rune
		RightBorder      rune
		Count            progressBarAddon
		Elapsed          progressBarAddon
		Estimate         progressBarAddon
		Percentage       progressBarAddon
		RenderCount      func(pos, max int, bar ProgressBar) string
		RenderElapsed    func(elapsed time.Duration, bar ProgressBar) string
		RenderEstimate   func(forecast time.Duration, bar ProgressBar) string
		RenderPercentage func(percent float32, bar ProgressBar) string
		RenderPrefix     func(count, elapsed, estimate, percentage string) string
		RenderSuffix     func(count, elapsed, estimate, percentage string) string
	}
)

const (
	PROGRESS_BAR_ADDON_OFF progressBarAddon = iota
	PROGRESS_BAR_ADDON_PREPEND
	PROGRESS_BAR_ADDON_APPEND
)

var (
	// ProgressBarDefaultRenderElapsed
	ProgressBarDefaultRenderCount = func(pos, max int, bar ProgressBar) string {
		l := len(fmt.Sprintf("%d", max))
		return fmt.Sprintf("%"+fmt.Sprintf("%d", l)+"d/%d", pos, max)
	}
	ProgressBarDefaultRenderElapsed = func(elapsed time.Duration, bar ProgressBar) string {
		return fmt.Sprintf("@%s", RenderFixedSizeDuration(elapsed))
	}
	ProgressBarDefaultRenderEstimated = func(forecast time.Duration, bar ProgressBar) string {
		return fmt.Sprintf("~%s", RenderFixedSizeDuration(forecast))
	}
	ProgressBarDefaultRenderPercentage = func(percent float32, bar ProgressBar) string {
		if percent > 99.99 {
			return " 100%"
		} else {
			return fmt.Sprintf("%4.1f%%", percent)
		}
	}
	ProgressBarDefaultRenderPrefix = func(count, elapsed, estimate, percentage string) string {
		out := []string{}
		for _, s := range []string{count, elapsed, estimate, percentage} {
			if s != "" {
				out = append(out, s)
			}
		}
		if len(out) == 0 {
			return ""
		}
		return strings.Join(out, " / ") + " "
	}
	ProgressBarDefaultRenderSuffix = func(count, elapsed, estimate, percentage string) string {
		out := []string{}
		for _, s := range []string{count, elapsed, estimate, percentage} {
			if s != "" {
				out = append(out, s)
			}
		}
		if len(out) == 0 {
			return ""
		}
		return " " + strings.Join(out, " / ")
	}

	// ProgressBarStyleAscii is an ASCII encoding based style for rendering the progress bar
	ProgressBarStyleAscii = &ProgressBarStyle{
		Empty:       '-',
		Progress:    '=',
		Rightmost:   '>',
		None:        '-',
		LeftBorder:  '[',
		RightBorder: ']',
		Percentage:  PROGRESS_BAR_ADDON_APPEND,
		Elapsed:     PROGRESS_BAR_ADDON_PREPEND,
	}

	// ProgressBarStyleUtf8 is an UTF-8 encoding based style for rendering the progress bar
	ProgressBarStyleUtf8 = &ProgressBarStyle{
		Empty:       ' ',
		Progress:    '█',
		Rightmost:   '▓',
		None:        '░',
		LeftBorder:  '▕',
		RightBorder: '▏',
		Percentage:  PROGRESS_BAR_ADDON_APPEND,
		Elapsed:     PROGRESS_BAR_ADDON_PREPEND,
	}
)

func CloneProgressBarStyle(from *ProgressBarStyle) *ProgressBarStyle {
	to := *from
	return &to
}

func init() {
	for _, s := range []*ProgressBarStyle{ProgressBarStyleAscii, ProgressBarStyleUtf8} {
		s.RenderCount = ProgressBarDefaultRenderCount
		s.RenderElapsed = ProgressBarDefaultRenderElapsed
		s.RenderEstimate = ProgressBarDefaultRenderEstimated
		s.RenderPercentage = ProgressBarDefaultRenderPercentage
		s.RenderPrefix = ProgressBarDefaultRenderPrefix
		s.RenderSuffix = ProgressBarDefaultRenderSuffix
	}

}
