package output

import (
	"fmt"
)

func NewTableRow(cols []string) *TableRow {
	l := uint(len(cols))
	this := &TableRow{
		ColAmount: l,
		Cols:      make([]*TableCol, l),
	}
	for i, col := range cols {
		this.SetCol(uint(i), NewTableCol(col))
	}
	return this
}

func (this *TableRow) Render(totalWidth uint) (rendered []string, maxLineCount uint) {
	widths := this.CalculateWidths(totalWidth)
	return this.RenderWithWidths(widths)
}

func (this *TableRow) CalculateWidths(totalWidth uint) (colSize []uint) {
	colSize = make([]uint, this.ColAmount)
	factor := float64(1.0)
	if totalWidth > 0 && totalWidth != this.TotalWidth {
		factor = float64(totalWidth) / float64(this.TotalWidth)
	}
	//fmt.Printf("\nTOTAL WIDTH: %d -> %d, FACTOR: %.2f\n", this.TotalWidth, totalWidth, factor)
	usedWidth := uint(0)
	lastCol := int(this.ColAmount - 1)
	for idx, col := range this.Cols {
		width := uint(factor * float64(col.Width()))
		//widthPre := width
		if totalWidth == 0 {
			width = 0
		} else {
			if idx == lastCol { // last col get's it all
				width += totalWidth - (usedWidth + width)
			} else if width == 0 { //
				width = 1
			} else { // collected used width to spend all in last col
				usedWidth += width
			}
		}
		colSize[idx] = width
	}
	return
}

func (this *TableRow) RenderWithWidths(colWidths []uint) (rendered []string, maxLineCount uint) {
	rendered = make([]string, this.ColAmount)
	//fmt.Printf("\nTOTAL WIDTH: %d -> %d, FACTOR: %.2f\n", this.TotalWidth, totalWidth, factor)
	for idx, col := range this.Cols {
		content, _, lineCount := col.Render(colWidths[idx])
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

func (this *TableRow) SetCol(idx uint, col *TableCol) error {
	if idx >= this.ColAmount {
		return fmt.Errorf("Column index %d is beyond colum size %d", idx, this.ColAmount)
	}
	if col.lineCount > this.MaxLineCount {
		this.MaxLineCount = col.lineCount
	}
	if this.Cols[idx] != nil {
		this.TotalWidth -= this.Cols[idx].width
	}
	this.TotalWidth += col.width
	this.Cols[idx] = col
	col.SetRow(this)
	return nil
}
