package utils

import (
	"fmt"
	"strings"
)

type ColorTerminalEnum int

const (
	FgBlack ColorTerminalEnum = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

func (c ColorTerminalEnum) String() string {
	switch c {
	case FgBlack:
		return "black"
	case FgRed:
		return "red"
	case FgGreen:
		return "green"
	case FgYellow:
		return "yellow"
	case FgBlue:
		return "blue"
	case FgMagenta:
		return "magenta"
	case FgCyan:
		return "cyan"
	case FgWhite:
		return "white"
	default:
		return "unknown"
	}
}

var (
	colourStringMap = map[string]ColorTerminalEnum{
		"black":   FgBlack,
		"red":     FgRed,
		"green":   FgGreen,
		"yellow":  FgYellow,
		"blue":    FgBlue,
		"magenta": FgMagenta,
		"cyan":    FgCyan,
		"white":   FgWhite,
	}
)

func ParseStringToColorTerminal(str string) (ColorTerminalEnum, bool) {
	c, ok := colourStringMap[strings.ToLower(str)]
	return c, ok
}

func ConsolePrintColoredText(textMessage string, color ColorTerminalEnum) {
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, textMessage)
	fmt.Println(colored)
}
