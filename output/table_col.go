package output

import (
	"regexp"
	"strings"
)

var (
	rxWords       = regexp.MustCompile(`\s+`)
	rxPrefixSplit = regexp.MustCompile(`^(\s+)(.+)$`)
	rxRightTrim   = regexp.MustCompile(`\s+$`)
)


// NewTableCol creates a new table column object from given content
func NewTableCol(content string) *TableCol {
	this := &TableCol{}
	return this.SetContent(content)
}

// Render returns wrapped content, max content line width and content line count.
// See `LineCount()`, `Width()` and `Content()` for more informations.
func (this *TableCol) Render(maxWidth uint) (content string, width, lineCount uint) {
	if this.lineCount == 1 && (maxWidth == 0 || maxWidth >= this.width) {
		content = *this.content
		width = this.width
		if maxWidth > 0 {
			if diff := maxWidth - this.width; diff > 0 {
				content += strings.Repeat(" ", int(diff))
				width += diff
			}
		}
		lineCount = this.lineCount
	} else {
		content = this.Content(maxWidth)
		lines := strings.Split(content, "\n")
		for idx, line := range lines {
			l := uint(StringLength(line))
			if l > width {
				width = l
			}
			if maxWidth > 0 {
				if diff := maxWidth - l; diff > 0 {
					lines[idx] += strings.Repeat(" ", int(diff))
				}
			}
			lineCount++
		}
		content = strings.Join(lines, "\n")
	}
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
	if len(maxWidth) == 0 || maxWidth[0] == 0 || maxWidth[0] >= this.width {
		return *this.content
	}
	return WrapStringExtreme(*this.content, maxWidth[0])
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
	if len(maxWidth) == 0 || maxWidth[0] == 0 {
		return this.width
	}
	lines := strings.Split(this.Content(maxWidth...), "\n")
	return this.maxLineWidth(lines)
}

func (this *TableCol) SetContent(content string) *TableCol {
	content = rxRightTrim.ReplaceAllString(content, "")
	this.content = &content
	lines := strings.Split(content, "\n")
	this.width = this.maxLineWidth(lines)
	this.lineCount = uint(len(lines))
	return this
}

func (this *TableCol) maxLineWidth(lines []string) (width uint) {
	for _, line := range lines {
		if l := uint(StringLength(line)); l > width {
			width = l
		}
	}
	return
}

func (this *TableCol) SetRow(row *TableRow) *TableCol {
	this.row = row
	return this
}
