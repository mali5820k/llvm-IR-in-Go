// Code generated by "string2enum -linecomment -type TLSModel ../../ir/enum"; DO NOT EDIT.

package enum

import "fmt"
import "github.com/llir/llvm/ir/enum"

const _TLSModel_name = "nonegenericinitialexeclocaldynamiclocalexec"

var _TLSModel_index = [...]uint8{0, 4, 11, 22, 34, 43}

func TLSModelFromString(s string) enum.TLSModel {
	if len(s) == 0 {
		return 0
	}
	for i := range _TLSModel_index[:len(_TLSModel_index)-1] {
		if s == _TLSModel_name[_TLSModel_index[i]:_TLSModel_index[i+1]] {
			return enum.TLSModel(i)
		}
	}
	panic(fmt.Errorf("unable to locate TLSModel enum corresponding to %q", s))
}
