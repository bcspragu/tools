package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type state int

func (s state) Type() string {
	switch s {
	case Dark:
		return "Dark"
	case Light:
		return "Light"
	}
	return ""
}

const (
	Unknown state = iota
	Dark
	Light
)

var home string

var startRE = regexp.MustCompile(`! --- Start \w+ ---`)
var endRE = regexp.MustCompile(`! --- End \w+ ---`)
var vimRE = regexp.MustCompile(`set background=(\w+)`)

func main() {
	home = os.Getenv("HOME")

	fX, err := os.OpenFile(home+"/.Xresources", os.O_RDWR, 0640)
	if err != nil {
		panic(err)
	}
	defer fX.Close()

	err = switchResources(fX)
	if err != nil {
		panic(err)
	}

	// Shut off vim switching
	/*
	 *
	 * fV, err := os.OpenFile(home+"/.vimrc", os.O_RDWR, 0640)
	 * if err != nil {
	 * 	panic(err)
	 * }
	 * defer fV.Close()
	 *
	 * err = switchVimrc(fV)
	 * if err != nil {
	 * 	panic(err)
	 * }
	 */

	err = exec.Command("xrdb", "-merge", home+"/.Xresources").Run()
	if err != nil {
		panic(err)
	}
}

func switchResources(fX *os.File) error {
	flip := false
	var out bytes.Buffer
	scanner := bufio.NewScanner(fX)

	for scanner.Scan() {
		if endRE.MatchString(scanner.Text()) {
			flip = false
		}

		line := scanner.Text()
		if flip {
			if strings.HasPrefix(line, "! ") {
				line = line[2:]
			} else {
				line = "! " + line
			}
		}

		if startRE.MatchString(scanner.Text()) {
			flip = true
		}
		out.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	d := out.Bytes()
	// Remove the last newline
	d = d[:len(d)-2]
	fX.Seek(0, 0)
	n, err := fX.Write(d)
	if err == nil && n < len(d) {
		return io.ErrShortWrite
	}

	return err
}

func switchVimrc(fV *os.File) error {
	var out bytes.Buffer
	scanner := bufio.NewScanner(fV)

	for scanner.Scan() {
		line := scanner.Text()
		if vimRE.MatchString(line) {
			color := vimRE.FindStringSubmatch(line)[1]
			switch color {
			case "dark":
				line = "set background=light"
			case "light":
				line = "set background=dark"
			}
		}

		out.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	d := out.Bytes()
	// Remove the last newline
	d = d[:len(d)-2]
	fV.Seek(0, 0)
	n, err := fV.Write(d)
	if err == nil && n < len(d) {
		return io.ErrShortWrite
	}

	return err
	return nil
}
