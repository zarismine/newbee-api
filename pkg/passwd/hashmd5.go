package passwd

import "fmt"

func Hash(s []byte) string {
	if len(s) == 0 {
		return "0"
	}
	x := uint64(s[0]) << 7
	for _, v := range s {
		x = (1000003 * x) ^ uint64(v)
	}
	return Pad(fmt.Sprintf("%x",x),16)
}

func Pad(s string,l int) string {
	res := ""
	for i:=0;i<(l-len(s));i++{
		res += "0"
	}
	return res + s
}
