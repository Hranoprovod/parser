## go-hranoprovod-parser [![Build Status](https://travis-ci.org/Hranoprovod/parser.svg)](https://travis-ci.org/Hranoprovod/parser) [![GoDoc](https://godoc.org/github.com/Hranoprovod/parser?status.svg)](https://godoc.org/github.com/Hranoprovod/parser)

Go hranoprovod file parser

Sample usage:

```go
prsr := parser.NewParser(parser.NewDefaultOptions())
go prsr.ParseFile("data.yaml")
err := func() error {
	for {
		select {
		case node := <-prsr.Nodes:
			print(node.Header)
		case err := <-prsr.Errors:
			return err
		case <-prsr.Done:
			return nil
		}
	}
}()

```
