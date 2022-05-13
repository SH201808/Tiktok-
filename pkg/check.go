package pkg

func Check(username, password string) bool {
	if username == "" || password == "" {
		return false
	}
	return true
}
