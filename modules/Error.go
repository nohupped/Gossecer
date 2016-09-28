package modules

// CheckError will check whether the passed error variable is null, and if not, panics with the error.
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
