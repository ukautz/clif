
## v1 (2015-12)

* Allow options in the form `--flag` or `--name=value` before command name
* Fix: input `Ask()` failed to handle `io.EOF`, which can happen when input is not `io.Stdin`, but eg file or network
* Added `c.NewFlag(..)` to command, because `c.AddOption(clif.NewOption(...).IsFlag())` is too tedious
* Added table rendering, via `output.Table(header []string) *Table`


## v0

* Initial release
