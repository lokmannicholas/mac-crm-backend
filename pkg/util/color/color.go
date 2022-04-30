package color

import (
	"fmt"
	"strconv"
	"strings"
)

type Color string

const Reset Color = "\033[0m"
const Red Color = "\033[31m"
const Green Color = "\033[32m"
const Yellow Color = "\033[33m"
const Blue Color = "\033[34m"
const Purple Color = "\033[35m"
const Cyan Color = "\033[36m"
const Gray Color = "\033[37m"
const White Color = "\033[97m"

func Colorize(color Color, content string) string {
	return fmt.Sprintf("%s%s%s", (color), content, Reset)
}

type Hex string

type RGB struct {
	Red   int
	Green int
	Blue  int
}

func (h Hex) toRGB() (RGB, error) {
	return Hex2RGB(h)
}

func Hex2RGB(hex Hex) (RGB, error) {
	h := strings.Replace(string(hex), "#", "", 1)
	var rgb RGB
	values, err := strconv.ParseUint(h, 16, 32)

	if err != nil {
		return RGB{}, err
	}

	rgb = RGB{
		Red:   int(values >> 16),
		Green: int((values >> 8) & 0xFF),
		Blue:  int(values & 0xFF),
	}

	return rgb, nil
}
