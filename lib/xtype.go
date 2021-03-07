package lib

// XType represents the type of the X coordinate (centre, left or right)
type XType int

// XType
const (
	XLeft XType = iota
	XCentre
	XRight
)

// String representation of XType
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
