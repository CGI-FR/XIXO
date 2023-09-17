package xixo

import (
	"io"
)

type Driver struct {
	parser    *XMLParser
	processes []*Process
}

func NewDriver(reader io.Reader, writer io.Writer, callbacks map[string]string) Driver {
	parser := NewXMLParser(reader, writer)
	processes := []*Process{}

	for elementName, shell := range callbacks {
		process := NewProcess(shell)
		parser.RegisterJSONCallback(elementName, process.Callback())
		processes = append(processes, process)
	}

	return Driver{parser: parser, processes: processes}
}

func (d Driver) Stream() error {
	for _, process := range d.processes {
		if err := process.Start(); err != nil {
			return err
		}
	}

	err := d.parser.Stream()

	return err
}
