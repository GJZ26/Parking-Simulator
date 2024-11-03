package geo

type PointData struct {
	Id     uint32
	X      float64
	Y      float64
	Width  float64
	Height float64
	Type   string
}

type PointMap struct {
	Road             []PointData
	LeftParkingSlot  []PointData
	RightParkingSlot []PointData
	ParkingRoad      []PointData
}

type Sides byte

const (
	LeftSide  Sides = 0
	RightSide Sides = 1
)

type SlotInfo struct {
	Id       uint32
	X        float64
	Y        float64
	Width    float64
	Height   float64
	Type     string
	Occupied bool
	Side     Sides
}
