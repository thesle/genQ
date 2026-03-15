package ui

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type Preview struct {
	widget.BaseWidget
	image    *canvas.Image
	minSize  fyne.Size
}

func NewPreview(size float32) *Preview {
	p := &Preview{
		minSize: fyne.NewSize(size, size),
	}
	p.image = canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 1, 1)))
	p.image.FillMode = canvas.ImageFillContain
	p.image.ScaleMode = canvas.ImageScalePixels
	p.ExtendBaseWidget(p)
	return p
}

func (p *Preview) SetImage(img image.Image) {
	p.image.Image = img
	p.image.Refresh()
}

func (p *Preview) CreateRenderer() fyne.WidgetRenderer {
	return &previewRenderer{preview: p}
}

type previewRenderer struct {
	preview *Preview
}

func (r *previewRenderer) Layout(size fyne.Size) {
	r.preview.image.Resize(size)
	r.preview.image.Move(fyne.NewPos(0, 0))
}

func (r *previewRenderer) MinSize() fyne.Size {
	return r.preview.minSize
}

func (r *previewRenderer) Refresh() {
	r.preview.image.Refresh()
}

func (r *previewRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.preview.image}
}

func (r *previewRenderer) Destroy() {}
