package output

import (
	"fmt"
	"strings"
)

const (
	DEFAULT_RENDER_WIDTH = uint(80)
)

var (
	DefaultTableStyle = &TableStyle{
		Bottom:          "─",
		ContentRenderer: func(content string) string { return content },
		CrossBottom:     "┴",
		CrossInner:      "┼",
		CrossLeft:       "├",
		CrossRight:      "┤",
		CrossTop:        "┬",
		HeaderRenderer:  func(content string) string { return fmt.Sprintf("\033[1;4m%s\033[0m", content) },
		InnerHorizontal: "─",
		InnerVertical:   "│",
		Left:            "│",
		LeftBottom:      "└",
		LeftTop:         "┌",
		Prefix:          " ",
		Right:           "│",
		RightBottom:     "┘",
		RightTop:        "┐",
		Suffix:          " ",
		Top:             "─",
	}
)

func NewDefaultTableStyle() *TableStyle {
	return copyTableStyle(DefaultTableStyle)
}

func copyTableStyle(from *TableStyle) *TableStyle {
	to := new(TableStyle)
	to.Bottom = from.Bottom
	to.ContentRenderer = from.ContentRenderer
	to.CrossBottom = from.CrossBottom
	to.CrossInner = from.CrossInner
	to.CrossLeft = from.CrossLeft
	to.CrossRight = from.CrossRight
	to.CrossTop = from.CrossTop
	to.HeaderRenderer = from.HeaderRenderer
	to.InnerHorizontal = from.InnerHorizontal
	to.InnerVertical = from.InnerVertical
	to.Left = from.Left
	to.LeftBottom = from.LeftBottom
	to.LeftTop = from.LeftTop
	to.Prefix = from.Prefix
	to.Right = from.Right
	to.RightBottom = from.RightBottom
	to.RightTop = from.RightTop
	to.Suffix = from.Suffix
	to.Top = from.Top
	return to
}

func (this *TableStyle) Waste(colCount uint) uint {
	if colCount <= 1 {
		colCount = 0
	} else {
		colCount--
	}
	perCol := uint(StringLength(this.Prefix) + StringLength(this.Suffix))
	return uint(StringLength(this.Left)) +
		uint(StringLength(this.Right)) +
		perCol*colCount +
		(colCount-1)*uint(StringLength(this.InnerVertical))
}

func (this *TableStyle) Render(table *Table, maxWidth uint) string {
	colWidths := this.CalculateColWidths(table, maxWidth)
	out := this.renderTopRow(colWidths)
	out += this.renderHeaderRow(table.Headers, colWidths)
	for _, row := range table.Rows {
		out += this.renderDataRow(row, colWidths)
	}
	//out += "--"
	out += this.renderBottomRow(colWidths)
	return strings.TrimRight(out, "\n") + "\n"
}

func (this *TableStyle) CalculateColWidths(table *Table, totalTableWidth uint) []uint {
	if totalTableWidth == 0 {
		totalTableWidth = DEFAULT_RENDER_WIDTH
	}
	waste := this.Waste(table.colAmount)
	if waste >= totalTableWidth {
		totalTableWidth = 0
	} else {
		totalTableWidth -= waste
	}
	colWidths := make([]uint, table.colAmount)
	sumColWidth := uint(0)
	for _, row := range table.Rows {
		widths := row.CalculateWidths(totalTableWidth)
		for idx, wd := range widths {
			if wd > colWidths[idx] {
				sumColWidth -= colWidths[idx]
				colWidths[idx] = wd
				sumColWidth += wd
			}
		}
	}
	if sumColWidth == 0 {
		sumColWidth = 1
	}
	if totalTableWidth > 0 {
		factor := float64(totalTableWidth) / float64(sumColWidth)
		usedWidth := uint(0)
		lastWidthIdx := len(colWidths) - 1
		for idx, width := range colWidths {
			colWidths[idx] = uint(float64(width) * factor)
			usedWidth += colWidths[idx]
			if idx == lastWidthIdx {
				colWidths[idx] += totalTableWidth - usedWidth
			}
		}
	}
	return colWidths
}

func (this *TableStyle) renderBorderRow(first, prefix, content, suffix, cross, last string, colWidths []uint) string {
	row := first
	lastColIdx := len(colWidths) - 1
	for idx, colWidth := range colWidths {
		row += strings.Repeat(prefix, StringLength(this.Prefix))
		row += strings.Repeat(content, int(colWidth))
		row += strings.Repeat(suffix, StringLength(this.Suffix))
		if idx < lastColIdx {
			row += cross
		}
	}
	row += last
	return row
}

func (this *TableStyle) renderContentRow(first, cross, last string, colContents []string) string {

	// transform (i in slice[string] to (i, j in slice[string][string]) in which each row i has the
	// same amount of elements j (hence normalized)
	normalizedCols := make([][]string, len(colContents))
	maxLines := 0
	for idx, colContent := range colContents {
		normalizedCols[idx] = strings.Split(colContent, "\n")
		if l := len(normalizedCols[idx]); l > maxLines {
			maxLines = l
		}
	}
	for idx, _ := range normalizedCols {
		for i := len(normalizedCols[idx]); i < maxLines; i++ {
			normalizedCols[idx] = append(normalizedCols[idx], strings.Repeat(" ", len(colContents[idx])))
		}
		//fmt.Printf("\nCOL %d: LINES = %d\n", idx, len(normalizedCols[idx]))
	}

	colAmount := len(normalizedCols) - 1
	lineAmount := len(normalizedCols[0]) - 1
	out := ""

	for ldx := 0; ldx <= lineAmount; ldx++ {
		for cdx := 0; cdx <= colAmount; cdx++ {
			if cdx == 0 {
				out += first
			}
			out += this.Prefix
			out += normalizedCols[cdx][ldx]
			out += this.Suffix
			if cdx == colAmount {
				out += last
			} else {
				out += cross
			}
			if cdx == colAmount {
				//fmt.Printf(" LINE BREAK ON L=%d, C=%d = \"%s\"\n", ldx, cdx, normalizedCols[cdx][ldx])
				out += "\n"
			}
		}
	}
	return out
}

func (this *TableStyle) renderTopRow(colWidths []uint) string {
	if this.Top != "" {
		return this.renderBorderRow(this.LeftTop, this.Top, this.Top, this.Top, this.CrossTop, this.RightTop, colWidths) + "\n"
	}
	return ""
}

func (this *TableStyle) renderHeaderRow(row *TableRow, colWidths []uint) string {
	rendered, _ := row.RenderWithWidths(colWidths)
	for idx, text := range rendered {
		rendered[idx] = this.HeaderRenderer(text)
	}
	return this.renderContentRow(this.Left, this.InnerVertical, this.Right, rendered)
}

func (this *TableStyle) renderDataRow(row *TableRow, colWidths []uint) string {
	rendered, _ := row.RenderWithWidths(colWidths)
	out := ""
	out += this.renderBorderRow(this.CrossLeft, this.InnerHorizontal, this.InnerHorizontal, this.InnerHorizontal, this.CrossInner,
		this.CrossRight, colWidths)
	for idx, text := range rendered {
		rendered[idx] = this.ContentRenderer(text)
	}
	out += "\n" + this.renderContentRow(this.Left, this.InnerVertical, this.Right, rendered)
	return out
}

func (this *TableStyle) renderBottomRow(colWidths []uint) string {
	if this.Bottom != "" {
		return this.renderBorderRow(this.LeftBottom, this.Bottom, this.Bottom, this.Bottom, this.CrossBottom, this.RightBottom, colWidths)
	}
	return ""
}
