package clif

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	//"time"
	"time"
)

func _testRenderProgressBar(pb ProgressBar, pos int) string {
	pb.Set(pos)
	rendered := pb.Render()
	//fmt.Printf("OUT (%d):\n---\n%s\n---\n\n", pos, rendered)
	So(StringLength(rendered), ShouldEqual, pb.RenderWidth())
	return rendered
}

func TestProgressBarRender(t *testing.T) {
	Convey("Rendering progress", t, func() {
		pb := NewProgressBar(200).SetStyle(ProgressBarStyleAscii).(*ProgressBarSimple)
		pb.SetRenderWidth(80)

		Convey("Render without info", func() {
			pb.Style().Percentage = PROGRESS_BAR_ADDON_OFF
			pb.Style().Estimate = PROGRESS_BAR_ADDON_OFF
			pb.Style().Elapsed = PROGRESS_BAR_ADDON_OFF
			pb.Style().Count = PROGRESS_BAR_ADDON_OFF

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
			pb.Style().Percentage = PROGRESS_BAR_ADDON_APPEND
			pb.Style().Estimate = PROGRESS_BAR_ADDON_APPEND
			pb.Style().Elapsed = PROGRESS_BAR_ADDON_APPEND
			pb.Style().Count = PROGRESS_BAR_ADDON_PREPEND
			pb.started = time.Now().Add(time.Minute * -2)

			rendered := _testRenderProgressBar(pb, 0)
			So(rendered, ShouldEqual, "  0/200 [--------------------------------------------] @02m00s / ~00m00s /  0.0%")

			rendered = _testRenderProgressBar(pb, 50)
			So(rendered, ShouldEqual, " 50/200 [==========>---------------------------------] @02m00s / ~06m00s / 25.0%")

			rendered = _testRenderProgressBar(pb, 100)
			So(rendered, ShouldEqual, "100/200 [=====================>----------------------] @02m00s / ~02m00s / 50.0%")

			rendered = _testRenderProgressBar(pb, 200)
			So(rendered, ShouldEqual, "200/200 [============================================] @02m00s / ~00m00s /  100%")

			Convey("Render length", func() {
				//for _, pos := range []int{0, 40, 80, 120, 160, 200} {
				//for _, pos := range []int{3, 4} {
				for pos := 0; pos <= 200; pos++ {
					_testRenderProgressBar(pb, pos)
				}
			})
		})

		Convey("Render unicode style", func() {
			pb.SetStyle(ProgressBarStyleUtf8)
			pb.Style().Percentage = PROGRESS_BAR_ADDON_OFF
			pb.Style().Estimate = PROGRESS_BAR_ADDON_OFF
			pb.Style().Elapsed = PROGRESS_BAR_ADDON_OFF
			pb.Style().Count = PROGRESS_BAR_ADDON_OFF
			//pb.started = time.Now().Add(time.Minute * -2)

			rendered := _testRenderProgressBar(pb, 0)
			So(rendered, ShouldEqual, "▕░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▏")

			rendered = _testRenderProgressBar(pb, 50)
			So(rendered, ShouldEqual, "▕██████████████████▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▏")

			rendered = _testRenderProgressBar(pb, 100)
			So(rendered, ShouldEqual, "▕██████████████████████████████████████▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▏")

			rendered = _testRenderProgressBar(pb, 200)
			So(rendered, ShouldEqual, "▕██████████████████████████████████████████████████████████████████████████████▏")

			Convey("Render unicode style with info", func() {
				pb.Style().Percentage = PROGRESS_BAR_ADDON_APPEND
				pb.Style().Estimate = PROGRESS_BAR_ADDON_APPEND
				pb.Style().Elapsed = PROGRESS_BAR_ADDON_PREPEND
				pb.Style().Count = PROGRESS_BAR_ADDON_PREPEND
				//pb.started = time.Now().Add(time.Minute * -2)

				rendered = _testRenderProgressBar(pb, 50)
				So(rendered, ShouldEqual, " 50/200 / @00m00s ▕██████████▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▏ ~00m00s / 25.0%")
			})
		})
	})
}
