package colormap

import (
	"testing"
)

func TestNewColorFromString(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{in: "000000", out: "#000000"},
		{in: "#0F12Bc", out: "#0f12bc"},
		{in: "FFffFF", out: "#ffffff"},
	}

	for _, c := range cases {
		col := NewColorFromString(c.in).String()
		if col != c.out {
			t.Errorf("NewColorFromString(%s) returns '%s' instead of '%s'", c.in, col, c.out)
		}
	}
}

func TestColorMap_ColorAt1(t *testing.T) {
	c2 := NewColorFromString("f0ff0f")

	cm := NewColorMap()
	cm.Add(1, c2)

	cases := []struct {
		r   float64
		out string
	}{
		{r: 0, out: "#f0ff0f"},
		{r: 0.5, out: "#f0ff0f"},
		{r: 1, out: "#f0ff0f"},
		{r: 2, out: "#f0ff0f"},
	}
	for _, c := range cases {
		col := cm.ColorAt(c.r).String()
		if col != c.out {
			t.Errorf("NewColorFromString(%0.2f) returns '%s' instead of '%s'", c.r, col, c.out)
		}
	}
}

func TestColorMap_ColorAt2(t *testing.T) {
	c1 := NewColorFromString("000000")
	c2 := NewColorFromString("f0ff0f")

	cm := NewColorMap()
	cm.Add(0, c1)
	cm.Add(1, c2)

	cases := []struct {
		r   float64
		out string
	}{
		{r: -1, out: "#000000"},
		{r: 0, out: "#000000"},
		{r: 0.5, out: "#787f07"},
		{r: 1, out: "#f0ff0f"},
		{r: 2, out: "#f0ff0f"},
	}
	for _, c := range cases {
		col := cm.ColorAt(c.r).String()
		if col != c.out {
			t.Errorf("NewColorFromString(%0.2f) returns '%s' instead of '%s'", c.r, col, c.out)
		}
	}
}

func TestColorMap_ColorAt3(t *testing.T) {
	c1 := NewColorFromString("000000")
	c2 := NewColorFromString("00ff00")
	c3 := NewColorFromString("ff00ff")

	cm := NewColorMap()
	cm.Add(2, c3)
	cm.Add(0, c1)
	cm.Add(1, c2)

	cases := []struct {
		r   float64
		out string
	}{
		{r: -1, out: "#000000"},
		{r: 0, out: "#000000"},
		{r: 0.5, out: "#007f00"},
		{r: 1, out: "#00ff00"},
		{r: 1.5, out: "#7f807f"},
		{r: 2, out: "#ff00ff"},
		{r: 3, out: "#ff00ff"},
	}
	for _, c := range cases {
		col := cm.ColorAt(c.r).String()
		if col != c.out {
			t.Errorf("NewColorFromString(%0.2f) returns '%s' instead of '%s'", c.r, col, c.out)
		}
	}
}
