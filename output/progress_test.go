package output

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
	"unicode/utf8"
)

func _testRenderProgressBar(pb *ProgressBar, pos int) string {
	pb.Set(pos)
	rendered := pb.Render()
	//fmt.Printf("OUT (%d):\n---\n%s\n---\n\n", pos, rendered)
	So(StringLength(rendered), ShouldEqual, pb.renderWidth)
	return rendered
}

func TestProgressBarRender(t *testing.T) {
	Convey("Rendering progress", t, func() {
		pb := NewProgressBar(200).SetStyle(PbStyleAscii)

		Convey("Render without info", func() {
			pb.RenderProgressCount = false
			pb.RenderProgressPercentage = false
			pb.RenderTimeEstimate = false
			pb.started = time.Now().Add(time.Minute * -2)

			rendered := _testRenderProgressBar(pb, 0)
			So(rendered, ShouldEqual, "[------------------------------------------------------------------------------]")

			rendered = _testRenderProgressBar(pb, 50)
			So(rendered, ShouldEqual, "[==================>-----------------------------------------------------------]")

			rendered = _testRenderProgressBar(pb, 100)
			So(rendered, ShouldEqual, "[======================================>---------------------------------------]")

			rendered = _testRenderProgressBar(pb, 200)
			So(rendered, ShouldEqual, "[==============================================================================]")

			Convey("Render length", func() {
				//for _, pos := range []int{0, 3, 4, 6, 200} {
				//for _, pos := range []int{0, 40, 80, 120, 160, 200} {
				for pos := 0; pos <= 200; pos++ {

					//fmt.Printf("\n\n*******************************\n")
					_testRenderProgressBar(pb, pos)
					//rendered := _testRenderProgressBar(pb, pos)
					//fmt.Printf("OUT (%d):\n---\n%s\n---\n\n", pos, rendered)
				}
			})
		})

		Convey("Render with info", func() {
			pb.RenderProgressCount = true
			pb.RenderProgressPercentage = true
			pb.RenderTimeEstimate = true
			pb.started = time.Now().Add(time.Minute * -2)

			rendered := _testRenderProgressBar(pb, 0)
			So(rendered, ShouldEqual, "[  0/200] [---------------------------------------------------] [ 0.0% - 00m00s]")

			rendered = _testRenderProgressBar(pb, 50)
			So(rendered, ShouldEqual, "[ 50/200] [===========>---------------------------------------] [25.0% - 06m00s]")

			rendered = _testRenderProgressBar(pb, 100)
			So(rendered, ShouldEqual, "[100/200] [========================>--------------------------] [50.0% - 02m00s]")

			rendered = _testRenderProgressBar(pb, 200)
			So(rendered, ShouldEqual, "[200/200] [===================================================] [100.% - 00m00s]")

			Convey("Render length", func() {
				//for _, pos := range []int{0, 40, 80, 120, 160, 200} {
				//for _, pos := range []int{3, 4} {
				for pos := 0; pos <= 200; pos++ {
					_testRenderProgressBar(pb, pos)
				}
			})
		})

		Convey("Render unicode style", func() {
			pb.RenderProgressCount = false
			pb.RenderProgressPercentage = false
			pb.RenderTimeEstimate = false
			pb.SetStyle(PbStyleUtf8)
			pb.started = time.Now().Add(time.Minute * -2)

			rendered := _testRenderProgressBar(pb, 0)
			So(rendered, ShouldEqual, "▕░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▏")

			rendered = _testRenderProgressBar(pb, 50)
			So(rendered, ShouldEqual, "▕██████████████████▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▏")

			rendered = _testRenderProgressBar(pb, 100)
			So(rendered, ShouldEqual, "▕██████████████████████████████████████▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▏")

			rendered = _testRenderProgressBar(pb, 200)
			So(rendered, ShouldEqual, "▕██████████████████████████████████████████████████████████████████████████████▏")

			Convey("Render unicode style with info", func() {
				pb.RenderProgressCount = true
				pb.RenderProgressPercentage = true
				pb.RenderTimeEstimate = true
				pb.started = time.Now().Add(time.Minute * -2)

				rendered = _testRenderProgressBar(pb, 50)
				So(rendered, ShouldEqual, "[ 50/200] ▕███████████▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▏ [25.0% - 06m00s]")
			})
		})
	})
}
