package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type UiVerticalPosition string
type UiHorizontalPosition string
type EventListener func()

const (
	HCentered UiHorizontalPosition = "hcenter"
	Top       UiHorizontalPosition = "top"
	Right     UiVerticalPosition   = "right"

	VCentered UiVerticalPosition   = "vcenter"
	Left      UiVerticalPosition   = "left"
	Bottom    UiHorizontalPosition = "bottom"
)

type UiElement struct {
	Rectangle       rl.Rectangle
	BackgroundColor rl.Color
	BorderColor     rl.Color
	BorderWidth     int
	Margin          UiMargin
	VPosition       UiVerticalPosition
	HPosition       UiHorizontalPosition
	Container 		rl.Rectangle
	Childs          []*UiElement
	Listeners       map[string][]EventListener
	Text            string
}

type UiMargin struct {
	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

type NewUiElementInput struct {
	Width           float32
	Height          float32
	BackgroundColor rl.Color
	BorderColor     rl.Color
	BorderWidth     int
	HPosition       UiHorizontalPosition
	VPosition       UiVerticalPosition
	Margin          UiMargin
	Text            string
}

func NewUiElement(properties NewUiElementInput) UiElement {
	container := rl.NewRectangle(0, 0, float32(VIRTUAL_WINDOW_WIDTH), float32(VIRTUAL_WINDOW_HEIGHT))
	xPosition, yPosition := getRectanglePosition(container, properties.Width, properties.Height, properties.HPosition, properties.VPosition, properties.Margin)
	rectangle := rl.NewRectangle(xPosition, yPosition, properties.Width, properties.Height)

	listeners := make(map[string][]EventListener)

	listeners["click"] = []EventListener{}
	listeners["hover"] = []EventListener{}

	return UiElement{
		Rectangle: rectangle,
		BackgroundColor: properties.BackgroundColor,
		BorderColor: properties.BorderColor,
		BorderWidth: properties.BorderWidth,
		Margin: properties.Margin,
		HPosition: properties.HPosition,
		VPosition: properties.VPosition,
		Container: container,
		Listeners: listeners,
		Text: properties.Text,
	}
}

func (e *UiElement) Tick() {
	e.ComputePosition()
	e.RunPreDrawEvents()
	e.Draw()
	e.DrawText()
	e.RunPostDrawEvents()
}

func (e *UiElement) ComputePosition() {
	xPosition, yPosition := getRectanglePosition(e.Container, e.Rectangle.Width, e.Rectangle.Height, e.HPosition, e.VPosition, e.Margin)

	e.Rectangle.X = xPosition
	e.Rectangle.Y = yPosition
}

func (e *UiElement) Draw() {
	rl.DrawRectangleRec(e.Rectangle, e.BackgroundColor)
	rl.DrawRectangleLinesEx(e.Rectangle, float32(e.BorderWidth), e.BorderColor)

	for _, child := range e.Childs {
		child.Tick()
	}
}

func (e *UiElement) DrawText() {
	if e.Text == "" {
		return
	}

	xPosition, yPosition := getRectanglePosition(e.Rectangle, e.Rectangle.Width / 2, e.Rectangle.Height / 2, HCentered, VCentered, UiMargin{})
	rl.DrawText(e.Text, int32(xPosition), int32(yPosition), 7, rl.White)
}

func (e *UiElement) RunPreDrawEvents() {
	if e.IsMouseHovering() {
		for _, listener := range e.Listeners["hover"] {
			listener()
		}
	}
}

func (e *UiElement) RunPostDrawEvents() {
	if e.IsMouseHovering() && rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		for _, listener := range e.Listeners["click"] {
			listener()
		}
	}
}

func (e *UiElement) AddChild(child *UiElement) {
	child.Container = e.Rectangle
	e.Childs = append(e.Childs, child)
}

func (e *UiElement) AddEventListener(event string, listener EventListener) {
	e.Listeners[event] = append(e.Listeners[event], listener)
}

func (e *UiElement) SetBackgroundColor(color rl.Color) {
	e.BackgroundColor = color
}

func getRectanglePosition(container rl.Rectangle, width, height float32, hposition UiHorizontalPosition, vposition UiVerticalPosition, margin UiMargin) (float32, float32) {
	xPosition := container.X // Left by default
	yPosition := container.Y // Top by default

	if hposition == Bottom {
		yPosition = container.Y + container.Height - height
	}

	if hposition == HCentered {
		yPosition = container.Y + (container.Height / 2 - height / 2)
	}

	if vposition == Right {
		xPosition = container.X + container.Width - width
	}

	if vposition == VCentered {
		xPosition = container.X + (container.Width / 2 - width / 2)
	}

	xPosition += margin.Left - margin.Right
	yPosition += margin.Top - margin.Bottom

	return xPosition, yPosition
}

func (e *UiElement) IsMouseHovering() bool {
	mousePosition := rl.GetMousePosition()
	mouseX := mousePosition.X
	mouseY := mousePosition.Y
	mousePoint := rl.Vector2{X: mouseX, Y: mouseY}

	return rl.CheckCollisionPointRec(mousePoint, e.Rectangle)
}
