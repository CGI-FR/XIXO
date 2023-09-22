package xixo

import (
	"io"
)

type ShellDriver struct {
	parser    *XMLParser
	processes []*Process
}

func NewShellDriver(reader io.Reader, writer io.Writer, callbacks map[string]string) ShellDriver {
	parser := NewXMLParser(reader, writer).EnableXpath()
	processes := []*Process{}

	for elementName, shell := range callbacks {
		process := NewProcess(shell)
		parser.RegisterJSONCallback(elementName, process.Callback())
		processes = append(processes, process)
	}

	return ShellDriver{parser: parser, processes: processes}
}

func (d ShellDriver) Stream() error {
	for _, process := range d.processes {
		if err := process.Start(); err != nil {
			return err
		}
	}

	err := d.parser.Stream()

	return err
}

type FuncDriver struct {
	parser *XMLParser
}

func NewFuncDriver(reader io.Reader, writer io.Writer, callbacks map[string]CallbackMap) FuncDriver {
	parser := NewXMLParser(reader, writer).EnableXpath()

	for elementName, function := range callbacks {
		parser.RegisterMapCallback(elementName, function)
	}

	return FuncDriver{parser: parser}
}

func (d FuncDriver) Stream() error {
	err := d.parser.Stream()

	return err
}
