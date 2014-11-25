package parser

import (
	"bufio"
	"fmt"
	"github.com/Hranoprovod/shared"
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
	Nodes         chan *shared.Node
	Errors        chan *BreakingError
	Done          chan bool
}

// NewParser returns new parser
func NewParser(parserOptions *ParserOptions) *Parser {
	return &Parser{
		parserOptions,
		make(chan *shared.Node),
		make(chan *BreakingError),
		make(chan bool),
	}
}

func (p *Parser) ParseFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		p.Errors <- NewBreakingError(err.Error(), exitErrorOpeningFile)
		return
	}
	defer f.Close()
	p.ParseStream(f)
}

func (p *Parser) ParseStream(reader io.Reader) {
	var node *shared.Node
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
			node = shared.NewNode(trimmedLine)
			continue
		}

		if node != nil {
			separator := strings.LastIndexAny(trimmedLine, "\t ")

			if separator == -1 {
				p.Errors <- NewBreakingError(
					fmt.Sprintf("Bad syntax on line %d, \"%s\".", lineNumber, line),
					exitErrorBadSyntax,
				)
				return
			}
			ename := mytrim(trimmedLine[0:separator])

			//get element value
			snum := mytrim(trimmedLine[separator:])
			enum, err := strconv.ParseFloat(snum, 32)
			if err != nil {
				p.Errors <- NewBreakingError(
					fmt.Sprintf("Error converting \"%s\" to float on line %d \"%s\".", snum, lineNumber, line),
					exitErrorConversion,
				)
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
