package helper

import "log"

func PanicIfError(err error) {
	if err != nil {
		log.Println("An error occurred " + err.Error())

		panic(err)
	}
}
