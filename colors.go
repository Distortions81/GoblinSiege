package main

import (
	"image/color"
	"strconv"
	"strings"
)

func Hex2Color(input string) (color.RGBA, error) {
	input = strings.TrimPrefix(input, "#")

	var rgb color.RGBA
	values, err := strconv.ParseUint(string(input), 16, 32)

	if err != nil {
		return color.RGBA{}, err
	}

	rgb = color.RGBA{
		R: uint8(values >> 16),
		G: uint8((values >> 8) & 0xFF),
		B: uint8(values & 0xFF),
	}

	return rgb, nil
}

var (
	ColorVeryDarkGray   = color.NRGBA{64, 64, 64, 255}
	ColorReallyDarkGray = color.NRGBA{32, 32, 32, 255}
	ColorSmoke          = color.NRGBA{32, 32, 32, 200}
	ColorVeryDarkGreen  = color.NRGBA{11, 73, 12, 255}
	ColorVeryDarkRed    = color.NRGBA{73, 11, 11, 255}

	ColorGreenC = color.NRGBA{40, 180, 99, 64}

	ColorBlack   = color.NRGBA{0, 0, 0, 255}
	ColorRed     = color.NRGBA{203, 67, 53, 255}
	ColorGreen   = color.NRGBA{40, 180, 99, 255}
	ColorBlue    = color.NRGBA{41, 128, 185, 255}
	ColorYellow  = color.NRGBA{244, 208, 63, 255}
	ColorGray    = color.NRGBA{128, 128, 128, 255}
	ColorOrange  = color.NRGBA{243, 156, 18, 255}
	ColorPink    = color.NRGBA{255, 151, 197, 255}
	ColorPurple  = color.NRGBA{165, 105, 189, 255}
	ColorSilver  = color.NRGBA{209, 209, 209, 255}
	ColorTeal    = color.NRGBA{64, 199, 178, 255}
	ColorMaroon  = color.NRGBA{199, 54, 103, 255}
	ColorNavy    = color.NRGBA{99, 114, 166, 255}
	ColorOlive   = color.NRGBA{134, 166, 99, 255}
	ColorLime    = color.NRGBA{206, 231, 114, 255}
	ColorFuchsia = color.NRGBA{209, 114, 231, 255}
	ColorAqua    = color.NRGBA{114, 228, 231, 255}
	ColorBrown   = color.NRGBA{176, 116, 78, 255}

	ColorLightRed     = color.NRGBA{255, 83, 83, 255}
	ColorLightGreen   = color.NRGBA{170, 255, 159, 255}
	ColorLightBlue    = color.NRGBA{159, 186, 255, 255}
	ColorLightYellow  = color.NRGBA{255, 251, 159, 255}
	ColorLightGray    = color.NRGBA{236, 236, 236, 255}
	ColorLightOrange  = color.NRGBA{252, 213, 134, 255}
	ColorLightPink    = color.NRGBA{254, 163, 182, 255}
	ColorLightPurple  = color.NRGBA{254, 163, 245, 255}
	ColorLightSilver  = color.NRGBA{228, 228, 228, 255}
	ColorLightTeal    = color.NRGBA{152, 221, 210, 255}
	ColorLightMaroon  = color.NRGBA{215, 124, 143, 255}
	ColorLightNavy    = color.NRGBA{128, 152, 197, 255}
	ColorLightOlive   = color.NRGBA{186, 228, 144, 255}
	ColorLightLime    = color.NRGBA{219, 243, 153, 255}
	ColorLightFuchsia = color.NRGBA{239, 196, 253, 255}
	ColorLightAqua    = color.NRGBA{196, 246, 253, 255}
	ColorLightBrown   = color.NRGBA{211, 140, 94, 255}

	ColorDarkRed     = color.NRGBA{146, 22, 22, 255}
	ColorDarkGreen   = color.NRGBA{22, 146, 24, 255}
	ColorDarkBlue    = color.NRGBA{22, 98, 146, 255}
	ColorDarkYellow  = color.NRGBA{139, 146, 22, 255}
	ColorDarkGray    = color.NRGBA{111, 111, 111, 255}
	ColorCharcoal    = color.NRGBA{16, 16, 16, 255}
	ColorDarkOrange  = color.NRGBA{175, 117, 32, 255}
	ColorDarkPink    = color.NRGBA{128, 64, 64, 255}
	ColorDarkPurple  = color.NRGBA{137, 32, 175, 255}
	ColorDarkSilver  = color.NRGBA{162, 162, 162, 255}
	ColorDarkTeal    = color.NRGBA{27, 110, 86, 255}
	ColorDarkMaroon  = color.NRGBA{110, 27, 55, 255}
	ColorDarkNavy    = color.NRGBA{16, 46, 85, 255}
	ColorDarkOlive   = color.NRGBA{60, 101, 19, 255}
	ColorDarkLime    = color.NRGBA{122, 154, 45, 255}
	ColorDarkFuchsia = color.NRGBA{154, 45, 141, 255}
	ColorDarkAqua    = color.NRGBA{45, 154, 154, 255}
)
