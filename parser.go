package parser

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	runeTab   = '\t'
	runeSpace = ' '
)

// ParserOptions contains the parser related options
type ParserOptions struct {
	commentChar uint8
}

// NewDefaultParserOptions returns the default set of parser options
func NewDefaultParserOptions() *ParserOptions {
	return &ParserOptions{'#'}
}

// Parser is the parser data structure
type Parser struct {
	parserOptions *ParserOptions
	Nodes         chan *Node
	Errors        chan error
	Done          chan bool
}

// NewParser returns new parser
func NewParser(parserOptions *ParserOptions) *Parser {
	return &Parser{
		parserOptions,
		make(chan *Node),
		make(chan error),
		make(chan bool),
	}
}

func (p *Parser) ParseFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		p.Errors <- NewParserErrorIO(err, fileName)
		return
	}
	defer f.Close()
	p.ParseStream(f)
}

func (p *Parser) ParseStream(reader io.Reader) {
	var node *Node
	lineNumber := 0
	lineScanner := bufio.NewScanner(reader)
	for lineScanner.Scan() {
		lineNumber++
		line := lineScanner.Text()
		trimmedLine := mytrim(line)

		//skip empty lines and lines starting with #
		if trimmedLine == "" || line[0] == p.parserOptions.commentChar {
			continue
		}

		//new nodes start at the beginning of the line
		if line[0] != runeSpace && line[0] != runeTab {
			if node != nil {
				p.Nodes <- node
			}
			node = NewNode(trimmedLine)
			continue
		}

		if node != nil {
			separator := strings.LastIndexAny(trimmedLine, "\t ")

			if separator == -1 {
				p.Errors <- NewParserErrorBadSyntax(lineNumber, line)
				return
			}
			ename := mytrim(trimmedLine[0:separator])

			//get element value
			snum := mytrim(trimmedLine[separator:])
			enum, err := strconv.ParseFloat(snum, 32)
			if err != nil {
				p.Errors <- NewParserErrorConversion(snum, lineNumber, line)
				return
			}

			if ndx, exists := node.Elements.Index(ename); exists {
				(*node.Elements)[ndx].Val += float32(enum)
			} else {
				node.Elements.Add(ename, float32(enum))
			}
		}
	}
	// push last node
	if node != nil {
		p.Nodes <- node
	}
	p.Done <- true
}
