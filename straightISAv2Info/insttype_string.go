// Code generated by "stringer -type=InstType isainfo.go"; DO NOT EDIT.

package straightISAv1Info

import "strconv"

const _InstType_name = "ZeroRegOneRegTwoReg"

var _InstType_index = [...]uint8{0, 7, 13, 19}

func (i InstType) String() string {
	if i < 0 || i >= InstType(len(_InstType_index)-1) {
		return "InstType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _InstType_name[_InstType_index[i]:_InstType_index[i+1]]
}
