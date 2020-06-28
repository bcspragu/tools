package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type state interface {
	Advance(txt string) (state, error)
}

type stateStart struct{}

func (s stateStart) Advance(txt string) (state, error) {
	if txt == "! Dark" || txt == "! Light" {
		return stateColorTitle{}, nil
	}
	return s, nil
}

type stateColorTitle struct{}

func (s stateColorTitle) Advance(txt string) (state, error) {
	if !strings.Contains(txt, "#define ") {
		return nil, errors.New("expected color definitions")
	}
	return stateColorBody{}, nil
}

type stateColorBody struct{}

func (s stateColorBody) Advance(txt string) (state, error) {
	if !strings.Contains(txt, "#define ") {
		return stateStart{}, nil
	}
	return s, nil
}

func main() {
	xresources := filepath.Join(os.Getenv("HOME"), ".Xresources")

	fx, err := os.OpenFile(xresources, os.O_RDWR, 0640)
	if err != nil {
		panic(err)
	}
	defer fx.Close()

	err = switchResources(fx)
	if err != nil {
		panic(err)
	}

	err = exec.Command("xrdb", "-merge", xresources).Run()
	if err != nil {
		panic(err)
	}
}

func switchResources(fx *os.File) error {
	var out bytes.Buffer
	scanner := bufio.NewScanner(fx)

	var s state = stateStart{}
	var err error
	for scanner.Scan() {
		txt := scanner.Text()
		s, err = s.Advance(txt)
		if err != nil {
			return fmt.Errorf("invalid state transition: %v", err)
		}
		if _, ok := s.(stateColorBody); ok {
			txt = flip(txt)
		}
		out.WriteString(txt + "\n")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	d := bytes.TrimSuffix(out.Bytes(), []byte("\n"))

	fx.Seek(0, 0)
	n, err := fx.Write(d)
	if err == nil && n < len(d) {
		return io.ErrShortWrite
	}

	return err
}

func flip(txt string) string {
	if strings.HasPrefix(txt, "! ") {
		return txt[2:]
	}
	return "! " + txt
}
