package colormap

import (
	"strconv"
	"strings"
)

type Color struct {
	R int
	G int
	B int
	A int
}

// NewColorFromString returns a Color from string web color description #rrggbb
func NewColorFromString(c string) Color {
	c = strings.TrimLeft(c, "#")
	ic := str2int(c)
	return Color{
		R: int((ic & 0xff0000) >> 16),
		G: int((ic & 0x00ff00) >> 8),
		B: int(ic & 0x0000ff),
		A: 0,
	}
}

func NewColorRGB(r, g, b int) Color {
	return Color{r, g, b, 0}
}

func (c Color) String() string {
	return "#" + rgb2str(c.R, c.G, c.B)
}

func str2int(s string) int {
	i, err := strconv.ParseInt(s, 16, 32)
	if err != nil {
		i = 0
	}
	return int(i)
}

func rgb2str(r, g, b int) string {
	return strconv.FormatInt(int64(1<<24+r<<16+g<<8+b), 16)[1:]
}

type ColorMap struct {
	Range  []float64
	Colors []Color
}

func NewColorMap() ColorMap {
	return ColorMap{Range: []float64{}, Colors: []Color{}}
}

func (cm *ColorMap) Add(r float64, c Color) {
	var i int
	for i = 0; i < len(cm.Range); i++ {
		if r < cm.Range[i] {
			break
		}
	}

	cm.Range = append(cm.Range, 0)
	copy(cm.Range[i+1:], cm.Range[i:])
	cm.Range[i] = r

	cm.Colors = append(cm.Colors, Color{})
	copy(cm.Colors[i+1:], cm.Colors[i:])
	cm.Colors[i] = c
}

func (cm *ColorMap) ColorAt(r float64) Color {
	if len(cm.Range) == 1 {
		return cm.Colors[0]
	}
	if r <= cm.Range[0] {
		return cm.Colors[0]
	}
	if r >= cm.Range[len(cm.Range)-1] {
		return cm.Colors[len(cm.Colors)-1]
	}
	var i int
	for i = 1; i < len(cm.Range); i++ {
		if r < cm.Range[i] {
			break
		}
	}

	dr := (r - cm.Range[i-1])/(cm.Range[i] - cm.Range[i-1])

	return NewColorRGB(
		cm.Colors[i-1].R+int(dr*float64(cm.Colors[i].R-cm.Colors[i-1].R)),
		cm.Colors[i-1].G+int(dr*float64(cm.Colors[i].G-cm.Colors[i-1].G)),
		cm.Colors[i-1].B+int(dr*float64(cm.Colors[i].B-cm.Colors[i-1].B)),
	)
}
