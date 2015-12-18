package clif

import (
	"fmt"
)

type (
	Table struct {

		// AllowEmptyFill decides whether SetRow() and SetColumn() with row indices
		// bigger than the amount of rows creates additional, empty rows in between.
		// USE WITH CAUTION!
		AllowEmptyFill bool

		// set of output headers
		Headers *TableRow

		// colAmount is the amount of cols per row (fixed size)
		colAmount int

		// rowAmount is the amount of rows
		rowAmount int

		// set of row, col, lines
		Rows []*TableRow

		// Style for rendering table
		style *TableStyle
	}

	TableStyle struct {
		LeftTop         string
		Left            string
		LeftBottom      string
		Right           string
		RightTop        string
		RightBottom     string
		CrossInner      string
		CrossTop        string
		CrossBottom     string
		CrossLeft       string
		CrossRight      string
		InnerHorizontal string
		InnerVertical   string
		Top             string
		Bottom          string
		Prefix          string
		Suffix          string
		HeaderRenderer  func(content string) string
		ContentRenderer func(content string) string
	}

	TableRow struct {
		// MaxLineCount is the maximum amount of
		MaxLineCount int
		ColAmount    int
		Cols         []*TableCol
		table        *Table
	}

	TableCol struct {
		// Content is to the column content
		content *string

		// LineCount contains the amount of lines in the content
		lineCount int

		// renderer is the content renderer.. see `TableStyle.(Content|Header)Renderer`
		renderer func(content string) string

		// row is back-reference to row
		row *TableRow
	}
)

var (
	ErrHeadersNotSetYet = fmt.Errorf("Cannot add/set data when headers are not set")
)

// NewTable constructs new Table with optional list of headers
func NewTable(headers []string, style ...*TableStyle) *Table {
	if len(style) == 0 {
		style = []*TableStyle{NewDefaultTableStyle()}
	}
	this := &Table{
		Rows:  make([]*TableRow, 0),
		style: style[0],
	}
	if headers != nil {
		this.SetHeaders(headers)
	}
	return this
}

// Render prints the table into a string
func (this *Table) Render(maxWidth ...int) string {
	return this.style.Render(this, maxWidth...)
}

// Reset clears all (row) data of the table
func (this *Table) Reset() {
	this.Rows = make([]*TableRow, 0)
}

// SetHeaders sets the headers of the table. Can only be called before any
// data has been added.
func (this *Table) SetHeaders(headers []string) error {
	if this.Headers != nil && this.rowAmount > 0 {
		return fmt.Errorf("Cannot set headers after data has been added")
	}
	this.colAmount = len(headers)
	this.Headers = NewTableRow(headers)
	return nil
}

// SetStyle changes the table style
func (this *Table) SetStyle(style *TableStyle) {
	this.style = style
}

// AddRow adds another row to the table. Headers must be set beforehand.
func (this *Table) AddRow(cols []string) error {
	if err := this.checkAddCols(cols); err != nil {
		return err
	}
	this.addRow(cols)
	return nil
}

// AddRows adds multiple row to the table. Headers must be set beforehand.
func (this *Table) AddRows(rows [][]string) error {
	for idx, cols:= range rows {
		if err := this.checkAddCols(cols); err != nil {
			return fmt.Errorf("Row %d: %s", idx+1, err)
		}
	}
	for _, cols:= range rows {
		this.addRow(cols)
	}

	return nil
}

// SetRow sets columns in a specific row.
//
// If `AllowEmptyFill` is true, then the row index can be arbitrary and empty
// columns will be automatically created, if needed.
// Otherwise the row index must be within the bounds of existing data or an
// error is returned.
func (this *Table) SetRow(idx int, cols []string) error {
	if err := this.checkAddCols(cols); err != nil {
		return err
	}
	if idx < this.colAmount {
		row := NewTableRow(cols).SetTable(this)
		this.Rows[idx] = row
	} else if idx > this.colAmount {
		if this.AllowEmptyFill {
			empty := make([]string, this.colAmount)
			diff := int(this.colAmount - idx)
			for i := 0; i < diff; i++ {
				this.addRow(empty)
			}
			this.addRow(cols)
		} else {
			return fmt.Errorf("Cannot set row at index %d -> Only %d rows in data", idx, this.rowAmount)
		}
	} else { // == this.cols
		this.addRow(cols)
	}
	return nil
}

// SetColumn sets the contents of a specific column in a specific row. See `SetRow`
// for limitations on the row index.
// The column index must be within the bounds of the column amount.
func (this *Table) SetColumn(rowIdx, colIdx int, content string) error {
	if this.Headers == nil {
		return ErrHeadersNotSetYet
	} else if colIdx >= this.colAmount {
		return fmt.Errorf("Cannot set row at index %d -> Only %d rows in data", rowIdx, this.rowAmount)
	}
	if rowIdx < this.colAmount {
		this.Rows[rowIdx].Cols[colIdx] = NewTableCol(this.Rows[rowIdx], content)
	} else {
		cols := make([]string, this.colAmount)
		cols[colIdx] = content
		return this.SetRow(rowIdx, cols)
	}
	return nil
}

func (this *Table) addRow(cols []string) {
	row := NewTableRow(cols).SetTable(this)
	this.Rows = append(this.Rows, row)
	this.rowAmount++
}

func (this *Table) checkAddCols(cols []string) error {
	if this.Headers == nil {
		return ErrHeadersNotSetYet
	}
	if l := len(cols); l != this.colAmount {
		return fmt.Errorf("Cannot add %d cols. Expected width is %d", l, this.colAmount)
	}
	return nil
}
