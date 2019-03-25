// Code generated by "string2enum -linecomment -type Visibility ../../ir/enum"; DO NOT EDIT.

package enum

import "fmt"
import "github.com/llir/llvm/ir/enum"

const _Visibility_name = "nonedefaulthiddenprotected"

var _Visibility_index = [...]uint8{0, 4, 11, 17, 26}

func VisibilityFromString(s string) enum.Visibility {
	if len(s) == 0 {
		return 0
	}
	for i := range _Visibility_index[:len(_Visibility_index)-1] {
		if s == _Visibility_name[_Visibility_index[i]:_Visibility_index[i+1]] {
			return enum.Visibility(i)
		}
	}
	panic(fmt.Errorf("unable to locate Visibility enum corresponding to %q", s))
}
