// Code generated by "stringer -type=Action"; DO NOT EDIT.

package state

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ActionNone-0]
	_ = x[ActionLay-1]
	_ = x[ActionPull-2]
	_ = x[ActionEndTurn-3]
	_ = x[ActionNextTurn-4]
}

const _Action_name = "ActionNoneActionLayActionPullActionEndTurnActionNextTurn"

var _Action_index = [...]uint8{0, 10, 19, 29, 42, 56}

func (i Action) String() string {
	if i >= Action(len(_Action_index)-1) {
		return "Action(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Action_name[_Action_index[i]:_Action_index[i+1]]
}
