// Code generated by "stringer -type=BackendType"; DO NOT EDIT.

package backend

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LocalStorage-1]
	_ = x[S3-2]
	_ = x[SFTP-4]
}

const (
	_BackendType_name_0 = "LocalStorageS3"
	_BackendType_name_1 = "SFTP"
)

var (
	_BackendType_index_0 = [...]uint8{0, 12, 14}
)

func (i BackendType) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _BackendType_name_0[_BackendType_index_0[i]:_BackendType_index_0[i+1]]
	case i == 4:
		return _BackendType_name_1
	default:
		return "BackendType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
