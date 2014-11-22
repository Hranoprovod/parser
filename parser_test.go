package parser

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

type NodeList map[string]*Node

func NewNodeList() *NodeList {
	return &NodeList{}
}

func (db *NodeList) push(node *Node) {
	(*db)[(*node).Header] = node
}

func readChannels(parser *Parser) (*NodeList, error) {
	nodeList := NewNodeList()
	for {
		select {
		case node := <-parser.Nodes:
			nodeList.push(node)
		case breakingError := <-parser.Errors:
			return nil, breakingError
		case <-parser.Done:
			return nodeList, nil
		}
	}
}

func TestParser(t *testing.T) {
	Convey("Given new parser", t, func() {
		parser := NewParser(NewDefaultParserOptions())
		Convey("It completes successfully on empty string", func() {
			go parser.ParseStream(strings.NewReader(""))
			nodeList, error := readChannels(parser)
			So(len(*nodeList), ShouldEqual, 0)
			So(error, ShouldBeNil)
		})

		Convey("It processes valid node", func() {
			file := `2011/07/17:
  el1: 1.22
  ел 2:  4
  el/3:  3

2011/07/18:
  el1: 1.33
  ел 5:  5
  el/7:  4
  el1: 1.35
  `
			go parser.ParseStream(strings.NewReader(file))
			nodeList, err := readChannels(parser)
			So(len(*nodeList), ShouldEqual, 2)
			So(err, ShouldBeNil)
			node := (*nodeList)["2011/07/17"]
			So(node.Header, ShouldEqual, "2011/07/17")
			elements := node.Elements
			So(elements, ShouldNotBeNil)
			So(len(*elements), ShouldEqual, 3)
			So((*elements)[0].Name, ShouldEqual, "el1")
			So((*elements)[0].Val, ShouldEqual, 1.22)
			So((*elements)[1].Name, ShouldEqual, "ел 2")
			So((*elements)[1].Val, ShouldEqual, 4.0)
			So((*elements)[2].Name, ShouldEqual, "el/3")
			So((*elements)[2].Val, ShouldEqual, 3.0)
		})

		Convey("It raises bad syntax error", func() {
			file := `asdasd
  asdasd2`
			go parser.ParseStream(strings.NewReader(file))
			_, err := readChannels(parser)
			So(err, ShouldNotBeNil)
			bsError, ok := err.(*ParserErrorBadSyntax)
			So(ok, ShouldBeTrue)
			So(err.Error(), ShouldEqual, "Bad syntax on line 2, \"  asdasd2\".")
			So(bsError.LineNumber, ShouldEqual, 2)
			So(bsError.Line, ShouldEqual, "  asdasd2")
		})

		Convey("It raises conversion error", func() {
			file := `asdasd
  asdasd2 s`
			go parser.ParseStream(strings.NewReader(file))
			_, err := readChannels(parser)
			So(err, ShouldNotBeNil)
			cErr, ok := err.(*ParserErrorConversion)
			So(ok, ShouldBeTrue)
			So(err.Error(), ShouldEqual, "Error converting \"s\" to float on line 2 \"  asdasd2 s\".")
			So(cErr.LineNumber, ShouldEqual, 2)
			So(cErr.Text, ShouldEqual, "s")
			So(cErr.Line, ShouldEqual, "  asdasd2 s")
		})
	})
}
