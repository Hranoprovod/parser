package parser

import (
	"fmt"
)

type ParserErrorIO struct {
	err      error
	FileName string
}

func NewParserErrorIO(err error, fileName string) *ParserErrorIO {
	return &ParserErrorIO{ err, fileName}
}

// Error returns the error message
func (e *ParserErrorIO) Error() string {
	return e.err.Error()
}

type ParserErrorBadSyntax struct {
	LineNumber int 
	Line string
}

func NewParserErrorBadSyntax(lineNumber int, line string) *ParserErrorBadSyntax {
	return &ParserErrorBadSyntax{lineNumber, line}
}

func (e *ParserErrorBadSyntax) Error() string {
	return fmt.Sprintf("Bad syntax on line %d, \"%s\".", e.LineNumber, e.Line)
}

type ParserErrorConversion struct {
	Text string
	LineNumber int
	Line string
}

func NewParserErrorConversion(text string, lineNumber int, line string) *ParserErrorConversion {
	return &ParserErrorConversion{text, lineNumber, line}
}

func (e *ParserErrorConversion) Error() string {
	return 	fmt.Sprintf("Error converting \"%s\" to float on line %d \"%s\".", e.Text, e.LineNumber, e.Line)

}
