package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
)

var (
	vendor = flag.String("vendor", "intel_backlight", "The backlight vendor to modify")
	amount = flag.Int("amount", 1, "Percentage to increase, decrease, or set the brightness")
	inc    = flag.Bool("inc", false, "Increment the brightness")
	dec    = flag.Bool("dec", false, "Decrement the brightness")
	set    = flag.Bool("set", false, "Set the brightness")
)

type backlight struct {
	min     int
	max     int
	current int
}

func main() {
	flag.Parse()
	if !*inc && !*dec && !*set {
		fmt.Println("No option set")
		return
	}

	b, err := makeBacklight()
	if err != nil {
		fmt.Println(err)
		return
	}

	var modFunc func(int) error

	if *inc {
		modFunc = b.inc
	}

	if *dec {
		modFunc = b.dec
	}

	if *set {
		modFunc = b.set
	}

	if err := modFunc(*amount); err != nil {
		fmt.Println(err)
	}
	return
}

func makeBacklight() (*backlight, error) {
	bl := &backlight{}

	if b, err := getBrightness("max_brightness"); err == nil {
		bl.max = b
	} else {
		return nil, err
	}

	if b, err := getBrightness("brightness"); err == nil {
		bl.current = b
	} else {
		return nil, err
	}

	return bl, nil
}

func getBrightness(filename string) (int, error) {
	dat, err := ioutil.ReadFile(getDirectory() + filename)
	if err != nil {
		return 0, err
	}

	s := string(bytes.TrimSpace(dat))
	b, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return b, nil
}

func (b *backlight) inc(i int) error {
	return b.setBrightness(b.current + int(b.max*i/100.0))
}

func (b *backlight) dec(i int) error {
	return b.setBrightness(b.current - int(b.max*i/100.0))
}

func (b *backlight) set(i int) error {
	return b.setBrightness(int(b.max * i / 100.0))
}

func (b *backlight) setBrightness(i int) error {
	i = min(max(b.min, i), b.max)
	d := []byte(strconv.Itoa(i) + "\n")
	return ioutil.WriteFile(getDirectory()+"brightness", d, 0644)
}

func getDirectory() string {
	return "/sys/class/backlight/" + *vendor + "/"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
