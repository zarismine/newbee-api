package passwd

func Verfiy(password1, password2 string) bool {
	if Hash([]byte(password1)) == password2 {
		return true
	} else {
		return false
	}
}

