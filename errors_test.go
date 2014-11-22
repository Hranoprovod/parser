package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBreakingError(t *testing.T) {
	Convey("Given new IO error", t, func() {
		err := NewParserErrorIO(nil, "file_name")
		Convey("Error is of the right type", func() {
			So(err.FileName, ShouldEqual, "file_name")
		})
	})
}
