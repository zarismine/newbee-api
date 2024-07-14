package stringopt

func SubStrLen(str string, length int) string {
	nameRune := []rune(str)
	if len(str) > length {
		return string(nameRune[:length-1]) + "..."
	}
	return string(nameRune)
}