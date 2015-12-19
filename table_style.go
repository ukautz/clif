package clif

import (
	"fmt"
	"strings"
)

var (
	ClosedTableStyle = &TableStyle{
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
	ClosedTableStyleLight = &TableStyle{
		Bottom:          "\033[38;5;234m─\033[0m",
		ContentRenderer: func(content string) string { return content },
		CrossBottom:     "\033[38;5;234m┴\033[0m",
		CrossInner:      "\033[38;5;234m┼\033[0m",
		CrossLeft:       "\033[38;5;234m├\033[0m",
		CrossRight:      "\033[38;5;234m┤\033[0m",
		CrossTop:        "\033[38;5;234m┬\033[0m",
		HeaderRenderer:  func(content string) string { return fmt.Sprintf("\033[1;4m%s\033[0m", content) },
		InnerHorizontal: "\033[38;5;234m─\033[0m",
		InnerVertical:   "\033[38;5;234m│\033[0m",
		Left:            "\033[38;5;234m│\033[0m",
		LeftBottom:      "\033[38;5;234m└\033[0m",
		LeftTop:         "\033[38;5;234m┌\033[0m",
		Prefix:          " ",
		Right:           "\033[38;5;234m│\033[0m",
		RightBottom:     "\033[38;5;234m┘\033[0m",
		RightTop:        "\033[38;5;234m┐\033[0m",
		Suffix:          " ",
		Top:             "\033[38;5;234m─\033[0m",
	}
	OpenTableStyle = &TableStyle{
		Bottom:          "",
		ContentRenderer: func(content string) string { return content },
		CrossBottom:     "┴",
		CrossInner:      "┼",
		CrossLeft:       "",
		CrossRight:      "",
		CrossTop:        "┬",
		HeaderRenderer:  func(content string) string { return fmt.Sprintf("\033[1m%s\033[0m", content) },
		InnerHorizontal: "─",
		InnerVertical:   "│",
		Left:            "",
		LeftBottom:      "",
		LeftTop:         "",
		Prefix:          " ",
		Right:           "",
		RightBottom:     "",
		RightTop:        "",
		Suffix:          " ",
		Top:             "",
	}
	OpenTableStyleLight = &TableStyle{
		Bottom:          "",
		ContentRenderer: func(content string) string { return content },
		CrossBottom:     "\033[38;5;234m┴\033[0m",
		CrossInner:      "\033[38;5;234m┼\033[0m",
		CrossLeft:       "",
		CrossRight:      "",
		CrossTop:        "\033[38;5;234m┬\033[0m",
		HeaderRenderer:  func(content string) string { return fmt.Sprintf("\033[1m%s\033[0m", content) },
		InnerHorizontal: "\033[38;5;234m─\033[0m",
		InnerVertical:   "\033[38;5;234m│\033[0m",
		Left:            "",
		LeftBottom:      "",
		LeftTop:         "",
		Prefix:          " ",
		Right:           "",
		RightBottom:     "",
		RightTop:        "",
		Suffix:          " ",
		Top:             "",
	}
	DefaultTableStyle *TableStyle
)

func NewDefaultTableStyle() *TableStyle {
	return CopyTableStyle(DefaultTableStyle)
}

func CopyTableStyle(from *TableStyle) *TableStyle {
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

// Waste returns the amount of "wasted" characters (table render characters
// + prefix and suffix whitespaces) for given amount of columns
func (this *TableStyle) Waste(colCount int) int {
	if colCount <= 1 {
		colCount = 0
	} else {
		colCount--
	}
	perCol := StringLength(this.Prefix) + StringLength(this.Suffix)
	return StringLength(this.Left) +
		StringLength(this.Right) +
		perCol*colCount +
		(colCount-1)*StringLength(this.InnerVertical)
}

// Render renders given table with a given max width. Max width can not be lower
// than waste per row (see `Waste()`) and at least one character per column.
func (this *TableStyle) Render(table *Table, mw ...int) string {
	//fmt.Printf("\n-----------------------\n\n+ RENDER START\n")
	maxWidth := 0
	if len(mw) > 0 {
		maxWidth = mw[0]
	}
	colWidths := this.CalculateColWidths(table, maxWidth)
	out := this.renderTopRow(colWidths)
	out += this.renderHeaderRow(table.Headers, colWidths)
	for _, row := range table.Rows {
		out += this.renderDataRow(row, colWidths)
	}
	//out += "--"
	out += this.renderBottomRow(colWidths)
	//fmt.Printf("\n+ RENDER END\n")
	return strings.TrimRight(out, "\n") + "\n"
}

// CalculateColWidths returns the widths of the cols of the table, for given max width
func (this *TableStyle) CalculateColWidths(table *Table, totalTableWidth int) []int {
	if totalTableWidth == 0 {
		totalTableWidth = TermWidthCurrent
	}
	waste := this.Waste(table.colAmount)
	//fmt.Printf("\n+ CALC COL WIDTHS (MAX = %d, WASTE = %d)\n", totalTableWidth, waste)
	if waste >= totalTableWidth {
		totalTableWidth = 0
	} else {
		totalTableWidth -= waste
	}
	colWidths := make([]int, table.colAmount)
	sumColWidth := 0
	table.Headers.SetRenderer(this.HeaderRenderer)
	rows := make([]*TableRow, table.rowAmount+1)
	rows[0] = table.Headers
	for idx, row := range table.Rows {
		row.SetRenderer(this.ContentRenderer)
		rows[idx+1] = row
	}
	for _, row := range rows {
		widths := row.CalculateWidths(totalTableWidth)
		//fmt.Printf("  ## ROW %d -> %v\n", idx, widths)
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
		usedWidth := 0
		lastWidthIdx := len(colWidths) - 1
		for idx, width := range colWidths {
			colWidths[idx] = int(float64(width) * factor)
			usedWidth += colWidths[idx]
			if idx == lastWidthIdx {
				colWidths[idx] += totalTableWidth - usedWidth
			}
		}
	}
	//fmt.Printf("\n+ CALC COL WIDTHS DONE SUM=%d, (=%v)\n", sumColWidth, colWidths)
	return colWidths
}

func (this *TableStyle) renderBorderRow(first, prefix, content, suffix, cross, last string, colWidths []int) string {
	row := first
	lastColIdx := len(colWidths) - 1
	for idx, colWidth := range colWidths {
		row += strings.Repeat(prefix, StringLength(this.Prefix))
		if cw := colWidth; cw > 0 {
			row += strings.Repeat(content, cw)
		}
		row += strings.Repeat(suffix, StringLength(this.Suffix))
		if idx < lastColIdx {
			row += cross
		}
	}
	row += last
	return row
}

func (this *TableStyle) renderContentRow(first, cross, last string, colContents []string, colWidths []int) string {

	// transform (i in slice[string] to (i, j in slice[string][string]) in which each row i has the
	// same amount of elements j (hence normalized)
	colCount := len(colContents)
	normalizedCols := make([][]string, colCount)
	emptyCol := make([]string, colCount)
	maxLines := 0
	for idx, colContent := range colContents {
		normalizedCols[idx] = strings.Split(colContent, "\n")
		if l := len(normalizedCols[idx]); l > maxLines {
			maxLines = l
		}
		emptyCol[idx] = strings.Repeat(" ", int(colWidths[idx]))
	}
	for idx, _ := range normalizedCols {
		for i := len(normalizedCols[idx]); i < maxLines; i++ {
			normalizedCols[idx] = append(normalizedCols[idx], emptyCol[idx])
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

func (this *TableStyle) renderTopRow(colWidths []int) string {
	if this.Top != "" {
		return this.renderBorderRow(this.LeftTop, this.Top, this.Top, this.Top, this.CrossTop, this.RightTop, colWidths) + "\n"
	}
	return ""
}

func (this *TableStyle) renderHeaderRow(row *TableRow, colWidths []int) string {
	row.SetRenderer(this.HeaderRenderer)
	rendered, _ := row.RenderWithWidths(colWidths)
	/*for idx, text := range rendered {
		rendered[idx] = this.HeaderRenderer(text)
	}*/
	return this.renderContentRow(this.Left, this.InnerVertical, this.Right, rendered, colWidths)
}

func (this *TableStyle) renderDataRow(row *TableRow, colWidths []int) string {
	row.SetRenderer(this.ContentRenderer)
	rendered, _ := row.RenderWithWidths(colWidths)
	out := ""
	out += this.renderBorderRow(this.CrossLeft, this.InnerHorizontal, this.InnerHorizontal, this.InnerHorizontal, this.CrossInner,
		this.CrossRight, colWidths)
	/*for idx, text := range rendered {
		rendered[idx] = this.ContentRenderer(text)
	}*/
	out += "\n" + this.renderContentRow(this.Left, this.InnerVertical, this.Right, rendered, colWidths)
	return out
}

func (this *TableStyle) renderBottomRow(colWidths []int) string {
	if this.Bottom != "" {
		return this.renderBorderRow(this.LeftBottom, this.Bottom, this.Bottom, this.Bottom, this.CrossBottom, this.RightBottom, colWidths)
	}
	return ""
}

func init() {
	DefaultTableStyle = CopyTableStyle(ClosedTableStyle)
}
