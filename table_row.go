package clif

import (
	"fmt"
	"strings"
)

func NewTableRow(cols []string) *TableRow {
	l := len(cols)
	this := &TableRow{
		ColAmount: l,
		Cols:      make([]*TableCol, l),
	}
	for i, col := range cols {
		this.SetCol(i, NewTableCol(this, col))
	}
	return this
}

func (this *TableRow) Render(totalWidth int) (rendered []string, maxLineCount int) {
	widths := this.CalculateWidths(totalWidth)
	return this.RenderWithWidths(widths)
}

func (this *TableRow) Width(maxWidth ...int) int {
	width := 0
	for _, col := range this.Cols {
		width += col.Width(maxWidth...)
	}
	return width
}

func (this *TableRow) CalculateWidths(totalWidth int) (colSize []int) {
	colSize = make([]int, this.ColAmount)
	factor := float64(1.0)
	actualWidth := this.Width()
	if totalWidth > 0 && totalWidth != actualWidth {
		factor = float64(totalWidth) / float64(actualWidth)
	}
	//fmt.Printf("\nTOTAL WIDTH: %d -> %d, FACTOR: %.2f\n", this.TotalWidth, totalWidth, factor)
	usedWidth := 0
	lastCol := int(this.ColAmount - 1)
	for idx, col := range this.Cols {
		width := factor * float64(col.Width())
		//widthPre := width
		if totalWidth == 0 {
			width = 0
		} else {
			if idx == lastCol { // last col get's it all
				//fmt.Printf("  @@ LAST COL: TOTAL=%d, USED=%d, WIDTH=%d, FACTOR=%.3f\n", totalWidth, usedWidth, width, factor)
				width += float64(totalWidth - (usedWidth + int(width)))
			} else if width == 0 { //
				width = 1
			} else { // collected used width to spend all in last col
				usedWidth += int(width)
			}
		}
		colSize[idx] = int(width)
	}
	return
}

func (this *TableRow) RenderWithWidths(colWidths []int) (rendered []string, maxLineCount int) {
	rendered = make([]string, this.ColAmount)
	//fmt.Printf("\nTOTAL WIDTH: %d -> %d, FACTOR: %.2f\n", this.TotalWidth, totalWidth, factor)
	//fmt.Printf("\nRENDER ROW WITH COL WIDTHS %v\n", colWidths)
	for idx, col := range this.Cols {
		//content, renderWidth, lineCount := col.Render(colWidths[idx])
		content, _, lineCount := col.Render(colWidths[idx])
		m := 0
		for _, c := range strings.Split(content, "\n") {
			if l := len(c); l > m {
				m = l
			}
		}
		//fmt.Printf(" %d) with width = %d, lines = %d, max width = %d, x max width = %d\n", idx, colWidths[idx], lineCount, renderWidth, m)
		rendered[idx] = content
		if lineCount > maxLineCount {
			maxLineCount = lineCount
		}
	}
	return
}

func (this *TableRow) SetTable(table *Table) *TableRow {
	this.table = table
	return this
}

func (this *TableRow) SetCol(idx int, col *TableCol) error {
	if idx >= this.ColAmount {
		return fmt.Errorf("Column index %d is beyond colum size %d", idx, this.ColAmount)
	}
	if col.lineCount > this.MaxLineCount {
		this.MaxLineCount = col.lineCount
	}
	this.Cols[idx] = col
	col.SetRow(this)
	return nil
}

func (this *TableRow) SetRenderer(renderer func(string) string) *TableRow {
	for _, col := range this.Cols {
		col.SetRenderer(renderer)
	}
	return this
}
