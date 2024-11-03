package routines

import (
	"Parking_Simulator/src/core/manager/types/geo"
)

type SlotManager struct {
	slots map[uint32]geo.SlotInfo
}

func NewSlotManager(pointMap geo.PointMap) *SlotManager {
	slots := map[uint32]geo.SlotInfo{}
	for _, slot := range pointMap.LeftParkingSlot {
		slots[slot.Id] = geo.SlotInfo{
			Id:       slot.Id,
			X:        slot.X,
			Y:        slot.Y,
			Width:    slot.Width,
			Height:   slot.Height,
			Type:     slot.Type,
			Occupied: false,
			Side:     geo.LeftSide,
		}
	}

	for _, slot := range pointMap.RightParkingSlot {
		slots[slot.Id] = geo.SlotInfo{
			Id:       slot.Id,
			X:        slot.X,
			Y:        slot.Y,
			Width:    slot.Width,
			Height:   slot.Height,
			Type:     slot.Type,
			Occupied: false,
			Side:     geo.RightSide,
		}
	}
	return &SlotManager{slots}
}

func (s *SlotManager) Run(slotInfo chan geo.SlotInfo, freeSlotChannel chan []uint32) {
	for {
		select {
		case freeSlot := <-freeSlotChannel:
			s.freeSlots(freeSlot)
		default:
		}
		data := s.sendCurrentCarData()

		if data.X != -1 {
			slotInfo <- data
		}

	}
}

func (s *SlotManager) freeSlots(slots []uint32) {
	//println("[Slot Manager]: Free Slots")
	for _, id := range slots {
		if slot, exists := s.slots[id]; exists {
			slot.Occupied = false
			s.slots[id] = slot
		}
	}
}

func (s *SlotManager) sendCurrentCarData() geo.SlotInfo {
	for id, slot := range s.slots {
		if !slot.Occupied {
			slot.Occupied = true
			s.slots[id] = slot
			println("[Slot Manager]: Send free slot data")
			return slot
		}
	}
	return geo.SlotInfo{
		X: -1,
	}
}
