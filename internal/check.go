package internal

//Check panics if err != nil
func Check(e error) {
	if e != nil {
		panic(e)
	}
}
