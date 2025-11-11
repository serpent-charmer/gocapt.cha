package canvas

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math/rand"

	"github.com/google/uuid"
	"gocapt.cha/mask"
)

var captchaCache = make(map[string]*DummyCaptchaSolution)

type BrushFunc func(color.RGBA) color.RGBA

type ClientCaptchaSolution struct {
	Key      string      `json:"key"`
	Index    int         `json:"index"`
	Position image.Point `json:"pos"`
}

func (c ClientCaptchaSolution) Validate() error {
	defer delete(captchaCache, c.Key)
	captcha := captchaCache[c.Key]
	if captcha == nil {
		return errors.New("captcha not found")
	}
	if c.Index != captcha.index {
		return errors.New("wrong element")
	}
	if !c.Position.In(captcha.position) {
		return errors.New("not in bounds")
	}
	return nil
}

func SolutionFromJson(stream io.ReadCloser) (*ClientCaptchaSolution, error) {
	var captchaRequest ClientCaptchaSolution
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&captchaRequest)
	if err != nil {
		return nil, err
	}
	return &captchaRequest, nil
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

func EncodeImg(img image.Image) string {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func brush1(brush color.RGBA) color.RGBA {
	brush.R = uint8(rand.Intn(255))
	brush.B = uint8(rand.Intn(255))
	return brush
}

func brush2(brush color.RGBA) color.RGBA {
	brush.R = uint8(rand.Intn(255) + 120)
	brush.G = uint8(rand.Intn(255) + 120)
	return brush
}

func brush3(brush color.RGBA) color.RGBA {
	brush.R = uint8(rand.Intn(100) + 200)
	brush.G = uint8(rand.Intn(255) + 200)
	brush.B = uint8(rand.Intn(100) + 200)
	return brush
}

type DummyCaptchaElement struct {
	image    draw.Image
	position image.Rectangle
}

type DummyCaptchaSolution struct {
	index    int
	position image.Rectangle
}

type DummyCaptcha struct {
	canvas   draw.Image
	elements []DummyCaptchaElement
	solution string
}

type DummyCaptchaJson struct {
	Key      string   `json:"key"`
	Elements []string `json:"elements"`
	Canvas   string   `json:"canvas"`
}

func Make() DummyCaptcha {
	const imgWidth = 64
	const imgHeight = 64

	length := 10
	c := DummyCaptcha{}
	c.elements = make([]DummyCaptchaElement, length)
	brushes := make([]BrushFunc, length)

	for i, _ := range c.elements {
		v := &c.elements[i]
		xx := i * (512 / length)
		yy := 280
		yy += rand.Intn(64) - 128
		elementRect := image.Rect(xx, yy, xx+imgWidth, yy+imgHeight)
		v.position = elementRect
		brushes[i] = brush2
	}

	solution := rand.Intn(length)
	brushes[solution] = brush1
	el := c.elements[solution]
	captchaSolution := &DummyCaptchaSolution{index: solution, position: el.position}

	background := RandImg(512, 512, brush3)
	c.canvas = image.NewRGBA(image.Rect(0, 0, 512, 512))
	mask := mask.MakeCircle(image.Pt(imgWidth/2, imgWidth/2), imgWidth/4)
	draw.Draw(c.canvas, c.canvas.Bounds(), background, image.Point{}, draw.Over)

	for i, _ := range c.elements {
		v := &c.elements[i]
		surface := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
		img := RandImg(imgWidth, imgHeight, brushes[i])
		draw.DrawMask(c.canvas, v.position, img, image.Point{}, mask, image.Point{}, draw.Over)
		draw.DrawMask(surface, surface.Bounds(), img, image.Point{}, mask, image.Point{}, draw.Over)
		v.image = surface
	}

	captchaKey := uuid.NewString()
	captchaCache[captchaKey] = captchaSolution

	c.solution = captchaKey

	return c
}

func (c DummyCaptcha) ToJson() ([]byte, error) {
	if c.elements == nil {
		return nil, errors.New("elements are null")
	}
	if c.canvas == nil {
		return nil, errors.New("canvas is null")
	}
	l := len(c.elements)
	if l < 1 {
		return nil, errors.New("no elements")
	}
	encoded := make([]string, l)
	for i, v := range c.elements {
		encoded[i] = EncodeImg(v.image)
	}
	encodedCanvas := EncodeImg(c.canvas)
	captchaJson := DummyCaptchaJson{Key: c.solution, Elements: encoded, Canvas: encodedCanvas}
	bytes, err := json.Marshal(captchaJson)
	if err != nil {
		panic(err)
	}
	return bytes, nil
}
