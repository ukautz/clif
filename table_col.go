package clif

import (
	"regexp"
	"strings"
	"fmt"
)

var (
	rxWords       = regexp.MustCompile(`\s+`)
	rxPrefixSplit = regexp.MustCompile(`^(\s+)(.+)$`)
	rxRightTrim   = regexp.MustCompile(`\s+$`)
)

// NewTableCol creates a new table column object from given content
func NewTableCol(row *TableRow, content string) *TableCol {
	this := &TableCol{}
	return this.SetContent(content)
}

// Render returns wrapped content, max content line width and content line count.
// See `LineCount()`, `Width()` and `Content()` for more informations.
func (this *TableCol) Render(maxWidth uint) (content string, width, lineCount uint) {
	content = this.Content(maxWidth)
	lines := strings.Split(content, "\n")
	rendered := make([]string, len(lines))
	for idx, line := range lines {
		lineLen := StringLength(line)
		if uint(lineLen) > width {
			width = uint(lineLen)
		}
		if maxWidth > 0 {
			if diff := int(maxWidth) - lineLen; diff > 0 {
				//fmt.Printf(">> EXTEND LINE (%d VS %d) BY %d\n", lineLen, maxWidth, diff)
				line += strings.Repeat(" ", diff)
			}
		}
		lineCount++
		rendered[idx] = line
		fmt.Printf("\033[0mLINE %d (len: %d): \"%s\"\n\033[0m", idx, lineLen, line)
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
func (this *TableCol) LineCount(maxWidth ...uint) uint {
	if len(maxWidth) == 0 || maxWidth[0] == 0 {
		return this.lineCount
	}
	return uint(strings.Count(this.Content(maxWidth[0]), "\n") + 1)
}

// Content returns the string content of the column. If maxWidth is not given
// or 0 then the original content is returned. Otherwise the returned content
// is wrapped to the maxWidth limitation. Words will be split if they exceed
// the maxWidth value.
func (this *TableCol) Content(maxWidth ...uint) string {
	rendered := this.renderedContent()
	fmt.Printf("\033m-- RENDERED:\n%s\n\033[0m--\n", rendered)
	if len(maxWidth) > 0 && maxWidth[0] > 0 {
		return WrapStringExtreme(rendered, maxWidth[0])
	}
	return rendered
}

// ContentPrefixed
func (this *TableCol) ContentPrefixed(prefix string, maxWidth ...uint) string {
	lines := strings.Split(this.Content(maxWidth...), "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func (this *TableCol) Width(maxWidth ...uint) uint {
	return this.maxLineWidth(maxWidth...)
}

func (this *TableCol) SetContent(content string) *TableCol {
	content = rxRightTrim.ReplaceAllString(content, "")
	this.content = &content
	lines := strings.Split(content, "\n")
	this.lineCount = uint(len(lines))
	return this
}

func (this *TableCol) SetRenderer(renderer func(string) string) *TableCol {
	this.renderer = renderer
	return this
}

func (this *TableCol) maxLineWidth(maxWidth ...uint) (width uint) {
	rendered := this.renderedLines(maxWidth...)
	for _, line := range rendered {
		if l := uint(StringLength(line)); l > width {
			width = l
		}
	}
	return
}

func (this *TableCol) renderedLines(maxWidth ...uint) []string {
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
