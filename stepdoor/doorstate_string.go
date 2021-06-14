// Code generated by "stringer -type=DoorState"; DO NOT EDIT.

package stepdoor

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Closed-0]
	_ = x[Semi-1]
	_ = x[Open-2]
}

const _DoorState_name = "ClosedSemiOpen"

var _DoorState_index = [...]uint8{0, 6, 10, 14}

func (i DoorState) String() string {
	if i < 0 || i >= DoorState(len(_DoorState_index)-1) {
		return "DoorState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DoorState_name[_DoorState_index[i]:_DoorState_index[i+1]]
}
