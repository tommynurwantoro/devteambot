package constant

type Color struct {
	Blue   int
	Yellow int
	Green  int
}

func NewColor() Color {
	return Color{
		Blue:   39423,
		Yellow: 16772608,
		Green:  52373,
	}
}
