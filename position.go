package main

// Position is immutable. Changing coordinates returns a new Position.
type Position struct {
	absoluteX  float64
	absoluteY  float64
	halfWidth  float64
	halfHeight float64
}

// NewPositionAbsolute returns a new immutable Position instance using absolute coordinates
func NewPositionAbsolute(width, height, absoluteX, absoluteY float64) Position {
	return Position{
		halfWidth:  width / 2,
		halfHeight: height / 2,
		absoluteX:  absoluteX,
		absoluteY:  absoluteY,
	}
}

// NewPositionCentre returns a new immutable Position instance using centered coordinates
func NewPositionCentre(width, height, centreX, centreY float64) Position {
	return Position{
		halfWidth:  width / 2,
		halfHeight: height / 2,
		absoluteX:  centreX - (width / 2),
		absoluteY:  centreY - (height / 2),
	}
}

// Absolute returns the absolute position
func (p Position) Absolute() (x float64, y float64) {
	x = p.absoluteX
	y = p.absoluteY
	return
}

// AbsoluteX returns the absolute position on the X axis
func (p Position) AbsoluteX() float64 {
	return p.absoluteX
}

// AbsoluteY returns the absolute position on the Y axis
func (p Position) AbsoluteY() float64 {
	return p.absoluteY
}

// Centre returns the centered position
func (p Position) Centre() (x float64, y float64) {
	x = p.absoluteX + p.halfWidth
	y = p.absoluteY + p.halfHeight
	return
}

// CentreX returns the centered position
func (p Position) CentreX() float64 {
	return p.absoluteX + p.halfWidth
}

// CentreY returns the centered position
func (p Position) CentreY() float64 {
	return p.absoluteY + p.halfHeight
}

// MoveAbsolute moves to an absolute position and returns the new position
func (p Position) MoveAbsolute(x, y float64) Position {
	return Position{
		absoluteX:  x,
		absoluteY:  y,
		halfWidth:  p.halfWidth,
		halfHeight: p.halfHeight,
	}
}

// MoveCentre moves to a centered position and returns the new position
func (p Position) MoveCentre(x, y float64) Position {
	return Position{
		absoluteX:  x - p.halfWidth,
		absoluteY:  y - p.halfHeight,
		halfWidth:  p.halfWidth,
		halfHeight: p.halfHeight,
	}
}

// MoveRelative increments the current position values and returns the new position
func (p Position) MoveRelative(x, y float64) Position {
	return Position{
		absoluteX:  p.absoluteX + x,
		absoluteY:  p.absoluteY + y,
		halfWidth:  p.halfWidth,
		halfHeight: p.halfHeight,
	}
}
