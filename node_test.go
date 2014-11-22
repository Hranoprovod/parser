package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNodeList(t *testing.T) {
	Convey("Given Node", t, func() {
		node := NewNode("test")
		Convey("Is not null", func() {
			So(node, ShouldNotBeNil)
		})
	})
}
