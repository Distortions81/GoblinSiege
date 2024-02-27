package main

import "image/color"

var (
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

var colorList map[string]color.NRGBA

func init() {
	colorList = make(map[string]color.NRGBA)

	colorList["black"] = ColorBlack
	colorList["red"] = ColorRed
	colorList["green"] = ColorGreen
	colorList["blue"] = ColorBlue
	colorList["yellow"] = ColorYellow
	colorList["gray"] = ColorGray
	colorList["orange"] = ColorOrange
	colorList["pink"] = ColorPink
	colorList["purple"] = ColorPurple
	colorList["silver"] = ColorSilver
	colorList["teal"] = ColorTeal
	colorList["maroon"] = ColorMaroon
	colorList["navy"] = ColorNavy
	colorList["olive"] = ColorOlive
	colorList["Lime"] = ColorLime
	colorList["fuchsia"] = ColorFuchsia
	colorList["aqua"] = ColorAqua
	colorList["brown"] = ColorBrown
}
