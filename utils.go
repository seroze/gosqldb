package gosqldb

func assert(condition bool, msg string) {
	if condition == false {
		panic(msg)
	}
}
