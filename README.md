# poms

A command-line pomodoro timer to help you stay productive with focused work sessions and breaks.

## Features

- Traditional pomodoro technique with customizable work and break durations
- Two modes: simple CLI mode or interactive TUI mode
- Audio notifications when sessions complete
- Session tracking
- Clean terminal interface

## Usage

### CLI Mode (Default)
```bash
./poms -work 25 -break 5
```

### Interactive TUI Mode
```bash
./poms -tui -work 25 -break 5
```

### Options
- `-work`: Duration of work sessions in minutes (default: 25)
- `-break`: Duration of break sessions in minutes (default: 5)
- `-tui`: Start in interactive TUI mode

## Installation

```bash
go build -o poms cmd/pomodoro/main.go
```

## How it works

The tool follows the pomodoro technique:
1. Start a focused work session (default 25 minutes)
2. Take a short break (default 5 minutes)
3. Repeat the cycle

In CLI mode, it runs one work session followed by one break. In TUI mode, you get an interactive interface to manage multiple sessions.
