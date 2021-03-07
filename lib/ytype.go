package lib

// YType represents the type of the Y coordinate (centre, top or bottom)
type YType int

// YType
const (
	YTop YType = iota
	YCentre
	YBottom
)

// String representation of YType
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
