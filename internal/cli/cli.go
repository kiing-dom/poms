package cli

import "flag"

var workMinutes = flag.Int("work", 25, "duration of the work session")
var breakMinutes = flag.Int("break", 5, "duration of the break")

func ParseFlags() (int, int) {
	flag.Parse()
	return *workMinutes, *breakMinutes
}

