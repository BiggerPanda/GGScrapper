package utility

func check(e error) {
	if e != nil {
		panic(e)
	}
}
