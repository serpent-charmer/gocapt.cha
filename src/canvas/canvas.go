package canvas 

import (
	"log"
	"bytes"
	"image"
	"image/color"
	"image/png"
	"image/draw" 
	"math/rand"
	"gocapt.cha/mask"
	"encoding/base64"
)

type BrushFunc func(color.RGBA) color.RGBA

type CaptchaOut struct {
	Key string `json:"key"`
	Canvas string `json:"canvas"`
	Elements []string `json:"elements"`
	Solution CaptchaSolution `json:"-"`
}

type CaptchaSolution struct {
	Index int
	Solution image.Rectangle
}

type CaptchaRequest struct {
	Key string `json:"key"`
	Index int `json:"index"`
	Position image.Point `json:"pos"`
}

func RandImg(w int, h int, bf BrushFunc) *image.RGBA {
    img := image.NewRGBA(image.Rect(0, 0, w, h))
	
	brush := color.RGBA{R: 255, A: 255}
	
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, bf(brush))
		}
	}
	return img
}

func ImgToBytes(img image.Image) bytes.Buffer {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		log.Println(err)
	}
	return buf
}

func brush1(brush color.RGBA) color.RGBA {
	brush.R = uint8(rand.Intn(255))
	brush.B = uint8(rand.Intn(255))
	return brush
}

func brush2(brush color.RGBA) color.RGBA {
	brush.R = uint8(rand.Intn(255)+120)
	brush.G = uint8(rand.Intn(255)+120)
	return brush
}

func brush3(brush color.RGBA) color.RGBA {
	brush.R = uint8(rand.Intn(100)+200)
	brush.G = uint8(rand.Intn(255)+200)
	brush.B = uint8(rand.Intn(100)+200)
	return brush
}

func MakeCaptcha() *CaptchaOut {
	const imgWidth = 64
	const imgHeight = 64
	var dto CaptchaOut
	background := RandImg(512, 512, brush3)
	mask := mask.MakeCircle(image.Pt(imgWidth/2, imgHeight/2), imgWidth/4)
	canvas := image.NewRGBA(image.Rect(0, 0, 512, 512))
	draw.Draw(canvas, canvas.Bounds(), background, image.ZP, draw.Over)
	arrLen := 6
	elements := make([]string, 0)
	solution := rand.Intn(arrLen)
	for i := 0; i < arrLen; i++ {
		element := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
		xx := i * (512/arrLen)
		yy := 280
		yy += rand.Intn(64)-128
		elementRect := image.Rect(xx, yy, xx+imgWidth, yy+imgHeight)
		
		var img image.Image
		if i == solution {
			img = RandImg(imgWidth, imgHeight, brush1)
			dto.Solution = CaptchaSolution{i, elementRect}
		} else {
			img = RandImg(imgWidth, imgHeight, brush2)
		}
		draw.DrawMask(canvas, elementRect, img, image.ZP, mask, image.ZP, draw.Over)
		draw.DrawMask(element, element.Bounds(), img, image.ZP, mask, image.ZP, draw.Over)
		imgBytes := ImgToBytes(element)
		encodedImg := base64.StdEncoding.EncodeToString(imgBytes.Bytes())
		elements = append(elements, encodedImg)
	}
	buf := ImgToBytes(canvas)
	dto.Canvas = base64.StdEncoding.EncodeToString(buf.Bytes())
	dto.Elements = elements
	return &dto
}
