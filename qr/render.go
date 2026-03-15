package qr

import (
	"image"
	"image/color"
	"math"

	qrcode "github.com/skip2/go-qrcode"
)

type DotStyle int

const (
	DotSquare DotStyle = iota
	DotCircle
	DotRounded
)

type EyeStyle int

const (
	EyeSquare EyeStyle = iota
	EyeRounded
	EyeCircle
)

type RenderOptions struct {
	PixelSize     int
	DotStyle      DotStyle
	EyeFrameStyle EyeStyle
	EyePupilStyle EyeStyle
	Foreground    color.Color
	Background    color.Color
}

func GenerateMatrix(text string, level qrcode.RecoveryLevel) ([][]bool, error) {
	qr, err := qrcode.New(text, level)
	if err != nil {
		return nil, err
	}
	qr.DisableBorder = true
	return qr.Bitmap(), nil
}

// finderRegions returns the three 7x7 finder pattern bounding boxes as (row, col) origins.
func finderRegions(matrixSize int) [3][2]int {
	return [3][2]int{
		{0, 0},              // top-left
		{0, matrixSize - 7}, // top-right
		{matrixSize - 7, 0}, // bottom-left
	}
}

func isFinderModule(row, col, matrixSize int) bool {
	for _, origin := range finderRegions(matrixSize) {
		r0, c0 := origin[0], origin[1]
		if row >= r0 && row < r0+7 && col >= c0 && col < c0+7 {
			return true
		}
	}
	return false
}

func RenderImage(matrix [][]bool, opts RenderOptions) *image.RGBA {
	size := len(matrix)
	imgSize := size * opts.PixelSize
	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))

	bgCol := toRGBA(opts.Background)
	fgCol := toRGBA(opts.Foreground)

	// Fill background
	for y := 0; y < imgSize; y++ {
		for x := 0; x < imgSize; x++ {
			img.SetRGBA(x, y, bgCol)
		}
	}

	// Render data modules (skip finder pattern areas)
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			if isFinderModule(row, col, size) {
				continue
			}
			if !matrix[row][col] {
				continue
			}
			px := col * opts.PixelSize
			py := row * opts.PixelSize
			switch opts.DotStyle {
			case DotSquare:
				drawRect(img, px, py, opts.PixelSize, opts.PixelSize, fgCol)
			case DotCircle:
				drawFilledCircle(img, px, py, opts.PixelSize, opts.PixelSize, fgCol)
			case DotRounded:
				drawRoundedRect(img, px, py, opts.PixelSize, opts.PixelSize, float64(opts.PixelSize)*0.3, fgCol)
			}
		}
	}

	// Render finder patterns (eyes)
	for _, origin := range finderRegions(size) {
		r0, c0 := origin[0], origin[1]
		drawEye(img, c0*opts.PixelSize, r0*opts.PixelSize, 7*opts.PixelSize, opts, fgCol, bgCol)
	}

	return img
}

func drawEye(img *image.RGBA, x, y, totalSize int, opts RenderOptions, fg, bg color.RGBA) {
	// Eye frame (outer shape)
	drawEyeShape(img, x, y, totalSize, totalSize, opts.EyeFrameStyle, fg)

	// Gap (background ring) — inset by 1 module
	inset1 := opts.PixelSize
	gapSize := totalSize - 2*inset1
	drawEyeShape(img, x+inset1, y+inset1, gapSize, gapSize, opts.EyeFrameStyle, bg)

	// Pupil (inner solid) — inset by 2 modules from outer
	inset2 := 2 * opts.PixelSize
	pupilSize := totalSize - 2*inset2
	drawEyeShape(img, x+inset2, y+inset2, pupilSize, pupilSize, opts.EyePupilStyle, fg)
}

func drawEyeShape(img *image.RGBA, x, y, w, h int, style EyeStyle, col color.RGBA) {
	switch style {
	case EyeSquare:
		drawRect(img, x, y, w, h, col)
	case EyeRounded:
		drawRoundedRect(img, x, y, w, h, float64(w)*0.25, col)
	case EyeCircle:
		drawFilledCircle(img, x, y, w, h, col)
	}
}

// --- drawing primitives ---

func drawRect(img *image.RGBA, x, y, w, h int, col color.RGBA) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.SetRGBA(x+dx, y+dy, col)
		}
	}
}

func drawFilledCircle(img *image.RGBA, x, y, w, h int, col color.RGBA) {
	cx := float64(x) + float64(w)/2.0
	cy := float64(y) + float64(h)/2.0
	rx := float64(w) / 2.0
	ry := float64(h) / 2.0

	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			px := float64(x+dx) + 0.5
			py := float64(y+dy) + 0.5
			// Ellipse equation: ((px-cx)/rx)^2 + ((py-cy)/ry)^2 <= 1
			nx := (px - cx) / rx
			ny := (py - cy) / ry
			if nx*nx+ny*ny <= 1.0 {
				img.SetRGBA(x+dx, y+dy, col)
			}
		}
	}
}

func drawRoundedRect(img *image.RGBA, x, y, w, h int, radius float64, col color.RGBA) {
	maxR := math.Min(float64(w), float64(h)) / 2.0
	if radius > maxR {
		radius = maxR
	}

	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			px := float64(dx) + 0.5
			py := float64(dy) + 0.5

			inside := true
			// Check corners
			if px < radius && py < radius {
				// top-left corner
				inside = cornerDist(px, py, radius, radius, radius)
			} else if px > float64(w)-radius && py < radius {
				// top-right corner
				inside = cornerDist(px, py, float64(w)-radius, radius, radius)
			} else if px < radius && py > float64(h)-radius {
				// bottom-left corner
				inside = cornerDist(px, py, radius, float64(h)-radius, radius)
			} else if px > float64(w)-radius && py > float64(h)-radius {
				// bottom-right corner
				inside = cornerDist(px, py, float64(w)-radius, float64(h)-radius, radius)
			}

			if inside {
				img.SetRGBA(x+dx, y+dy, col)
			}
		}
	}
}

func cornerDist(px, py, cx, cy, r float64) bool {
	dx := px - cx
	dy := py - cy
	return dx*dx+dy*dy <= r*r
}

func toRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}
