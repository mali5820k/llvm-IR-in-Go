// Code generated by "string2enum -linecomment -type ReturnAttr ../../ir/enum"; DO NOT EDIT.

package enum

import "fmt"
import "github.com/llir/llvm/ir/enum"

const _ReturnAttr_name = "inregnoaliasnonnullsignextzeroext"

var _ReturnAttr_index = [...]uint8{0, 5, 12, 19, 26, 33}

func ReturnAttrFromString(s string) enum.ReturnAttr {
	if len(s) == 0 {
		return 0
	}
	for i := range _ReturnAttr_index[:len(_ReturnAttr_index)-1] {
		if s == _ReturnAttr_name[_ReturnAttr_index[i]:_ReturnAttr_index[i+1]] {
			return enum.ReturnAttr(i)
		}
	}
	panic(fmt.Errorf("unable to locate ReturnAttr enum corresponding to %q", s))
}
