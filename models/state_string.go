// Code generated by "stringer -type=State"; DO NOT EDIT.

package models

import "fmt"

const _State_name = "UnavailableKeyExchangedMuxedCommandsExchanged"

var _State_index = [...]uint8{0, 11, 23, 28, 45}

func (i State) String() string {
	if i >= State(len(_State_index)-1) {
		return fmt.Sprintf("State(%d)", i)
	}
	return _State_name[_State_index[i]:_State_index[i+1]]
}
