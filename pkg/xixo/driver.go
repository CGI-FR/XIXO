package xixo

import (
	"io"
)

// ShellDriver represents a driver that processes XML using shell commands.
type ShellDriver struct {
	parser    *XMLParser
	processes []*Process
}

// NewShellDriver creates a new ShellDriver instance with the given reader, writer, and callbacks.
func NewShellDriver(reader io.Reader, writer io.Writer, callbacks map[string]string) ShellDriver {
	// Create a new XML parser with XPath enabled and initialize processes.
	parser := NewXMLParser(reader, writer).EnableXpath()
	processes := []*Process{}

	// Iterate through the callbacks and create a process for each element.
	for elementName, shell := range callbacks {
		process := NewProcess(shell)
		// Register a JSON callback for the element and add the process to the list.
		parser.RegisterJSONCallback(elementName, process.Callback())
		processes = append(processes, process)
	}

	// Return the ShellDriver with the parser and processes.
	return ShellDriver{parser: parser, processes: processes}
}

// Stream processes the XML using the registered processes and returns any error encountered.
func (d ShellDriver) Stream() error {
	// Start each process in parallel.
	for _, process := range d.processes {
		if err := process.Start(); err != nil {
			return err
		}
	}

	// Stream the XML using the parser and return any error encountered.
	err := d.parser.Stream()

	return err
}

// FuncDriver represents a driver that processes XML using callback functions.
type FuncDriver struct {
	parser *XMLParser
}

// NewFuncDriver creates a new FuncDriver instance with the given reader, writer, and callbacks.
func NewFuncDriver(reader io.Reader, writer io.Writer, callbacks map[string]CallbackMap) FuncDriver {
	// Create a new XML parser with XPath enabled.
	parser := NewXMLParser(reader, writer).EnableXpath()

	// Register callback functions for each element.
	for elementName, function := range callbacks {
		parser.RegisterMapCallback(elementName, function)
	}

	// Return the FuncDriver with the parser.
	return FuncDriver{parser: parser}
}

// Stream processes the XML using registered callback functions and returns any error encountered.
func (d FuncDriver) Stream() error {
	// Stream the XML using the parser and return any error encountered.
	err := d.parser.Stream()

	return err
}
