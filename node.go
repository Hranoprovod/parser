package parser

// Node contains general node data
type Node struct {
	Header   string
	Elements *Elements
}

// NewNode creates new geneal node
func NewNode(header string) *Node {
	return &Node{
		header,
		NewElements(),
	}
}