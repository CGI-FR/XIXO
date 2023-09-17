package xixo

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type Process struct {
	command string
	cmd     *exec.Cmd
	stdout  io.ReadCloser
	stdin   io.WriteCloser
	scanner *bufio.Scanner
}

func NewProcess(command string) *Process {
	return &Process{command: command}
}

func (p *Process) Start() error {
	var err error

	//nolint: gosec
	p.cmd = exec.Command("/bin/sh", "-c", p.command)

	p.stdout, err = p.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	p.stdin, err = p.cmd.StdinPipe()
	if err != nil {
		return err
	}

	p.scanner = bufio.NewScanner(p.stdout)

	go func(p *Process) {
		//nolint
		p.cmd.Run()
	}(p)

	return nil
}

func (p *Process) Stop() error {
	errStdin := p.stdin.Close()

	if errStdin != nil {
		return fmt.Errorf("can't close pipe with process: %w", errStdin)
	}

	return nil
}

func (p *Process) Callback() CallbackJSON {
	return func(s string) (string, error) {
		log.Debug().Str("json", s).Msg("request edit json")

		_, err := p.stdin.Write([]byte(s + "\n"))
		if err != nil {
			return "", err
		}

		if !p.scanner.Scan() {
			return "", fmt.Errorf("command doesn't return line")
		}

		log.Debug().Str("read", p.scanner.Text()).Msg("reading from process")

		return p.scanner.Text(), nil
	}
}
