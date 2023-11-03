package xixo

import (
	"io"
)

// Driver represents a driver that processes XML using callback functions.
type Driver struct {
	parser *XMLParser
}

// NewDriver creates a new FuncDriver instance with the given reader, writer, and callbacks.
func NewDriver(reader io.Reader, writer io.Writer, callbacks map[string]CallbackMap) Driver {
	// Create a new XML parser with XPath enabled.
	parser := NewXMLParser(reader, writer).EnableXpath()

	// Register callback functions for each element.
	for elementName, function := range callbacks {
		parser.RegisterMapCallback(elementName, function)
	}

	// Return the FuncDriver with the parser.
	return Driver{parser: parser}
}

// Stream processes the XML using registered callback functions and returns any error encountered.
func (d Driver) Stream() error {
	// Stream the XML using the parser and return any error encountered.
	err := d.parser.Stream()

	return err
}
