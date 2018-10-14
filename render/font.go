package render

import (
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func loadFont(path string, fontSize float64) font.Face {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	ttf, err := truetype.Parse(c)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: fontSize,
	})
	return face
}
