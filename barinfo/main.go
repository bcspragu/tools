package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var outRE = regexp.MustCompile(`\s*((?:\w+\s?)+):\s*(\w[\w\.\s%\[\]-]+)`)
var volRE = regexp.MustCompile(`\w+ \w+ \[(\w+%)\] \[[\w-\.]+\] \[(\w+)\]`)

func main() {
	ticker := time.Tick(5*time.Second)
	for range ticker {
		printInfo()
	}
}

func printInfo() {
	pow, err := exec.Command("upower", "-i", "/org/freedesktop/UPower/devices/battery_BAT0").Output()
	if err != nil {
		fmt.Println("QUnknown")
		return
	}
	vol, err := exec.Command("amixer", "-c", "0", "get", "Master").Output()
	if err != nil {
		fmt.Println("QUnknown")
		return
	}
	fmt.Printf("Q%s | %s\n", parsePow(string(pow)), parseVol(string(vol)))
}

func parsePow(in string) string {
	lines := strings.Split(in, "\n")
	var percentage, state, time string
	for _, line := range lines {
		match := outRE.FindStringSubmatch(line)
		if len(match) == 0 {
			continue
		}
		key, value := match[1], match[2]
		switch key {
		case "percentage":
			percentage = value
		case "time to empty":
			time = value
		case "time to full":
			time = value
		case "state":
			state = value
		}
	}
	return fmt.Sprintf("Battery: %s - %s - %s", strings.TrimSpace(percentage), strings.TrimSpace(state), parseBatteryLife(time))
}

func parseVol(in string) string {
	lines := strings.Split(in, "\n")
	var vol, mute string
	for _, line := range lines {
		match := outRE.FindStringSubmatch(line)
		if len(match) == 0 {
			continue
		}
		key, value := match[1], match[2]
		switch key {
		case "Mono":
			match := volRE.FindStringSubmatch(value)
			vol = match[1]
			if match[2] == "off" {
				mute = "Muted"
			}
		}
	}
	if len(mute) > 0 {
		return fmt.Sprintf("Volume: %s %s", vol, mute)
	} else {
		return fmt.Sprintf("Volume: %s", vol)
	}
}

func parseBatteryLife(time string) string {
	time = strings.TrimSpace(time)
	timeParts := strings.Split(time, " ")
	if len(timeParts) < 2 {
		return "Unknown"
	}
	t, unit := timeParts[0], timeParts[1]
	val, err := strconv.ParseFloat(t, 64)
	if err != nil {
		return "Unknown"
	}
	switch unit {
	case "hours":
		if val > 1 {
			return fmt.Sprintf("%d hours %d minutes", int(val), int(val*60)%60)
		} else {
			return fmt.Sprintf("%d minutes", int(val*60)%60)
		}
	case "minutes":
		return fmt.Sprintf("%d minutes", int(val))
	}
	return "Unknown"
}
