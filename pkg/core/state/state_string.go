// Code generated by "stringer -type=State"; DO NOT EDIT.

package state

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[StateWaitForTurn-0]
	_ = x[StateMustLayOrPull-1]
	_ = x[StateMustLay-2]
	_ = x[StateCanLay-3]
}

const _State_name = "StateWaitForTurnStateMustLayOrPullStateMustLayStateCanLay"

var _State_index = [...]uint8{0, 16, 34, 46, 57}

func (i State) String() string {
	if i >= State(len(_State_index)-1) {
		return "State(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _State_name[_State_index[i]:_State_index[i+1]]
}
