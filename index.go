package carp

// IndexDescriptor
type IndexDescriptor struct {
	Name   string   // index name
	Fields []string // field columns
	Index  IndexDef
}
