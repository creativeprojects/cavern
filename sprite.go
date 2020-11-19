package main

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// XType represents the type of the X coordinate (centre, left or right)
type XType int

// XType
const (
	XLeft XType = iota
	XCentre
	XRight
)

func (x XType) String() string {
	switch x {
	case XCentre:
		return "X center"
	case XRight:
		return "X right"
	default:
		return "X left"
	}
}

// YType represents the type of the Y coordinate (centre, top or bottom)
type YType int

// YType
const (
	YTop YType = iota
	YCentre
	YBottom
)

func (y YType) String() string {
	switch y {
	case YCentre:
		return "Y center"
	case YBottom:
		return "Y bottom"
	default:
		return "Y top"
	}
}

// SequenceFunc is used as a callback to decide which image to draw
type SequenceFunc func(int) int

// Sprite manages sprite movement and animation
type Sprite struct {
	xType        XType
	yType        YType
	x            float64
	y            float64
	frame        int             // current frame counter (animation mode)
	width        int             // fixed size only
	height       int             // fixed size only
	image        *ebiten.Image   // single image (no animation)
	animation    []*ebiten.Image // animation images
	sequence     []int           // animation sequence (in case same images need reused)
	sequenceFunc SequenceFunc    // use a callback to decide on animation sequence instead
	rate         int             // change image every n frame per second
	loop         bool            // animation loop
	started      bool            // is animation running?
	op           *ebiten.DrawImageOptions
}

// NewSprite creates a new Sprite with default coordinate type
func NewSprite(xType XType, yType YType) *Sprite {
	return &Sprite{
		xType:   xType,
		yType:   yType,
		op:      &ebiten.DrawImageOptions{},
		started: false,
	}
}

// String returns a debug string representation of the sprite current state
func (s *Sprite) String() string {
	return fmt.Sprintf("x: %.0f, y: %.0f, frame: %d",
		s.x,
		s.y,
		s.frame,
	)
}

// SetImage sets sprite image
func (s *Sprite) SetImage(image *ebiten.Image) *Sprite {
	s.image = image
	return s
}

// SetSize forces width and height of the image. This is used for calculation only (does not resize the images)
func (s *Sprite) SetSize(width, height int) *Sprite {
	s.width = width
	s.height = height
	return s
}

// Update animation (if needed)
func (s *Sprite) Update() {
	if !s.started {
		return
	}
	s.frame++
	if !s.loop && s.frame >= s.rate*len(s.animation) {
		// animation is finished
		s.started = false
		return
	}
	// update current frame
	frameID := s.getFrameID()
	s.image = s.animation[frameID]
}

// Draw the current image, or the animation to the screen. If no image or animation has been set, it does nothing
func (s *Sprite) Draw(screen *ebiten.Image) {
	if s.image == nil {
		log.Println("Sprite.Draw: no image to draw")
		return
	}
	width, height := s.image.Size()
	s.op.GeoM.Reset()
	s.op.GeoM.Translate(s.xleft(float64(width)), s.ytop(float64(height)))
	screen.DrawImage(s.image, s.op)
}

func (s *Sprite) getFrameID() int {
	if s.sequenceFunc != nil {
		return s.sequenceFunc(s.frame)
	}
	current := float64(s.frame / s.rate)
	if s.sequence == nil || len(s.sequence) == 0 {
		// no sequence, we just go through all images one by one
		frameID := int(math.Mod(current, float64(len(s.animation))))
		return frameID
	}
	frameID := int(math.Mod(current, float64(len(s.sequence))))
	// if frameID is over the size of animation, we pick the last one
	if frameID >= len(s.animation) {
		frameID = len(s.animation) - 1
	}
	return frameID
}

// Start (or restart) an animation
func (s *Sprite) Start() *Sprite {
	// only start if the animation is well defined
	if s.animation != nil && len(s.animation) > 0 && s.rate > 0 {
		s.frame = 0
		s.started = true
	}
	return s
}

// Stop animation
func (s *Sprite) Stop() *Sprite {
	s.started = false
	return s
}

// Animation defines a new animation (but does not start it yet)
func (s *Sprite) Animation(animation []*ebiten.Image, sequence []int, rate int, loop bool) *Sprite {
	s.animation = animation
	s.sequence = sequence
	s.rate = rate
	s.loop = loop
	if s.animation != nil && len(s.animation) > 0 {
		s.image = s.animation[0]
	}
	return s
}

// Animate defines a new animation and starts it
func (s *Sprite) Animate(animation []*ebiten.Image, sequence []int, rate int, loop bool) *Sprite {
	s.Animation(animation, sequence, rate, loop)
	return s.Start()
}

// SetSequenceFunc registers a callback to calculate the next image.
// the value returned should be a valid index in the animation slice passed to the Animation method
func (s *Sprite) SetSequenceFunc(sequenceFunc SequenceFunc) *Sprite {
	s.sequenceFunc = sequenceFunc
	return s
}

// IsFinished returns true when the *Sprite animation has finished.
// An animation with loop = true will never finish
func (s *Sprite) IsFinished() bool {
	return !s.started
}

// Move to relative coordinates (adds coordinates to the current position)
func (s *Sprite) Move(x, y float64) *Sprite {
	s.x += x
	s.y += y
	return s
}

// MoveTo the new coordinates using the default coordinates type defined at instantiation
func (s *Sprite) MoveTo(x, y float64) *Sprite {
	s.x = x
	s.y = y
	return s
}

// MoveToType moves to the new coordinates using the specified coordinate types
func (s *Sprite) MoveToType(x, y float64, xType XType, yType YType) *Sprite {
	if s.xType == xType {
		s.x = x
	} else {
		panic(fmt.Sprintf("mixing different types of coordinates is not yet supported: want to set %s but is %s", xType.String(), s.xType.String()))
	}
	if s.yType == yType {
		s.y = y
	} else if s.yType == YCentre && yType == YBottom {
		_, height := s.image.Size()
		s.y = y - (float64(height) / 2)
	} else {
		panic(fmt.Sprintf("mixing different types of coordinates is not yet supported: want to set %s but is %s", yType.String(), s.yType.String()))
	}
	return s
}

// X returns x position. If not image is available to calculate width, it returns -1
func (s *Sprite) X(xType XType) float64 {
	if s.image == nil {
		return -1
	}
	width, _ := s.image.Size()
	switch xType {
	case XCentre:
		return s.xcentre(float64(width))
	case XRight:
		return s.xright(float64(width))
	default:
		return s.xleft(float64(width))
	}
}

// Y returns y position. If not image is available to calculate height, it returns -1
func (s *Sprite) Y(yType YType) float64 {
	if s.image == nil {
		return -1
	}
	_, height := s.image.Size()
	switch yType {
	case YCentre:
		return s.ycentre(float64(height))
	case YBottom:
		return s.ybottom(float64(height))
	default:
		return s.ytop(float64(height))
	}
}

// CollidePoint returns true when the coordinates are "touching" the sprite
func (s *Sprite) CollidePoint(x, y float64) bool {
	return s.X(XLeft) <= x && x <= s.X(XRight) &&
		s.Y(YTop) <= y && y <= s.Y(YBottom)
}

func (s *Sprite) xleft(width float64) float64 {
	switch s.xType {
	case XCentre:
		return s.x - (width / 2)
	case XRight:
		return s.x - width
	default:
		return s.x
	}
}

func (s *Sprite) xcentre(width float64) float64 {
	switch s.xType {
	case XCentre:
		return s.x
	case XRight:
		return s.x - (width / 2)
	default:
		return s.x + (width / 2)
	}
}

func (s *Sprite) xright(width float64) float64 {
	switch s.xType {
	case XCentre:
		return s.x + (width / 2)
	case XRight:
		return s.x
	default:
		return s.x + width
	}
}

func (s *Sprite) ytop(height float64) float64 {
	switch s.yType {
	case YCentre:
		return s.y - (height / 2)
	case YBottom:
		return s.y - height
	default:
		return s.y
	}
}

func (s *Sprite) ycentre(height float64) float64 {
	switch s.yType {
	case YCentre:
		return s.y
	case YBottom:
		return s.y - (height / 2)
	default:
		return s.y + (height / 2)
	}
}

func (s *Sprite) ybottom(height float64) float64 {
	switch s.yType {
	case YCentre:
		return s.y + (height / 2)
	case YBottom:
		return s.y
	default:
		return s.y + height
	}
}
