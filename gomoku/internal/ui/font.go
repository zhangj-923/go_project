package ui

import (
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed cjk_font.ttf
var cjkFontData []byte

var (
	smallFont  font.Face
	normalFont font.Face
	titleFont  font.Face
)

func InitFonts() {
	fontData := cjkFontData
	if len(fontData) == 0 {
		fontData = fonts.MPlus1pRegular_ttf
	}

	tt, err := opentype.Parse(fontData)
	if err != nil {
		log.Printf("failed to parse custom font, falling back to mplus: %v", err)
		tt, err = opentype.Parse(fonts.MPlus1pRegular_ttf)
		if err != nil {
			log.Fatalf("failed to parse fallback font: %v", err)
		}
	}

	const dpi = 72
	smallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    13,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("failed to create small font: %v", err)
	}

	normalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("failed to create normal font: %v", err)
	}

	titleFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    30,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("failed to create title font: %v", err)
	}
}
