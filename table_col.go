package clif

import (
	"fmt"
	"regexp"
	"strings"
	"os"
)

var (
	rxWords       = regexp.MustCompile(`\s+`)
	rxPrefixSplit = regexp.MustCompile(`^(\s+)(.+)$`)
	rxRightTrim   = regexp.MustCompile(`\s+$`)

	TableColWrapper = func(limit int) *Wrapper {
		wrapper := NewWrapper(uint(limit))
		wrapper.WhitespaceMode = WRAP_WHITESPACE_CONTRACT
		wrapper.TrimMode = WRAP_TRIM_RIGHT
		wrapper.BreakWords = true
		return wrapper
	}
)

// NewTableCol creates a new table column object from given content
func NewTableCol(row *TableRow, content string) *TableCol {
	this := &TableCol{}
	return this.SetContent(content)
}

// Render returns wrapped content, max content line width and content line count.
// See `LineCount()`, `Width()` and `Content()` for more informations.
func (this *TableCol) Render(maxWidth int) (content string, width, lineCount int) {
	content = this.Content(maxWidth)
	lines := strings.Split(content, "\n")
	rendered := make([]string, len(lines))
	for idx, line := range lines {
		lineLen := StringLength(line)
		if lineLen > width {
			width = lineLen
		}
		if maxWidth > 0 {
			if diff := int(maxWidth) - lineLen; diff > 0 {
				line += strings.Repeat(" ", diff)
			}
		}
		lineCount++
		rendered[idx] = line
		_dbgTableCol("\033[0mLINE %d (len: %d): \"%s\"\n\033[0m", idx, lineLen, line)
	}
	content = strings.Join(rendered, "\n")
	return
}

// LineCount returns the amount of lines of the content. If a `maxWidth` value
// is provided then the returned line count is calculated based on the existing
// content.
//
// Example 1:
//  maxWidth not given
//  content = "foo bar baz"
//  line count = 1
//
// Example 2:
//  maxWidth  = 3
//  content = "foo bar baz"
//  line count = 3
func (this *TableCol) LineCount(maxWidth ...int) int {
	if len(maxWidth) == 0 || maxWidth[0] == 0 {
		return this.lineCount
	}
	return strings.Count(this.Content(maxWidth[0]), "\n") + 1
}

// Content returns the string content of the column. If maxWidth is not given
// or 0 then the original content is returned. Otherwise the returned content
// is wrapped to the maxWidth limitation. Words will be split if they exceed
// the maxWidth value.
func (this *TableCol) Content(maxWidth ...int) string {
	rendered := this.renderedContent()
	_dbgTableCol("\033m-- RENDERED:\n%s\n\033[0m--\n", rendered)
	if len(maxWidth) > 0 && maxWidth[0] > 0 {
		return TableColWrapper(maxWidth[0]).Wrap(rendered)
	}
	return rendered
}

// ContentPrefixed
func (this *TableCol) ContentPrefixed(prefix string, maxWidth ...int) string {
	lines := strings.Split(this.Content(maxWidth...), "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func (this *TableCol) Width(maxWidth ...int) int {
	return this.maxLineWidth(maxWidth...)
}

func (this *TableCol) SetContent(content string) *TableCol {
	content = rxRightTrim.ReplaceAllString(content, "")
	this.content = &content
	lines := strings.Split(content, "\n")
	this.lineCount = len(lines)
	return this
}

func (this *TableCol) SetRenderer(renderer func(string) string) *TableCol {
	this.renderer = renderer
	return this
}

func (this *TableCol) maxLineWidth(maxWidth ...int) (width int) {
	rendered := this.renderedLines(maxWidth...)
	for _, line := range rendered {
		if l := StringLength(line); l > width {
			width = l
		}
	}
	return
}

func (this *TableCol) renderedLines(maxWidth ...int) []string {
	return strings.Split(this.Content(maxWidth...), "\n")
}

func (this *TableCol) renderedContent() string {
	if this.renderer != nil {
		return this.renderer(*this.content)
	}
	return *this.content
}

func (this *TableCol) SetRow(row *TableRow) *TableCol {
	this.row = row
	return this
}

func _dbgTableCol(str string, args ...interface{}) {
	if dbg := os.Getenv("DEBUG_TABLE_COL"); dbg == "yes" || dbg == "1" {
		fmt.Printf(str, args...)
	}
}
