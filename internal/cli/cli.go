package cli

import "flag"

var (
	workMinutes  = flag.Int("work", 25, "duration of the work session")
	breakMinutes = flag.Int("break", 5, "duration of the break")
	interactive  = flag.Bool("tui", false, "start interactive TUI mode")
)

type Config struct {
	WorkMinutes int
	BreakMinutes int
	Interactive bool
}

func ParseFlags() Config {
	flag.Parse()
	return Config {
		WorkMinutes: *workMinutes,
		BreakMinutes: *breakMinutes,
		Interactive: *interactive,
	}
}

func (c Config) ShouldUseTUI() bool {
	return c.Interactive
}
