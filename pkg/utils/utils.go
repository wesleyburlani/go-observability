package utils

func PanicOnNotNil(err error) {
	if err != nil {
		panic(err)
	}
}
